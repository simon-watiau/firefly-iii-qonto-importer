package importer

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/simon-watiau/firefly-iii-qonto-importer/firefly"
	"github.com/simon-watiau/firefly-iii-qonto-importer/qonto"
	"go.uber.org/zap"
)

type Importer struct {
	logger                  *zap.Logger
	httpClient              *http.Client
	fireflyClient           *firefly.ClientWithResponses
	qontoClient             *qonto.ClientWithResponses
	qontoIban               string
	authQontoRequest        func(*http.Request)
	fireflyAssetAccountId   string
	fireflyRevenueAccountId string
	fireflyExpenseAccountId string
}

func (i *Importer) List(page *int) (*int, []qonto.Transaction, error) {
	resp, err := i.qontoClient.ListTransactionsWithResponse(
		context.Background(),
		&qonto.ListTransactionsParams{
			Iban: i.qontoIban,
			Page: page,
		},
		func(ctx context.Context, req *http.Request) error {
			if req.URL.RawQuery != "" {
				req.URL.RawQuery += "&"
			}
			req.URL.RawQuery += "includes[]=attachments"
			return nil
		},
	)

	if err != nil {
		return nil, nil, fmt.Errorf("API call failed: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, nil, fmt.Errorf("expected status code 200 got %d: %s", resp.StatusCode(), string(resp.Body))
	}

	return resp.JSON200.Meta.NextPage, *resp.JSON200.Transactions, nil
}

func (i *Importer) Import() error {
	i.logger.Info("Sync started")

	var page *int
	var transactions []qonto.Transaction
	var err error

	for next := true; next; next = page != nil {
		page, transactions, err = i.List(page)
		if err != nil {
			return fmt.Errorf("failed to list transactions: %w", err)
		}

		for _, t := range transactions {

			err = i.processQontoTransaction(t)

			if err != nil {
				i.logger.Error(fmt.Sprintf(
					"Failed to process transaction %s: %s",
					*t.TransactionId,
					err,
				),
				)
				continue
			}
		}
	}

	return nil
}

func (i *Importer) processQontoTransaction(t qonto.Transaction) error {
	if t.TransactionId == nil {
		return errors.New("Qonto transaction has no transactionId")
	}

	existingTransactions, err := i.listFirefly(*t.TransactionId)
	if err != nil {
		return fmt.Errorf(
			"failed to list existing Firefly transactions for Qonto transaction %s: %w",
			*t.TransactionId,
			err,
		)
	}

	var isNew = false
	var newAttachmentsAdded = 0

	if len(existingTransactions) == 0 {
		isNew = true
		newTransaction, err := i.createFireflyTransactionFromQonto(t)
		if err != nil {
			return fmt.Errorf(
				"failed to create Firefly transaction for transaction %s: %w",
				*t.TransactionId,
				err,
			)
		}

		existingTransactions = []firefly.TransactionRead{
			*newTransaction,
		}
	}

	for _, existingTransaction := range existingTransactions {
		if t.Attachments == nil {
			continue
		}

		for _, attachment := range *t.Attachments {
			if attachment.Url == nil {
				i.logger.Error(fmt.Sprintf(
					"Attachment has no id for Qonto transaction %s, skipping",
					*t.TransactionId,
				),
				)
				continue
			}

			hasAttachment, err := i.hasFireflyTransactionAttachment(
				existingTransaction.Id,
				*attachment.FileName,
			)

			if err != nil {
				return fmt.Errorf(
					"failed to check for existing attachments for Qonto transaction %s: %s",
					*t.TransactionId,
					err,
				)
			}

			if !hasAttachment {
				newAttachmentsAdded++
				content, err := i.DownloadQontoAttachment(*attachment.Url)

				if err != nil {
					i.logger.Error(
						fmt.Sprintf(
							"Failed to download attachment for Qonto transaction %s at %s, skipping",
							*t.TransactionId,
							*attachment.Url,
						),
					)
					continue
				}

				resp, err := i.fireflyClient.StoreAttachmentWithResponse(
					context.Background(),
					&firefly.StoreAttachmentParams{},
					firefly.AttachmentStore{
						AttachableId:   existingTransaction.Id,
						Filename:       *attachment.FileName,
						AttachableType: firefly.TransactionJournal,
						Notes:          stringPtr("imported"),
					},
				)

				if err != nil {
					i.logger.Error(
						fmt.Sprintf(
							"Failed to create attachment on Firefly transaction %s for Qonto transaction %s, skipping: %s",
							existingTransaction.Id,
							*t.TransactionId,
							err,
						),
					)
					continue
				}

				if resp.StatusCode() != 200 {
					i.logger.Error(
						fmt.Sprintf(
							"Failed to create attachment on Firefly transaction %s for Qonto transaction %s, skipping: %d!= 200 %s",
							existingTransaction.Id,
							*t.TransactionId,
							resp.StatusCode(),
							string(resp.Body),
						),
					)
					continue
				}

				_, err = i.fireflyClient.UploadAttachmentWithBody(
					context.Background(),
					resp.ApplicationvndApiJSON200.Data.Id,
					&firefly.UploadAttachmentParams{},
					"application/octet-stream",
					bytes.NewReader(content),
				)

				if err != nil {
					i.logger.Error(
						fmt.Sprintf(
							"Failed to upload attachment %s on Firefly transaction %s for Qonto transaction %s, skipping: %s",
							resp.ApplicationvndApiJSON200.Data.Id,
							existingTransaction.Id,
							*t.TransactionId,
							err,
						),
					)
					continue
				}
			}
		}
	}
	var matchingTransactionIds string
	for idx, et := range existingTransactions {
		if idx != 0 {
			matchingTransactionIds += ", "
		}
		matchingTransactionIds += et.Id
	}

	if isNew {
		i.logger.Info(
			fmt.Sprintf(
				"New transaction %s imported as transaction %s with %d/%d documents",
				*t.TransactionId,
				matchingTransactionIds,
				newAttachmentsAdded,
				len(*t.Attachments),
			),
		)
	} else {

		if newAttachmentsAdded == 0 {
			i.logger.Info(
				fmt.Sprintf(
					"Qonto transaction %s already in sync with %d firefly transaction(s) (%s)",
					*t.TransactionId,
					len(existingTransactions),
					matchingTransactionIds,
				),
			)
		} else {

			i.logger.Info(
				fmt.Sprintf(
					"%d/%d new attachments added for Qonto transaction %s on %d transactions (%s)",
					newAttachmentsAdded,
					len(*t.Attachments),
					*t.TransactionId,
					len(existingTransactions),
					matchingTransactionIds,
				),
			)
		}

	}

	return nil
}

func (i *Importer) DownloadQontoAttachment(url string) ([]byte, error) {
	req, err := http.NewRequest(
		"GET",
		url,
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := i.httpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("HTTP call failed: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("expected status code 200, got %d", resp.StatusCode)
	}

	content, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return content, nil
}

func (i *Importer) hasFireflyTransactionAttachment(transactionId string, fileName string) (bool, error) {
	existingAttachments, err := i.fireflyClient.ListAttachmentByTransactionWithResponse(
		context.Background(),
		transactionId,
		&firefly.ListAttachmentByTransactionParams{},
	)

	if err != nil {
		return false, fmt.Errorf("API call failed: %w", err)
	}

	if existingAttachments.StatusCode() != 200 {
		return false, fmt.Errorf(
			"expected status code 200, got %d: %s",
			existingAttachments.StatusCode(),
			string(existingAttachments.Body),
		)
	}

	for _, a := range existingAttachments.ApplicationvndApiJSON200.Data {
		if a.Attributes.Filename == fileName && *a.Attributes.Size != 0 {
			return true, nil
		}
	}
	return false, nil
}

func (i *Importer) listFirefly(qontoTransactionId string) ([]firefly.TransactionRead, error) {
	existingTransactionsOnExternalId, err := i.fireflyClient.SearchTransactionsWithResponse(
		context.Background(),
		&firefly.SearchTransactionsParams{
			Query: fmt.Sprintf(
				"external_id_is:%s",
				qontoTransactionId,
			),
		},
	)

	if err != nil {
		return nil, fmt.Errorf("API call failed: %w", err)
	}

	if existingTransactionsOnExternalId.StatusCode() != 200 {
		return nil, fmt.Errorf(
			"expected status code 200, got %d: %s",
			existingTransactionsOnExternalId.StatusCode(),
			string(existingTransactionsOnExternalId.Body),
		)
	}

	existingTransactionsOnDescription, err := i.fireflyClient.SearchTransactionsWithResponse(
		context.Background(),
		&firefly.SearchTransactionsParams{
			Query: fmt.Sprintf(
				`notes_contains:"[%s]"`,
				qontoTransactionId,
			),
		},
	)

	if err != nil {
		return nil, fmt.Errorf("API call failed: %w", err)
	}

	if existingTransactionsOnDescription.StatusCode() != 200 {
		return nil, fmt.Errorf(
			"expected status code 200, got %d: %s",
			existingTransactionsOnExternalId.StatusCode(),
			string(existingTransactionsOnExternalId.Body),
		)
	}

	var allTransactions = make(
		[]firefly.TransactionRead,
		0,
		len(existingTransactionsOnDescription.ApplicationvndApiJSON200.Data)+len(existingTransactionsOnExternalId.ApplicationvndApiJSON200.Data),
	)

	allTransactions = append(allTransactions, existingTransactionsOnDescription.ApplicationvndApiJSON200.Data...)
	allTransactions = append(allTransactions, existingTransactionsOnExternalId.ApplicationvndApiJSON200.Data...)

	return allTransactions, nil
}

func (i *Importer) createFireflyTransactionFromQonto(t qonto.Transaction) (*firefly.TransactionRead, error) {
	var operationType firefly.TransactionTypeProperty
	var sourceAccountId, destinationAccountId string
	switch *t.Side {
	case "credit":
		operationType = firefly.Deposit
		sourceAccountId = i.fireflyRevenueAccountId
		destinationAccountId = i.fireflyAssetAccountId
	case "debit":
		operationType = firefly.Withdrawal
		sourceAccountId = i.fireflyAssetAccountId
		destinationAccountId = i.fireflyExpenseAccountId
	}
	date, err := time.Parse("2006-01-02T15:04:05.999Z", *t.SettledAt)

	if err != nil {
		return nil, fmt.Errorf("failed to parse settledAt date: %w", err)
	}

	resp, err := i.fireflyClient.StoreTransactionWithResponse(
		context.Background(),
		&firefly.StoreTransactionParams{},
		firefly.TransactionStore{
			ApplyRules:           boolPtr(false),
			ErrorIfDuplicateHash: boolPtr(false),
			Transactions: []firefly.TransactionSplitStore{
				{
					Type:          operationType,
					Date:          date,
					Order:         int32Ptr(0),
					Amount:        fmt.Sprint(*t.Amount),
					SourceId:      stringPtr(sourceAccountId),
					DestinationId: stringPtr(destinationAccountId),
					Description:   *t.Label,
					ExternalId:    t.TransactionId,
				},
			},
		},
	)

	if err != nil {
		return nil, fmt.Errorf("API call failed: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf(
			"expected status code 200 got %d with body %s",
			resp.StatusCode(),
			string(resp.Body),
		)
	}

	return &resp.ApplicationvndApiJSON200.Data, nil
}

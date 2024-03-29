// Package qonto provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.1.0 DO NOT EDIT.
package qonto

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/oapi-codegen/runtime"
)

const (
	OAuthScopes     = "OAuth.Scopes"
	SecretKeyScopes = "SecretKey.Scopes"
)

// Transaction defines model for Transaction.
type Transaction struct {
	Amount             *float32  `json:"amount,omitempty"`
	AmountCents        *int      `json:"amount_cents,omitempty"`
	AttachmentIds      *[]string `json:"attachment_ids,omitempty"`
	AttachmentLost     *bool     `json:"attachment_lost,omitempty"`
	AttachmentRequired *bool     `json:"attachment_required,omitempty"`
	Attachments        *[]struct {
		CreatedAt           *string `json:"created_at,omitempty"`
		FileContentType     *string `json:"file_content_type,omitempty"`
		FileName            *string `json:"file_name,omitempty"`
		FileSize            *string `json:"file_size,omitempty"`
		Id                  *string `json:"id,omitempty"`
		ProbativeAttachment *struct {
			Status *string `json:"status,omitempty"`
		} `json:"probative_attachment,omitempty"`
		Url *string `json:"url,omitempty"`
	} `json:"attachments,omitempty"`
	CardLastDigits *string `json:"card_last_digits"`
	Category       *string `json:"category,omitempty"`
	Check          *struct {
		CheckKey    *string `json:"check_key,omitempty"`
		CheckNumber *string `json:"check_number,omitempty"`
	} `json:"check"`
	Currency    *string `json:"currency,omitempty"`
	DirectDebit *struct {
		CounterpartyAccountNumber        *string `json:"counterparty_account_number,omitempty"`
		CounterpartyAccountNumberFormat  *string `json:"counterparty_account_number_format,omitempty"`
		CounterpartyBankIdentifier       *string `json:"counterparty_bank_identifier,omitempty"`
		CounterpartyBankIdentifierFormat *string `json:"counterparty_bank_identifier_format,omitempty"`
	} `json:"direct_debit"`
	EmittedAt            *string `json:"emitted_at,omitempty"`
	FinancingInstallment *struct {
		CurrentInstallmentNumber *int `json:"current_installment_number,omitempty"`
		TotalInstallmentsNumber  *int `json:"total_installments_number,omitempty"`
	} `json:"financing_installment"`
	Id     *string `json:"id,omitempty"`
	Income *struct {
		CounterpartyAccountNumber        *string `json:"counterparty_account_number,omitempty"`
		CounterpartyAccountNumberFormat  *string `json:"counterparty_account_number_format,omitempty"`
		CounterpartyBankIdentifier       *string `json:"counterparty_bank_identifier,omitempty"`
		CounterpartyBankIdentifierFormat *string `json:"counterparty_bank_identifier_format,omitempty"`
	} `json:"income"`
	InitiatorId *string   `json:"initiator_id"`
	Label       *string   `json:"label,omitempty"`
	LabelIds    *[]string `json:"label_ids,omitempty"`
	Labels      *[]struct {
		Id       *string      `json:"id,omitempty"`
		Name     *string      `json:"name,omitempty"`
		ParentId *interface{} `json:"parent_id"`
	} `json:"labels,omitempty"`
	LocalAmount      *float32 `json:"local_amount,omitempty"`
	LocalAmountCents *int     `json:"local_amount_cents,omitempty"`
	LocalCurrency    *string  `json:"local_currency,omitempty"`
	Note             *string  `json:"note,omitempty"`
	OperationType    *string  `json:"operation_type,omitempty"`
	PagopaPayment    *struct {
		CreditorFiscalCode *string `json:"creditor_fiscal_code,omitempty"`
		Iuv                *string `json:"iuv,omitempty"`
		NoticeNumber       *string `json:"notice_number,omitempty"`
	} `json:"pagopa_payment"`
	Reference           *interface{} `json:"reference"`
	SettledAt           *string      `json:"settled_at,omitempty"`
	SettledBalance      *float32     `json:"settled_balance,omitempty"`
	SettledBalanceCents *int         `json:"settled_balance_cents,omitempty"`
	Side                *string      `json:"side,omitempty"`
	Status              *string      `json:"status,omitempty"`
	SubjectType         *string      `json:"subject_type,omitempty"`
	SwiftIncome         *struct {
		CounterpartyAccountNumber        *string `json:"counterparty_account_number,omitempty"`
		CounterpartyAccountNumberFormat  *string `json:"counterparty_account_number_format,omitempty"`
		CounterpartyBankIdentifier       *string `json:"counterparty_bank_identifier,omitempty"`
		CounterpartyBankIdentifierFormat *string `json:"counterparty_bank_identifier_format,omitempty"`
	} `json:"swift_income"`
	TransactionId *string `json:"transaction_id,omitempty"`
	Transfer      *struct {
		CounterpartyAccountNumber        *string `json:"counterparty_account_number,omitempty"`
		CounterpartyAccountNumberFormat  *string `json:"counterparty_account_number_format,omitempty"`
		CounterpartyBankIdentifier       *string `json:"counterparty_bank_identifier,omitempty"`
		CounterpartyBankIdentifierFormat *string `json:"counterparty_bank_identifier_format,omitempty"`
	} `json:"transfer"`
	UpdatedAt      *string  `json:"updated_at,omitempty"`
	VatAmount      *float32 `json:"vat_amount"`
	VatAmountCents *int     `json:"vat_amount_cents"`
	VatRate        *float32 `json:"vat_rate"`
}

// ListTransactionsParams defines parameters for ListTransactions.
type ListTransactionsParams struct {
	Iban string `form:"iban" json:"iban"`
	Page *int   `form:"page,omitempty" json:"page,omitempty"`
}

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// GetAttachment request
	GetAttachment(ctx context.Context, id string, reqEditors ...RequestEditorFn) (*http.Response, error)

	// ListTransactions request
	ListTransactions(ctx context.Context, params *ListTransactionsParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// ListTransactionAttachments request
	ListTransactionAttachments(ctx context.Context, id string, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) GetAttachment(ctx context.Context, id string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetAttachmentRequest(c.Server, id)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) ListTransactions(ctx context.Context, params *ListTransactionsParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewListTransactionsRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) ListTransactionAttachments(ctx context.Context, id string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewListTransactionAttachmentsRequest(c.Server, id)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewGetAttachmentRequest generates requests for GetAttachment
func NewGetAttachmentRequest(server string, id string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "id", runtime.ParamLocationPath, id)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v2/attachments/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewListTransactionsRequest generates requests for ListTransactions
func NewListTransactionsRequest(server string, params *ListTransactionsParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v2/transactions")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if queryFrag, err := runtime.StyleParamWithLocation("form", true, "iban", runtime.ParamLocationQuery, params.Iban); err != nil {
			return nil, err
		} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
			return nil, err
		} else {
			for k, v := range parsed {
				for _, v2 := range v {
					queryValues.Add(k, v2)
				}
			}
		}

		if params.Page != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "page", runtime.ParamLocationQuery, *params.Page); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewListTransactionAttachmentsRequest generates requests for ListTransactionAttachments
func NewListTransactionAttachmentsRequest(server string, id string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "id", runtime.ParamLocationPath, id)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v2/transactions/%s/attachments", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// GetAttachmentWithResponse request
	GetAttachmentWithResponse(ctx context.Context, id string, reqEditors ...RequestEditorFn) (*GetAttachmentResponse, error)

	// ListTransactionsWithResponse request
	ListTransactionsWithResponse(ctx context.Context, params *ListTransactionsParams, reqEditors ...RequestEditorFn) (*ListTransactionsResponse, error)

	// ListTransactionAttachmentsWithResponse request
	ListTransactionAttachmentsWithResponse(ctx context.Context, id string, reqEditors ...RequestEditorFn) (*ListTransactionAttachmentsResponse, error)
}

type GetAttachmentResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *struct {
		Attachment *struct {
			CreatedAt           *string `json:"created_at,omitempty"`
			FileContentType     *string `json:"file_content_type,omitempty"`
			FileName            *string `json:"file_name,omitempty"`
			FileSize            *string `json:"file_size,omitempty"`
			Id                  *string `json:"id,omitempty"`
			ProbativeAttachment *struct {
				FileContentType *string `json:"file_content_type,omitempty"`
				FileName        *string `json:"file_name,omitempty"`
				FileSize        *string `json:"file_size,omitempty"`
				Status          *string `json:"status,omitempty"`
				Url             *string `json:"url,omitempty"`
			} `json:"probative_attachment,omitempty"`
			Url *string `json:"url,omitempty"`
		} `json:"attachment,omitempty"`
	}
	JSON401 *struct {
		Errors *[]struct {
			Code   *string `json:"code,omitempty"`
			Detail *string `json:"detail,omitempty"`
		} `json:"errors,omitempty"`
	}
}

// Status returns HTTPResponse.Status
func (r GetAttachmentResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetAttachmentResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type ListTransactionsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *struct {
		Meta *struct {
			CurrentPage *int `json:"current_page,omitempty"`
			NextPage    *int `json:"next_page"`
			PerPage     *int `json:"per_page,omitempty"`
			PrevPage    *int `json:"prev_page"`
			TotalCount  *int `json:"total_count,omitempty"`
			TotalPages  *int `json:"total_pages,omitempty"`
		} `json:"meta,omitempty"`
		Transactions *[]Transaction `json:"transactions,omitempty"`
	}
	JSON400 *struct {
		Errors *[]struct {
			Code    *string `json:"code,omitempty"`
			Message *string `json:"message,omitempty"`
		} `json:"errors,omitempty"`
	}
	JSON422 *struct {
		Errors []struct {
			Code   string  `json:"code"`
			Detail *string `json:"detail,omitempty"`
			Source *struct {
				Pointer *string `json:"pointer,omitempty"`
			} `json:"source,omitempty"`
		} `json:"errors"`
	}
}

// Status returns HTTPResponse.Status
func (r ListTransactionsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r ListTransactionsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type ListTransactionAttachmentsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *struct {
		Attachments *[]struct {
			CreatedAt           *string `json:"created_at,omitempty"`
			FileContentType     *string `json:"file_content_type,omitempty"`
			FileName            *string `json:"file_name,omitempty"`
			FileSize            *string `json:"file_size,omitempty"`
			Id                  *string `json:"id,omitempty"`
			ProbativeAttachment *struct {
				FileContentType *string `json:"file_content_type,omitempty"`
				FileName        *string `json:"file_name,omitempty"`
				FileSize        *string `json:"file_size,omitempty"`
				Status          *string `json:"status,omitempty"`
				Url             *string `json:"url,omitempty"`
			} `json:"probative_attachment,omitempty"`
			Url *string `json:"url,omitempty"`
		} `json:"attachments,omitempty"`
	}
	JSON401 *struct {
		Errors *[]struct {
			Code   *string `json:"code,omitempty"`
			Detail *string `json:"detail,omitempty"`
		} `json:"errors,omitempty"`
	}
	JSON404 *struct {
		Errors *[]struct {
			Code   *string `json:"code,omitempty"`
			Detail *string `json:"detail,omitempty"`
			Source *struct {
				Parameter *string `json:"parameter,omitempty"`
			} `json:"source,omitempty"`
		} `json:"errors,omitempty"`
	}
}

// Status returns HTTPResponse.Status
func (r ListTransactionAttachmentsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r ListTransactionAttachmentsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// GetAttachmentWithResponse request returning *GetAttachmentResponse
func (c *ClientWithResponses) GetAttachmentWithResponse(ctx context.Context, id string, reqEditors ...RequestEditorFn) (*GetAttachmentResponse, error) {
	rsp, err := c.GetAttachment(ctx, id, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetAttachmentResponse(rsp)
}

// ListTransactionsWithResponse request returning *ListTransactionsResponse
func (c *ClientWithResponses) ListTransactionsWithResponse(ctx context.Context, params *ListTransactionsParams, reqEditors ...RequestEditorFn) (*ListTransactionsResponse, error) {
	rsp, err := c.ListTransactions(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseListTransactionsResponse(rsp)
}

// ListTransactionAttachmentsWithResponse request returning *ListTransactionAttachmentsResponse
func (c *ClientWithResponses) ListTransactionAttachmentsWithResponse(ctx context.Context, id string, reqEditors ...RequestEditorFn) (*ListTransactionAttachmentsResponse, error) {
	rsp, err := c.ListTransactionAttachments(ctx, id, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseListTransactionAttachmentsResponse(rsp)
}

// ParseGetAttachmentResponse parses an HTTP response from a GetAttachmentWithResponse call
func ParseGetAttachmentResponse(rsp *http.Response) (*GetAttachmentResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetAttachmentResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest struct {
			Attachment *struct {
				CreatedAt           *string `json:"created_at,omitempty"`
				FileContentType     *string `json:"file_content_type,omitempty"`
				FileName            *string `json:"file_name,omitempty"`
				FileSize            *string `json:"file_size,omitempty"`
				Id                  *string `json:"id,omitempty"`
				ProbativeAttachment *struct {
					FileContentType *string `json:"file_content_type,omitempty"`
					FileName        *string `json:"file_name,omitempty"`
					FileSize        *string `json:"file_size,omitempty"`
					Status          *string `json:"status,omitempty"`
					Url             *string `json:"url,omitempty"`
				} `json:"probative_attachment,omitempty"`
				Url *string `json:"url,omitempty"`
			} `json:"attachment,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest struct {
			Errors *[]struct {
				Code   *string `json:"code,omitempty"`
				Detail *string `json:"detail,omitempty"`
			} `json:"errors,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON401 = &dest

	}

	return response, nil
}

// ParseListTransactionsResponse parses an HTTP response from a ListTransactionsWithResponse call
func ParseListTransactionsResponse(rsp *http.Response) (*ListTransactionsResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &ListTransactionsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest struct {
			Meta *struct {
				CurrentPage *int `json:"current_page,omitempty"`
				NextPage    *int `json:"next_page"`
				PerPage     *int `json:"per_page,omitempty"`
				PrevPage    *int `json:"prev_page"`
				TotalCount  *int `json:"total_count,omitempty"`
				TotalPages  *int `json:"total_pages,omitempty"`
			} `json:"meta,omitempty"`
			Transactions *[]Transaction `json:"transactions,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest struct {
			Errors *[]struct {
				Code    *string `json:"code,omitempty"`
				Message *string `json:"message,omitempty"`
			} `json:"errors,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 422:
		var dest struct {
			Errors []struct {
				Code   string  `json:"code"`
				Detail *string `json:"detail,omitempty"`
				Source *struct {
					Pointer *string `json:"pointer,omitempty"`
				} `json:"source,omitempty"`
			} `json:"errors"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON422 = &dest

	}

	return response, nil
}

// ParseListTransactionAttachmentsResponse parses an HTTP response from a ListTransactionAttachmentsWithResponse call
func ParseListTransactionAttachmentsResponse(rsp *http.Response) (*ListTransactionAttachmentsResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &ListTransactionAttachmentsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest struct {
			Attachments *[]struct {
				CreatedAt           *string `json:"created_at,omitempty"`
				FileContentType     *string `json:"file_content_type,omitempty"`
				FileName            *string `json:"file_name,omitempty"`
				FileSize            *string `json:"file_size,omitempty"`
				Id                  *string `json:"id,omitempty"`
				ProbativeAttachment *struct {
					FileContentType *string `json:"file_content_type,omitempty"`
					FileName        *string `json:"file_name,omitempty"`
					FileSize        *string `json:"file_size,omitempty"`
					Status          *string `json:"status,omitempty"`
					Url             *string `json:"url,omitempty"`
				} `json:"probative_attachment,omitempty"`
				Url *string `json:"url,omitempty"`
			} `json:"attachments,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest struct {
			Errors *[]struct {
				Code   *string `json:"code,omitempty"`
				Detail *string `json:"detail,omitempty"`
			} `json:"errors,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON401 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest struct {
			Errors *[]struct {
				Code   *string `json:"code,omitempty"`
				Detail *string `json:"detail,omitempty"`
				Source *struct {
					Parameter *string `json:"parameter,omitempty"`
				} `json:"source,omitempty"`
			} `json:"errors,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	}

	return response, nil
}

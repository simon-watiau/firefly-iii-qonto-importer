package importer

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
	"github.com/simon-watiau/firefly-iii-qonto-importer/firefly"
	"github.com/simon-watiau/firefly-iii-qonto-importer/qonto"
	"go.uber.org/zap"
)

type FireflyConfig struct {
	BaseUrl          string
	Token            string
	AssetAccountId   string
	RevenueAccountId string
	ExpenseAccountId string
}

type QontoConfig struct {
	BaseUrl  string
	Login    string
	Password string
	Iban     string
}

type ImporterConfig struct {
	Firefly FireflyConfig
	Qonto   QontoConfig
}

func NewImporter(
	logger *zap.Logger,
	httpClient *http.Client,
	config ImporterConfig,
) (*Importer, error) {

	if logger == nil {
		return nil, errors.New("missing logger")
	}

	fireflyClient, err := createFireflyClient(
		config.Firefly.BaseUrl,
		config.Firefly.Token,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create Firefly client: %w", err)
	}

	authQontoRequest := func(r *http.Request) {
		r.Header.Set("Authorization", config.Qonto.Login+":"+config.Qonto.Password)
	}

	qontoClient, err := createQontoClient(
		config.Qonto.BaseUrl,
		authQontoRequest,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create Qonto client: %w", err)
	}

	return &Importer{
		logger:                  logger,
		httpClient:              httpClient,
		fireflyClient:           fireflyClient,
		qontoClient:             qontoClient,
		authQontoRequest:        authQontoRequest,
		fireflyAssetAccountId:   config.Firefly.AssetAccountId,
		fireflyRevenueAccountId: config.Firefly.RevenueAccountId,
		fireflyExpenseAccountId: config.Firefly.ExpenseAccountId,
		qontoIban:               config.Qonto.Iban,
	}, nil
}

func createFireflyClient(
	baseUrl string,
	token string,
) (*firefly.ClientWithResponses, error) {
	authProvider, err := securityprovider.NewSecurityProviderBearerToken(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create authorization provider: %w", err)
	}

	fireflyClient, err := firefly.NewClientWithResponses(
		baseUrl,
		firefly.WithRequestEditorFn(authProvider.Intercept),
		firefly.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
			req.Header.Set("accept", "application/json; charset=utf-8")
			req.Header.Set("Content-Type", "application/json; charset=utf-8")
			return nil
		}),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	return fireflyClient, nil
}

func createQontoClient(
	baseUrl string,
	authRequest func(r *http.Request),
) (*qonto.ClientWithResponses, error) {

	qontoClient, err := qonto.NewClientWithResponses(
		baseUrl,
		qonto.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
			authRequest(req)
			req.Header.Set("accept", "application/json; charset=utf-8")
			req.Header.Set("Content-Type", "application/json; charset=utf-8")
			return nil
		}),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	return qontoClient, nil
}

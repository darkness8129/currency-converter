package currapi

import (
	"context"
	"darkness8129/currency-converter/currency-service/app/service"
	"darkness8129/currency-converter/currency-service/config"
	"darkness8129/currency-converter/currency-service/packages/logging"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

var _ service.CurrencyAPI = (*oxrAPI)(nil)

type oxrAPI struct {
	logger     *logging.Logger
	cfg        *config.Config
	httpClient *resty.Client
}

func NewOxrAPI(logger *logging.Logger, cfg *config.Config) *oxrAPI {
	httpClient := resty.New().
		SetBaseURL(cfg.OXR.BaseURL).
		SetHeader("Content-type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Token %s", cfg.OXR.AppID))

	return &oxrAPI{
		logger:     logger.Named("oxrAPI"),
		cfg:        cfg,
		httpClient: httpClient,
	}
}

type getRateResult struct {
	Rates struct {
		UAH float64
	}
}

func (a *oxrAPI) GetRate(ctx context.Context) (float64, error) {
	logger := a.logger.Named("GetRate")

	var result getRateResult
	res, err := a.httpClient.R().
		SetQueryParams(map[string]string{
			"base":    "USD",
			"symbols": "UAH",
		}).
		SetResult(&result).
		Get("/latest.json")
	if err != nil {
		logger.Error("failed to send get rate request", "err", err)
		return 0, fmt.Errorf("failed to send get rate request: %w", err)
	}
	// TODO: parse errors and return custom
	if res.StatusCode() != http.StatusOK {
		logger.Error("failed to get rate", "resBody", res.String(), "statusCode", res.StatusCode())
		return 0, fmt.Errorf("failed to get rate: http status %d, body %s", res.StatusCode(), res.String())
	}

	rate := result.Rates.UAH

	logger.Info("successfully got rate", "rate", rate)
	return rate, nil
}

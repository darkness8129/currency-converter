package service

import (
	"context"
	"darkness8129/currency-converter/currency-service/packages/errs"
	"darkness8129/currency-converter/currency-service/packages/logging"
	"fmt"
)

var _ CurrencyService = (*currencyService)(nil)

type currencyService struct {
	logger *logging.Logger
	apis   APIs
}

func NewCurrencyService(logger *logging.Logger, apis APIs) *currencyService {
	return &currencyService{
		logger: logger.Named("currencyService"),
		apis:   apis,
	}
}

func (s *currencyService) GetRate(ctx context.Context) (float64, error) {
	logger := s.logger.Named("GetRate")

	rate, err := s.apis.Currency.GetRate(ctx)
	if err != nil {
		if errs.IsCustom(err) {
			logger.Info(err.Error())
			return 0, err
		}

		logger.Error("failed to get rate", "err", err)
		return 0, fmt.Errorf("failed to get rate: %w", err)
	}

	logger.Info("successfully got rate", "rate", rate)
	return rate, nil
}

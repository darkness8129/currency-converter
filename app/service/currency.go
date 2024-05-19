package service

import (
	"context"
	"darkness8129/currency-converter/packages/logging"
)

var _ CurrencyService = (*currencyService)(nil)

type currencyService struct {
	logger *logging.Logger
}

func NewCurrencyService(logger *logging.Logger) *currencyService {
	return &currencyService{
		logger: logger.Named("currencyService"),
	}
}

func (s *currencyService) GetRate(ctx context.Context, opt GetRateOpt) error {
	logger := s.logger.Named("GetRate")

	logger.Info("successfully got rate")
	return nil
}

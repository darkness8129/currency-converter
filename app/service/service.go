package service

import "context"

type Services struct {
	Currency CurrencyService
}

type CurrencyService interface {
	GetRate(ctx context.Context, opt GetRateOpt) error
}

type GetRateOpt struct {
}

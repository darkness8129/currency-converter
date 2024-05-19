package service

import "context"

type Services struct {
	Currency CurrencyService
}

type CurrencyService interface {
	GetRate(ctx context.Context) (float64, error)
}

type APIs struct {
	Currency CurrencyAPI
}

type CurrencyAPI interface {
	GetRate(ctx context.Context) (float64, error)
}

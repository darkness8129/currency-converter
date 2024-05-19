package service

import (
	"context"
	"darkness8129/currency-converter/app/entity"
	"darkness8129/currency-converter/packages/errs"
)

type Services struct {
	Currency     CurrencyService
	Subscription SubscriptionService
}

type CurrencyService interface {
	GetRate(ctx context.Context) (float64, error)
}

type SubscriptionService interface {
	Subscribe(ctx context.Context, email string) error
	StartMailing(ctx context.Context)
}

const (
	already_exists = "already_exists"
)

var (
	ErrSubscriberAlreadyExists = errs.New(errs.Options{Message: "subscriber already exists", Code: already_exists})
)

type Storages struct {
	Subscriber SubscriberStorage
}

type SubscriberStorage interface {
	Create(ctx context.Context, sub *entity.Subscriber) (*entity.Subscriber, error)
	Get(ctx context.Context, email string) (*entity.Subscriber, error)
	List(ctx context.Context) ([]entity.Subscriber, error)
}

type APIs struct {
	Currency CurrencyAPI
	Email    EmailAPI
}

type CurrencyAPI interface {
	GetRate(ctx context.Context) (float64, error)
}

type EmailAPI interface {
	Send(ctx context.Context, opt SendOpt) error
}

type SendOpt struct {
	To      string
	Subject string
	Body    string
}

package service

import (
	"context"
	"darkness8129/currency-converter/subscription-service/app/entity"
	"darkness8129/currency-converter/subscription-service/packages/errs"
)

type Services struct {
	Subscription SubscriptionService
}

type SubscriptionService interface {
	Subscribe(ctx context.Context, email string) error
}

const (
	already_exists = "already_exists"
)

var (
	ErrSubscriptionAlreadyExists = errs.New(errs.Options{Message: "subscription already exists", Code: already_exists})
)

type Storages struct {
	Subscription SubscriptionStorage
}

type SubscriptionStorage interface {
	Create(ctx context.Context, subscription *entity.Subscription) (*entity.Subscription, error)
	Get(ctx context.Context, email string) (*entity.Subscription, error)
}

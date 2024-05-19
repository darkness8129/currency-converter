package service

import (
	"context"
	"darkness8129/currency-converter/subscription-service/app/entity"
	"darkness8129/currency-converter/subscription-service/packages/logging"
	"fmt"
)

var _ SubscriptionService = (*subscriptionService)(nil)

type subscriptionService struct {
	logger   *logging.Logger
	storages Storages
}

func NewSubscriptionService(logger *logging.Logger, storages Storages) *subscriptionService {
	return &subscriptionService{
		logger:   logger.Named("subscriptionService"),
		storages: storages,
	}
}

func (s *subscriptionService) Subscribe(ctx context.Context, email string) error {
	logger := s.logger.Named("Subscribe")

	subscription, err := s.storages.Subscription.Get(ctx, email)
	if err != nil {
		logger.Error("failed to get subscription", "err", err)
		return fmt.Errorf("failed to get subscription: %w", err)
	}
	if subscription != nil {
		logger.Info("subscription already exists", "email", email)
		return ErrSubscriptionAlreadyExists
	}

	_, err = s.storages.Subscription.Create(ctx, &entity.Subscription{Email: email})
	if err != nil {
		logger.Error("failed to create subscription", "err", err)
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	logger.Info("successfully subscribed", "email", email)
	return nil
}

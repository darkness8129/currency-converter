package storage

import (
	"context"

	"errors"
	"fmt"

	"darkness8129/currency-converter/subscription-service/app/entity"
	"darkness8129/currency-converter/subscription-service/app/service"
	"darkness8129/currency-converter/subscription-service/packages/logging"

	"gorm.io/gorm"
)

var _ service.SubscriptionStorage = (*subscriptionStorage)(nil)

type subscriptionStorage struct {
	logger *logging.Logger
	db     *gorm.DB
}

func NewSubscriptionStorage(logger *logging.Logger, db *gorm.DB) *subscriptionStorage {
	return &subscriptionStorage{
		logger: logger.Named("subscriptionStorage"),
		db:     db,
	}
}

func (s *subscriptionStorage) Create(ctx context.Context, subscription *entity.Subscription) (*entity.Subscription, error) {
	logger := s.logger.Named("Create")

	err := s.db.Create(subscription).Error
	if err != nil {
		logger.Error("failed to create subscription", "err", err)
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	logger.Info("successfully created subscription", "subscription", subscription)
	return subscription, nil
}

func (s *subscriptionStorage) Get(ctx context.Context, email string) (*entity.Subscription, error) {
	logger := s.logger.Named("Get")

	var subscription entity.Subscription
	err := s.db.
		Where(entity.Subscription{Email: email}).
		First(&subscription).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Info("subscription not found", "email", email)
		return nil, nil
	}
	if err != nil {
		logger.Error("failed to get subscription", "err", err)
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	logger.Info("successfully got subscription", "subscription", subscription)
	return &subscription, nil
}

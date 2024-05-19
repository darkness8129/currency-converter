package storage

import (
	"context"

	"errors"
	"fmt"

	"darkness8129/currency-converter/app/entity"
	"darkness8129/currency-converter/app/service"
	"darkness8129/currency-converter/packages/logging"

	"gorm.io/gorm"
)

var _ service.SubscriberStorage = (*subscriberStorage)(nil)

type subscriberStorage struct {
	logger *logging.Logger
	db     *gorm.DB
}

func NewSubscriberStorage(logger *logging.Logger, db *gorm.DB) *subscriberStorage {
	return &subscriberStorage{
		logger: logger.Named("subscriberStorage"),
		db:     db,
	}
}

func (s *subscriberStorage) Create(ctx context.Context, sub *entity.Subscriber) (*entity.Subscriber, error) {
	logger := s.logger.Named("Create")

	err := s.db.Create(sub).Error
	if err != nil {
		logger.Error("failed to create subscriber", "err", err)
		return nil, fmt.Errorf("failed to create subscriber: %w", err)
	}

	logger.Info("successfully created subscriber", "sub", sub)
	return sub, nil
}

func (s *subscriberStorage) Get(ctx context.Context, email string) (*entity.Subscriber, error) {
	logger := s.logger.Named("Get")

	var sub entity.Subscriber
	err := s.db.
		Where(entity.Subscriber{Email: email}).
		First(&sub).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Info("subscriber not found", "email", email)
		return nil, nil
	}
	if err != nil {
		logger.Error("failed to get subscriber", "err", err)
		return nil, fmt.Errorf("failed to get subscriber: %w", err)
	}

	logger.Info("successfully got subscriber", "sub", sub)
	return &sub, nil
}

func (s *subscriberStorage) List(ctx context.Context) ([]entity.Subscriber, error) {
	logger := s.logger.Named("List")

	var subs []entity.Subscriber
	err := s.db.Find(&subs).Error
	if err != nil {
		logger.Error("failed to list subscribers", "err", err)
		return nil, fmt.Errorf("failed to list subscribers: %w", err)
	}

	logger.Info("successfully listed subscribers", "subs", subs)
	return subs, nil
}

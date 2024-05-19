package service

import (
	"context"
	"darkness8129/currency-converter/app/entity"
	"darkness8129/currency-converter/config"
	"darkness8129/currency-converter/packages/errs"
	"darkness8129/currency-converter/packages/logging"
	"fmt"
	"time"
)

var _ SubscriptionService = (*subscriptionService)(nil)

type subscriptionService struct {
	logger          *logging.Logger
	cfg             *config.Config
	storages        Storages
	apis            APIs
	currencyService CurrencyService
}

func NewSubscriptionService(logger *logging.Logger, cfg *config.Config, storages Storages, apis APIs, currencyService CurrencyService) *subscriptionService {
	return &subscriptionService{
		logger:          logger.Named("subscriptionService"),
		cfg:             cfg,
		storages:        storages,
		apis:            apis,
		currencyService: currencyService,
	}
}

func (s *subscriptionService) Subscribe(ctx context.Context, email string) error {
	logger := s.logger.Named("Subscribe")

	sub, err := s.storages.Subscriber.Get(ctx, email)
	if err != nil {
		logger.Error("failed to get subscriber", "err", err)
		return fmt.Errorf("failed to get subscriber: %w", err)
	}
	if sub != nil {
		logger.Info("subscriber already exists", "email", email)
		return ErrSubscriberAlreadyExists
	}

	_, err = s.storages.Subscriber.Create(ctx, &entity.Subscriber{Email: email})
	if err != nil {
		logger.Error("failed to create subscriber", "err", err)
		return fmt.Errorf("failed to create subscriber: %w", err)
	}

	logger.Info("successfully subscribed", "email", email)
	return nil
}

func (s *subscriptionService) StartMailing(ctx context.Context) {
	logger := s.logger.Named("StartMailing")

	ticker := time.NewTicker(s.cfg.App.MailingPeriod)
	defer ticker.Stop()

	for range ticker.C {
		err := s.sendEmails(ctx)
		if err != nil {
			logger.Error("failed to send emails for subscribers", "err", err)
		}
	}
}

func (s *subscriptionService) sendEmails(ctx context.Context) error {
	logger := s.logger.Named("sendEmails")

	rate, err := s.currencyService.GetRate(ctx)
	if err != nil {
		if errs.IsCustom(err) {
			logger.Info(err.Error())
			return err
		}

		logger.Error("failed to get rate", "err", err)
		return fmt.Errorf("failed to get rate: %w", err)
	}
	logger.Debug("got rate", "rate", rate)

	subs, err := s.storages.Subscriber.List(ctx)
	if err != nil {
		logger.Error("failed to list subscribers", "err", err)
		return fmt.Errorf("failed to get subscribers: %w", err)
	}
	logger.Debug("got subscribers", "subs", subs)

	// TODO: divide into batches and do it concurrently
	for _, sub := range subs {
		err := s.apis.Email.Send(ctx, SendOpt{
			To:      sub.Email,
			Subject: "Currency rate!",
			Body:    fmt.Sprintf("%f", rate),
		})
		if err != nil {
			logger.Error("failed to send email", "err", err)
		}
	}

	logger.Info("successfully sent emails")
	return nil
}

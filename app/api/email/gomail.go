package email

import (
	"context"
	"darkness8129/currency-converter/app/service"
	"darkness8129/currency-converter/config"
	"darkness8129/currency-converter/packages/logging"
	"fmt"

	"gopkg.in/gomail.v2"
)

var _ service.EmailAPI = (*gomailAPI)(nil)

type gomailAPI struct {
	logger *logging.Logger
	cfg    *config.Config
	dialer *gomail.Dialer
}

func NewGomailAPIAPI(logger *logging.Logger, cfg *config.Config) *gomailAPI {
	dialer := gomail.NewDialer(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.Username, cfg.SMTP.Password)

	return &gomailAPI{
		logger: logger.Named("gomailAPI"),
		cfg:    cfg,
		dialer: dialer,
	}
}

func (a *gomailAPI) Send(ctx context.Context, opt service.SendOpt) error {
	logger := a.logger.Named("Send")

	if a.cfg.App.TestMode {
		logger.Info("successfully sent email")
		return nil
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", a.cfg.SMTP.Username)
	msg.SetHeader("To", opt.To)
	msg.SetHeader("Subject", opt.Subject)
	msg.SetBody("text/plain", opt.Body)

	err := a.dialer.DialAndSend(msg)
	if err != nil {
		logger.Error("failed to send email", "err", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	logger.Info("successfully sent email")
	return nil
}

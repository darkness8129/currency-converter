package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App
		HTTP
		OXR
		PostgreSQL
		SMTP
	}

	App struct {
		TestMode      bool          `env:"APP_TEST_MODE" env-default:"true"`
		MailingPeriod time.Duration `env:"APP_MAILING_PERIOD" env-default:"24h"`
	}

	HTTP struct {
		Addr            string        `env:"HTTP_ADDR" env-default:":8080"`
		WriteTimeout    time.Duration `env:"HTTP_WRITE_TIMEOUT" env-default:"5s"`
		ReadTimeout     time.Duration `env:"HTTP_READ_TIMEOUT" env-default:"5s"`
		ShutdownTimeout time.Duration `env:"HTTP_SHUTDOWN_TIMEOUT" env-default:"3s"`
	}

	OXR struct {
		BaseURL string `env:"OXR_BASE_URL" env-default:"https://openexchangerates.org/api"`
		AppID   string `env:"OXR_APP_ID" env-default:"c4f1123da1fd4e4990e6ef32608c059a"`
	}

	PostgreSQL struct {
		User     string `env:"POSTGRESQL_USER" env-default:"postgres"`
		Password string `env:"POSTGRESQL_PASSWORD" env-default:"postgres"`
		Host     string `env:"POSTGRESQL_HOST" env-default:"localhost"`
		Database string `env:"POSTGRESQL_DATABASE" env-default:"subscription_service_db"`
		Port     string `env:"POSTGRESQL_PORT" env-default:"5432"`
	}

	SMTP struct {
		Host     string `env:"SMTP_HOST" env-default:"smtp.gmail.com"`
		Port     int    `env:"SMTP_PORT" env-default:"587"`
		Username string `env:"SMTP_USERNAME" env-default:""`
		Password string `env:"SMTP_PASSWORD" env-default:""`
	}
)

func New() (*Config, error) {
	var cfg Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to read env: %w", err)
	}

	return &cfg, nil
}

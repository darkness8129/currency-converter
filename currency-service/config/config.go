package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		HTTP
		OXR
	}

	HTTP struct {
		Addr            string        `env:"HTTP_ADDR" env-default:":8081"`
		WriteTimeout    time.Duration `env:"HTTP_WRITE_TIMEOUT" env-default:"5s"`
		ReadTimeout     time.Duration `env:"HTTP_READ_TIMEOUT" env-default:"5s"`
		ShutdownTimeout time.Duration `env:"HTTP_SHUTDOWN_TIMEOUT" env-default:"3s"`
	}

	OXR struct {
		BaseURL string `env:"OXR_BASE_URL" env-default:"https://openexchangerates.org/api"`
		AppID   string `env:"OXR_APP_ID" env-default:"c4f1123da1fd4e4990e6ef32608c059a"`
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

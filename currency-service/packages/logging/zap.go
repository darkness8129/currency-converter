package logging

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	logger *zap.SugaredLogger
}

func NewZapLogger() (*Logger, error) {
	config := zap.Config{
		Encoding:         "console",
		Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeName:     zapcore.FullNameEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			LevelKey:       "severity",
			CallerKey:      "caller",
			TimeKey:        "timestamp",
			NameKey:        "name",
			MessageKey:     "message",
			LineEnding:     "\n",
		},
	}

	logger, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to construct logger: %w", err)
	}

	return &Logger{logger.Sugar()}, nil
}

func (l *Logger) Named(name string) *Logger {
	return &Logger{l.logger.Named(name)}
}

func (l *Logger) Debug(message string, args ...interface{}) {
	l.logger.Debugw(message, args...)
}

func (l *Logger) Info(message string, args ...interface{}) {
	l.logger.Infow(message, args...)
}

func (l *Logger) Error(message string, args ...interface{}) {
	l.logger.Errorw(message, args...)
}

func (l *Logger) Fatal(message string, args ...interface{}) {
	l.logger.Fatalw(message, args...)
}

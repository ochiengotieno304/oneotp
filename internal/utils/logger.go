package utils

import (
	"context"

	config "github.com/ochiengotieno304/oneotp/internal/configs"
	"go.uber.org/zap"
)

type Fields map[string]interface{}

type Logger interface {
	Info(ctx context.Context, msg string, fields Fields)
	Error(ctx context.Context, msg string, fields Fields)
	Debug(ctx context.Context, msg string, fields Fields)
}

type logger struct {
	logger *zap.Logger
}

func (l *logger) Info(ctx context.Context, msg string, fields Fields) {
	l.logger.Info(msg, l.getFields(ctx, fields)...)
}

func (l *logger) Error(ctx context.Context, msg string, fields Fields) {
	l.logger.Error(msg, l.getFields(ctx, fields)...)
}

func (l *logger) Debug(ctx context.Context, msg string, fields Fields) {
	l.logger.Debug(msg, l.getFields(ctx, fields)...)
}

func (l *logger) getFields(ctx context.Context, fields Fields) []zap.Field {
	zapFields := []zap.Field{}

	if fields == nil {
		fields = Fields{}
	}

	if traceID := ctx.Value("trace_id"); traceID != nil {
		fields["trace_id"] = traceID
	}
	if url := ctx.Value("url"); url != nil {
		fields["url"] = url
	}

	for key, value := range fields {
		// Exclude sensitive fields from logging
		switch key {
		case "authorization", "secret", "api_key":
			continue
		}

		zapFields = append(zapFields, zap.Any(key, value))
	}

	return zapFields
}

func InitLogger() Logger {
	configs, err := config.LoadConfig()

	if err != nil {
		panic(err)
	}

	var zapLogger *zap.Logger

	if configs.Environment == "development" {
		zapLogger, _ = zap.NewDevelopment()
	} else {
		zapLogger, _ = zap.NewProduction()
	}

	return &logger{
		logger: zapLogger,
	}
}

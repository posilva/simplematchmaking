// Package logging implements a logger interface
package logging

import (
	"context"
	"os"

	"github.com/posilva/simplematchmaking/internal/core/ports"
	"github.com/rs/zerolog"
)

// NewSimpleLogger creates a new simpler logger
func NewSimpleLogger() *SimpleLogger {
	return &SimpleLogger{
		logger: zerolog.New(os.Stderr).With().Timestamp().Logger(),
	}
}

// SimpleLogger implements a simple logger interface
type SimpleLogger struct {
	logger zerolog.Logger
}

// Debug logs a debug message
func (log *SimpleLogger) Debug(msg string, keyvals ...interface{}) error {
	log.logger.Debug().Msgf(msg, keyvals...)
	return nil
}

// Info logs an info message
func (log *SimpleLogger) Info(msg string, keyvals ...interface{}) error {
	log.logger.Info().Msgf(msg, keyvals...)

	return nil
}

// Error logs an error message
func (log *SimpleLogger) Error(msg string, keyvals ...interface{}) error {
	log.logger.Error().Msgf(msg, keyvals...)
	return nil
}

// FromContext returns a logger from the context
func FromContext(ctx context.Context) ports.Logger {
	v := ctx.Value("logger")
	if v == nil {
		return NewSimpleLogger()
	}
	return v.(ports.Logger)

}

// ToContext adds a logger to the context
func ToContext(ctx context.Context, logger ports.Logger) context.Context {
	return context.WithValue(ctx, "logger", logger)
}

// WithContext returns a context with a logger
func WithContext(ctx context.Context, logger ports.Logger) context.Context {
	return ToContext(ctx, logger)
}

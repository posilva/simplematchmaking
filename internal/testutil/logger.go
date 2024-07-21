package testutil

import (
	"testing"

	"github.com/posilva/simplematchmaking/internal/core/ports"
)

// TestLogger is a logger implementation for testing
type TestLogger struct {
	t *testing.T
}

// Debug logs
func (l *TestLogger) Debug(msg string, v ...interface{}) error {
	l.t.Logf("[Test DEBUG]: "+msg, v...)
	return nil
}

// Info logs
func (l *TestLogger) Info(msg string, v ...interface{}) error {
	l.t.Logf("[Test INFO]: "+msg, v...)
	return nil
}

// Error logs
func (l *TestLogger) Error(msg string, v ...interface{}) error {
	l.t.Logf("[Test ERROR]: "+msg, v...)
	return nil

}

// NewLogger creates a new logger for testing
func NewLogger(t *testing.T) ports.Logger {
	return &TestLogger{
		t: t,
	}
}

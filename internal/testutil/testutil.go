// Package testutil is used to share test utilities
package testutil

import (
	"strings"
	"testing"

	uuid "github.com/segmentio/ksuid"
)

// Name returns the name of the test
func Name(t *testing.T) string {
	return t.Name()
}

// NewID returns an ID for tests using kuid package
func NewID() string {
	return strings.ToLower(uuid.New().String())
}

// NewUnique appends to a string a UUID to allow for uniqueness
func NewUnique(prefix string) string {
	return strings.ToLower(prefix + NewID())
}

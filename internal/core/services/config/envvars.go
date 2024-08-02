// Package config provides the configuration for the matchmaking service
package config

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/posilva/simplematchmaking/internal/core/domain"
)

const (
	// EnvVarName is the name of the environment variable that holds the configuration
	EnvVarName = "MATCHMAKING_CFG"
)

// EnvVar represents the configs loaded from environment variable for the matchmaking service
type EnvVar struct {
	config domain.MatchmakingConfig
}

// NewEnvVar creates a new EnvVars instance
func NewEnvVar() *EnvVar {

	return &EnvVar{}
}

// Load loads the configs from environment variables
func (e *EnvVar) Load() error {
	cfg, ok := os.LookupEnv(EnvVarName)
	if !ok {
		return fmt.Errorf("%s environment variable not set", EnvVarName)
	}

	typeAndData := strings.Split(cfg, ".")
	if len(typeAndData) != 2 {
		return fmt.Errorf("invalid format")
	}

	switch typeAndData[0] {
	case "json":
		data, err := base64.StdEncoding.DecodeString(typeAndData[1])
		if err != nil {
			return fmt.Errorf("error decoding base64: %w", err)
		}
		err = e.loadFromJSON(data)
		if err != nil {
			return fmt.Errorf("error loading configuration from JSON: %w", err)
		}
	default:
		return fmt.Errorf("invalid type: %s ", typeAndData[0])
	}
	return nil
}

// Get returns the loaded configuration
func (e *EnvVar) Get() domain.MatchmakingConfig {
	return e.config
}

func (e *EnvVar) loadFromJSON(data []byte) error {
	// TODO: Implement JSON validation
	e.config = domain.MatchmakingConfig{}
	err := json.Unmarshal(data, &e.config)
	if err != nil {
		return fmt.Errorf("Error unmarshalling JSON: %w", err)
	}
	return nil
}

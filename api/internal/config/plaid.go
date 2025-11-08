package config

import (
	"errors"
	"os"
)

// PlaidConfig holds Plaid API configuration
type PlaidConfig struct {
	ClientID    string `json:"client_id"`
	Secret      string `json:"secret"`
	Environment string `json:"environment"` // sandbox, development, production
}

// Validate ensures PlaidConfig has valid values
func (p *PlaidConfig) Validate() error {
	if p == nil {
		return errors.New("plaid config is nil")
	}

	// Allow environment variable overrides
	clientID := p.ClientID
	if envID := os.Getenv("PLAID_CLIENT_ID"); envID != "" {
		clientID = envID
	}

	secret := p.Secret
	if envSecret := os.Getenv("PLAID_SECRET"); envSecret != "" {
		secret = envSecret
	}

	environment := p.Environment
	if envEnv := os.Getenv("PLAID_ENV"); envEnv != "" {
		environment = envEnv
	}

	if clientID == "" {
		return errors.New("plaid client_id is required")
	}

	if secret == "" {
		return errors.New("plaid secret is required")
	}

	if environment == "" {
		return errors.New("plaid environment is required")
	}

	// Validate environment value
	switch environment {
	case "sandbox", "development", "production":
		// valid
	default:
		return errors.New("plaid environment must be 'sandbox', 'development', or 'production'")
	}

	return nil
}

// GetClientID returns the client ID with environment variable override
func (p *PlaidConfig) GetClientID() string {
	if envID := os.Getenv("PLAID_CLIENT_ID"); envID != "" {
		return envID
	}
	return p.ClientID
}

// GetSecret returns the secret with environment variable override
func (p *PlaidConfig) GetSecret() string {
	if envSecret := os.Getenv("PLAID_SECRET"); envSecret != "" {
		return envSecret
	}
	return p.Secret
}

// GetEnvironment returns the environment with environment variable override
func (p *PlaidConfig) GetEnvironment() string {
	if envEnv := os.Getenv("PLAID_ENV"); envEnv != "" {
		return envEnv
	}
	return p.Environment
}

package config

import "errors"

// CORS represents the CORS configuration
type CORS struct {
	AllowedOrigins []string `json:"allowed_origins"`
}

// Validate validates the CORS configuration
func (c *CORS) Validate() error {
	if c == nil {
		return errors.New("CORS config is nil")
	}

	if len(c.AllowedOrigins) == 0 {
		return errors.New("at least one allowed origin must be specified")
	}

	return nil
}

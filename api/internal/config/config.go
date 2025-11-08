package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sync"
)

// Version is the application version
// Injected at build time
var Version = "local"

const (
	defaultConfigDir = "/etc/lib/regulation/"
	configFileName   = "config.json"
	filePermissions  = 0644
)

// Config represents the application configuration
type Config struct {
	DB      *DB          `json:"database"`
	Redis   *Redis       `json:"redis"`
	Plaid   *PlaidConfig `json:"plaid"`
	OpenAI  *OpenAI      `json:"openai"`
	WebPush *WebPush     `json:"webpush,omitempty"`

	Debug bool `json:"debug"`

	mu sync.RWMutex
}

// getConfigDir returns the configuration directory path
func getConfigDir() string {
	if dir := os.Getenv("CONFIG_DIR"); dir != "" {
		return dir
	}
	return defaultConfigDir
}

// getConfigPath returns the full path to the config file
func getConfigPath() string {
	return filepath.Join(getConfigDir(), configFileName)
}

// Validate ensures all configuration values are valid
func (c *Config) Validate() error {
	if c == nil {
		return errors.New("config is nil")
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	val := reflect.ValueOf(c).Elem()
	typ := val.Type()

	for i := range val.NumField() {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Skip non-pointer fields and mutex
		if fieldType.Name == "mu" || field.Kind() != reflect.Ptr {
			continue
		}

		// Check if required pointer field is nil
		if field.IsNil() {
			return fmt.Errorf("missing %s configuration", fieldType.Name)
		}

		// If the field implements Validate() method, call it
		if validator, ok := field.Interface().(ValidatableConfig); ok {
			if err := validator.Validate(); err != nil {
				return fmt.Errorf("invalid %s configuration: %w", fieldType.Name, err)
			}
		}
	}

	return nil
}

// Save writes the configuration to disk
func (c *Config) Save() error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	configFile, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Ensure directory exists
	if err = os.MkdirAll(getConfigDir(), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if err = os.WriteFile(getConfigPath(), configFile, filePermissions); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// LoadConfig reads and parses the configuration file
func LoadConfig() (*Config, error) {
	configFile, err := os.ReadFile(getConfigPath())
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err = json.Unmarshal(configFile, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	if config.Plaid == nil {
		config.Plaid = &PlaidConfig{}
	}

	if err = config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &config, nil
}

type ValidatableConfig interface {
	Validate() error
}

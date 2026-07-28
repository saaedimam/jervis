package config

import (
	"fmt"

	"github.com/ioriimasu/jervis/internal/runtime/contracts"
	"github.com/ioriimasu/jervis/internal/runtime/errors"
)

var (
	_ contracts.Validator = (*Config)(nil)
	_ contracts.Freezable = (*Config)(nil)
)

// Config represents immutable runtime configuration.
type Config struct {
	environment string
	logLevel    string
	dataDir     string
	maxSessions int
	frozen      bool
}

// Default returns a new default Config.
func Default() Config {
	return Config{
		environment: "development",
		logLevel:    "info",
		dataDir:     "~/.jervis",
		maxSessions: 100,
		frozen:      false,
	}
}

// New creates a validated Config instance.
func New(env, logLevel, dataDir string, maxSessions int) (Config, error) {
	cfg := Config{
		environment: env,
		logLevel:    logLevel,
		dataDir:     dataDir,
		maxSessions: maxSessions,
		frozen:      false,
	}
	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

// Environment returns the target environment string.
func (c Config) Environment() string {
	return c.environment
}

// LogLevel returns the system log level.
func (c Config) LogLevel() string {
	return c.logLevel
}

// DataDir returns the path to the application data directory.
func (c Config) DataDir() string {
	return c.dataDir
}

// MaxSessions returns the maximum allowed concurrent sessions limit.
func (c Config) MaxSessions() int {
	return c.maxSessions
}

// IsFrozen reports whether the Config instance is frozen and immutable.
func (c Config) IsFrozen() bool {
	return c.frozen
}

// Validate checks that all configuration fields comply with invariants.
func (c Config) Validate() error {
	if c.environment == "" {
		return fmt.Errorf("%w: environment cannot be empty", errors.ErrConfiguration)
	}
	if c.logLevel == "" {
		return fmt.Errorf("%w: logLevel cannot be empty", errors.ErrConfiguration)
	}
	if c.dataDir == "" {
		return fmt.Errorf("%w: dataDir cannot be empty", errors.ErrConfiguration)
	}
	if c.maxSessions <= 0 {
		return fmt.Errorf("%w: maxSessions must be greater than zero", errors.ErrConfiguration)
	}
	return nil
}

// Clone returns a deep copy of the Config instance.
func (c Config) Clone() Config {
	return Config{
		environment: c.environment,
		logLevel:    c.logLevel,
		dataDir:     c.dataDir,
		maxSessions: c.maxSessions,
		frozen:      c.frozen,
	}
}

// Freeze returns a new Config instance marked as frozen and immutable.
func (c Config) Freeze() Config {
	cloned := c.Clone()
	cloned.frozen = true
	return cloned
}

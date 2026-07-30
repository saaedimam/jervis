package config_test

import (
	"errors"
	"testing"

	"github.com/saaedimam/jervis/internal/runtime/config"
	errs "github.com/saaedimam/jervis/internal/runtime/errors"
)

func TestConfigDefault(t *testing.T) {
	cfg := config.Default()

	if cfg.Environment() != "development" {
		t.Errorf("expected environment 'development', got %s", cfg.Environment())
	}
	if cfg.LogLevel() != "info" {
		t.Errorf("expected logLevel 'info', got %s", cfg.LogLevel())
	}
	if cfg.DataDir() != "~/.jervis" {
		t.Errorf("expected dataDir '~/.jervis', got %s", cfg.DataDir())
	}
	if cfg.MaxSessions() != 100 {
		t.Errorf("expected maxSessions 100, got %d", cfg.MaxSessions())
	}
	if cfg.IsFrozen() {
		t.Errorf("expected IsFrozen false for Default()")
	}

	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected Default() config to be valid: %v", err)
	}
}

func TestConfigNewAndValidation(t *testing.T) {
	tests := []struct {
		name        string
		env         string
		logLevel    string
		dataDir     string
		maxSessions int
		expectErr   bool
	}{
		{
			name:        "Valid Config",
			env:         "production",
			logLevel:    "warn",
			dataDir:     "/var/jervis",
			maxSessions: 50,
			expectErr:   false,
		},
		{
			name:        "Empty Environment",
			env:         "",
			logLevel:    "info",
			dataDir:     "/var/jervis",
			maxSessions: 50,
			expectErr:   true,
		},
		{
			name:        "Empty LogLevel",
			env:         "production",
			logLevel:    "",
			dataDir:     "/var/jervis",
			maxSessions: 50,
			expectErr:   true,
		},
		{
			name:        "Empty DataDir",
			env:         "production",
			logLevel:    "info",
			dataDir:     "",
			maxSessions: 50,
			expectErr:   true,
		},
		{
			name:        "Zero MaxSessions",
			env:         "production",
			logLevel:    "info",
			dataDir:     "/var/jervis",
			maxSessions: 0,
			expectErr:   true,
		},
		{
			name:        "Negative MaxSessions",
			env:         "production",
			logLevel:    "info",
			dataDir:     "/var/jervis",
			maxSessions: -10,
			expectErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := config.New(tt.env, tt.logLevel, tt.dataDir, tt.maxSessions)
			if tt.expectErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if !errors.Is(err, errs.ErrConfiguration) {
					t.Fatalf("expected ErrConfiguration wrapper, got %v", err)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if cfg.Environment() != tt.env || cfg.LogLevel() != tt.logLevel || cfg.DataDir() != tt.dataDir || cfg.MaxSessions() != tt.maxSessions {
					t.Fatalf("config field mismatch")
				}
			}
		})
	}
}

func TestConfigCloneAndFreeze(t *testing.T) {
	cfg := config.Default()
	cloned := cfg.Clone()

	if cloned != cfg {
		t.Fatalf("expected cloned config to equal original config")
	}

	frozen := cfg.Freeze()
	if !frozen.IsFrozen() {
		t.Fatalf("expected Freeze() to set frozen=true")
	}
	if cfg.IsFrozen() {
		t.Fatalf("expected original config to remain unfrozen")
	}
}

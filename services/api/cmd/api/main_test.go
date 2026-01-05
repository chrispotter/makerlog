package main

import (
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		fallback string
		envValue string
		setEnv   bool
		expected string
	}{
		{
			name:     "environment variable exists",
			key:      "TEST_VAR_EXISTS",
			fallback: "fallback-value",
			envValue: "actual-value",
			setEnv:   true,
			expected: "actual-value",
		},
		{
			name:     "environment variable does not exist",
			key:      "TEST_VAR_NOT_EXISTS",
			fallback: "fallback-value",
			envValue: "",
			setEnv:   false,
			expected: "fallback-value",
		},
		{
			name:     "environment variable is empty string",
			key:      "TEST_VAR_EMPTY",
			fallback: "fallback-value",
			envValue: "",
			setEnv:   true,
			expected: "fallback-value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			if tt.setEnv {
				if err := os.Setenv(tt.key, tt.envValue); err != nil {
					t.Fatalf("Failed to set env var: %v", err)
				}
				defer func() {
					if err := os.Unsetenv(tt.key); err != nil {
						t.Errorf("Failed to unset env var: %v", err)
					}
				}()
			}

			// Execute
			result := getEnv(tt.key, tt.fallback)

			// Verify
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

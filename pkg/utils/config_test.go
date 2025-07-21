package utils

import (
	"os"
	"testing"
)

func TestGetEnvWithDefault(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		setEnv       bool
		expected     string
	}{
		{
			name:         "returns environment variable when set",
			key:          "TEST_VAR_SET",
			defaultValue: "default",
			envValue:     "env_value",
			setEnv:       true,
			expected:     "env_value",
		},
		{
			name:         "returns default when environment variable not set",
			key:          "TEST_VAR_NOT_SET",
			defaultValue: "default_value",
			envValue:     "",
			setEnv:       false,
			expected:     "default_value",
		},
		{
			name:         "returns environment variable when set to empty string explicitly",
			key:          "TEST_VAR_EMPTY",
			defaultValue: "default",
			envValue:     "",
			setEnv:       true,
			expected:     "default",
		},
		{
			name:         "handles empty default value",
			key:          "TEST_VAR_EMPTY_DEFAULT",
			defaultValue: "",
			envValue:     "",
			setEnv:       false,
			expected:     "",
		},
		{
			name:         "handles special characters in environment variable",
			key:          "TEST_VAR_SPECIAL",
			defaultValue: "default",
			envValue:     "value with spaces and symbols !@#$%^&*()",
			setEnv:       true,
			expected:     "value with spaces and symbols !@#$%^&*()",
		},
		{
			name:         "handles multiline environment variable",
			key:          "TEST_VAR_MULTILINE",
			defaultValue: "default",
			envValue:     "line1\nline2\nline3",
			setEnv:       true,
			expected:     "line1\nline2\nline3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalValue := os.Getenv(tt.key)
			defer func() {
				if originalValue != "" {
					if err := os.Setenv(tt.key, originalValue); err != nil {
						t.Logf("Failed to restore environment variable: %v", err)
					}
				} else {
					if err := os.Unsetenv(tt.key); err != nil {
						t.Logf("Failed to unset environment variable: %v", err)
					}
				}
			}()

			if tt.setEnv {
				if err := os.Setenv(tt.key, tt.envValue); err != nil {
					t.Fatalf("Failed to set environment variable: %v", err)
				}
			} else {
				if err := os.Unsetenv(tt.key); err != nil {
					t.Logf("Failed to unset environment variable: %v", err)
				}
			}

			result := GetEnvWithDefault(tt.key, tt.defaultValue)

			if result != tt.expected {
				t.Errorf("GetEnvWithDefault(%q, %q) = %q, want %q", tt.key, tt.defaultValue, result, tt.expected)
			}
		})
	}
}

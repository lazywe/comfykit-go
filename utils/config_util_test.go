package utils

import (
	"os"
	"testing"
)

func TestGetEnvOrDefault(t *testing.T) {
	oldVal := os.Getenv("TEST_ENV_VAR_NOT_SET")
	defer os.Setenv("TEST_ENV_VAR_NOT_SET", oldVal)

	os.Unsetenv("TEST_ENV_VAR_NOT_SET")

	result := GetEnvOrDefault("TEST_ENV_VAR_NOT_SET", "default_value")
	if result != "default_value" {
		t.Errorf("GetEnvOrDefault() = %v, want default_value", result)
	}

	os.Setenv("TEST_ENV_VAR_NOT_SET", "set_value")
	result = GetEnvOrDefault("TEST_ENV_VAR_NOT_SET", "default_value")
	if result != "set_value" {
		t.Errorf("GetEnvOrDefault() = %v, want set_value", result)
	}
}

func TestGetEnvIntOrDefault(t *testing.T) {
	oldVal := os.Getenv("TEST_ENV_INT")
	defer os.Setenv("TEST_ENV_INT", oldVal)

	os.Unsetenv("TEST_ENV_INT")
	result := GetEnvIntOrDefault("TEST_ENV_INT", 42)
	if result != 42 {
		t.Errorf("GetEnvIntOrDefault() = %v, want 42", result)
	}

	os.Setenv("TEST_ENV_INT", "100")
	result = GetEnvIntOrDefault("TEST_ENV_INT", 42)
	if result != 100 {
		t.Errorf("GetEnvIntOrDefault() = %v, want 100", result)
	}

	os.Setenv("TEST_ENV_INT", "not_a_number")
	result = GetEnvIntOrDefault("TEST_ENV_INT", 42)
	if result != 42 {
		t.Errorf("GetEnvIntOrDefault() with invalid value = %v, want 42", result)
	}
}

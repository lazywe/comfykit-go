package utils

import (
	"testing"
)

func TestGetNestedValue(t *testing.T) {
	data := map[string]interface{}{
		"level1": map[string]interface{}{
			"level2": map[string]interface{}{
				"value": "target",
			},
		},
	}

	tests := []struct {
		name     string
		path     string
		expected interface{}
	}{
		{"simple path", "level1.level2.value", "target"},
		{"non-existent path", "nonexistent", nil},
		{"empty path", "", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetNestedValue(data, tt.path)
			if result != tt.expected {
				t.Errorf("GetNestedValue(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}

	t.Run("single level", func(t *testing.T) {
		result := GetNestedValue(data, "level1")
		if result == nil {
			t.Error("GetNestedValue(level1) should not be nil")
		}
		_, ok := result.(map[string]interface{})
		if !ok {
			t.Error("GetNestedValue(level1) should return a map")
		}
	})
}

func TestSetNestedValue(t *testing.T) {
	data := map[string]interface{}{
		"level1": map[string]interface{}{
			"level2": "old_value",
		},
	}

	SetNestedValue(data, "level1.level2", "new_value")

	result := data["level1"].(map[string]interface{})["level2"]
	if result != "new_value" {
		t.Errorf("SetNestedValue() = %v, want new_value", result)
	}
}

func TestToString(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"string", "hello", "hello"},
		{"int", int64(42), "42"},
		{"float", float64(3.14), "3.14"},
		{"nil", nil, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToString(tt.input)
			if result != tt.expected {
				t.Errorf("ToString(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToInt(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected int
		hasError bool
	}{
		{"int", int(42), 42, false},
		{"int64", int64(42), 42, false},
		{"float64", float64(42.5), 42, false},
		{"string", "123", 123, false},
		{"invalid", "not_a_number", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ToInt(tt.input)
			if tt.hasError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.hasError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("ToInt(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToBool(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected bool
		hasError bool
	}{
		{"bool true", true, true, false},
		{"bool false", false, false, false},
		{"string true", "true", true, false},
		{"string 1", "1", true, false},
		{"string false", "false", false, false},
		{"int 1", 1, true, false},
		{"int 0", 0, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ToBool(tt.input)
			if tt.hasError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.hasError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("ToBool(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

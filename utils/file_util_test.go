package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileExists(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_*.txt")
	if err != nil {
		t.Fatal(err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)
	tmpFile.Close()

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{"existing file", tmpPath, true},
		{"non-existing file", "/nonexistent/path/to/file.txt", false},
		{"empty path", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FileExists(tt.path)
			if result != tt.expected {
				t.Errorf("FileExists(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestReadJSONFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_*.json")
	if err != nil {
		t.Fatal(err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	testData := map[string]interface{}{"key": "value", "number": float64(42)}
	err = WriteJSONFile(tmpPath, testData)
	if err != nil {
		t.Fatal(err)
	}

	var data map[string]interface{}
	err = ReadJSONFile(tmpPath, &data)
	if err != nil {
		t.Errorf("ReadJSONFile() error = %v", err)
	}
	if data["key"] != "value" {
		t.Errorf("data[key] = %v, want value", data["key"])
	}
	if data["number"] != float64(42) {
		t.Errorf("data[number] = %v, want 42", data["number"])
	}

	invalidFile, _ := os.CreateTemp("", "invalid_*.json")
	invalidFile.WriteString("not valid json")
	invalidFile.Close()
	defer os.Remove(invalidFile.Name())

	var invalidData map[string]interface{}
	err = ReadJSONFile(invalidFile.Name(), &invalidData)
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestCreateDirIfNotExists(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "test_dir")
	defer os.RemoveAll(tmpDir)

	err := CreateDirIfNotExists(tmpDir)
	if err != nil {
		t.Errorf("CreateDirIfNotExists() error = %v", err)
	}

	info, err := os.Stat(tmpDir)
	if err != nil {
		t.Errorf("Directory was not created: %v", err)
	}
	if !info.IsDir() {
		t.Error("Created path is not a directory")
	}

	err = CreateDirIfNotExists(tmpDir)
	if err != nil {
		t.Errorf("CreateDirIfNotExists() should be idempotent: %v", err)
	}
}

func TestGetFileName(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{"full path", "/path/to/file.txt", "file.txt"},
		{"relative path", "workflow.json", "workflow.json"},
		{"empty", "", "."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetFileName(tt.path)
			if result != tt.expected {
				t.Errorf("GetFileName(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestGetFileNameWithoutExt(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{"full path", "/path/to/file.txt", "file"},
		{"relative path", "workflow.json", "workflow"},
		{"no extension", "readme", "readme"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetFileNameWithoutExt(tt.path)
			if result != tt.expected {
				t.Errorf("GetFileNameWithoutExt(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestGetFileExtension(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{"txt", "/path/to/file.txt", ".txt"},
		{"json", "workflow.json", ".json"},
		{"no extension", "readme", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetFileExtension(tt.path)
			if result != tt.expected {
				t.Errorf("GetFileExtension(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

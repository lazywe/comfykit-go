package comfykit

import (
	"os"
	"testing"
)

func TestNewComfyKit(t *testing.T) {
	tests := []struct {
		name         string
		opts         []ComfyKitOption
		wantURL      string
		wantExecType ExecutorType
	}{
		{
			name:         "default configuration",
			opts:         nil,
			wantURL:      "http://127.0.0.1:8188",
			wantExecType: ExecutorTypeHTTP,
		},
		{
			name: "custom configuration",
			opts: []ComfyKitOption{
				WithComfyUIBaseURL("http://custom:8188"),
				WithExecutorType(ExecutorTypeWebSocket),
			},
			wantURL:      "http://custom:8188",
			wantExecType: ExecutorTypeWebSocket,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kit := NewComfyKit(tt.opts...)
			if kit.GetComfyUIBaseURL() != tt.wantURL {
				t.Errorf("GetComfyUIBaseURL() = %v, want %v", kit.GetComfyUIBaseURL(), tt.wantURL)
			}
			if kit.GetExecutorType() != tt.wantExecType {
				t.Errorf("GetExecutorType() = %v, want %v", kit.GetExecutorType(), tt.wantExecType)
			}
		})
	}
}

func TestNewComfyKitWithEnvVars(t *testing.T) {
	// Save original env vars
	oldURL := os.Getenv("COMFYUI_BASE_URL")
	oldExecType := os.Getenv("COMFYUI_EXECUTOR_TYPE")
	oldAPIKey := os.Getenv("COMFYUI_API_KEY")

	// Cleanup
	defer func() {
		os.Setenv("COMFYUI_BASE_URL", oldURL)
		os.Setenv("COMFYUI_EXECUTOR_TYPE", oldExecType)
		os.Setenv("COMFYUI_API_KEY", oldAPIKey)
	}()

	// Set env vars
	os.Setenv("COMFYUI_BASE_URL", "http://env-server:8188")
	os.Setenv("COMFYUI_EXECUTOR_TYPE", "websocket")
	os.Setenv("COMFYUI_API_KEY", "env-api-key")

	kit := NewComfyKit()
	if kit.GetComfyUIBaseURL() != "http://env-server:8188" {
		t.Errorf("GetComfyUIBaseURL() = %v, want http://env-server:8188", kit.GetComfyUIBaseURL())
	}
	if kit.GetExecutorType() != ExecutorTypeWebSocket {
		t.Errorf("GetExecutorType() = %v, want websocket", kit.GetExecutorType())
	}
	if kit.GetAPIKey() != "env-api-key" {
		t.Errorf("GetAPIKey() = %v, want env-api-key", kit.GetAPIKey())
	}
}

func TestConfigurationPriority(t *testing.T) {
	// Set env var
	oldURL := os.Getenv("COMFYUI_BASE_URL")
	defer os.Setenv("COMFYUI_BASE_URL", oldURL)
	os.Setenv("COMFYUI_BASE_URL", "http://env:8188")

	// Param should override env
	kit := NewComfyKit(WithComfyUIBaseURL("http://param:8188"))
	if kit.GetComfyUIBaseURL() != "http://param:8188" {
		t.Errorf("Param should override env, got %v", kit.GetComfyUIBaseURL())
	}
}

func TestRunningHubConfig(t *testing.T) {
	kit := NewComfyKit(
		WithRunningHubAPIKey("test-key"),
		WithRunningHubTimeout(600),
	)

	if kit.GetRunningHubAPIKey() != "test-key" {
		t.Errorf("GetRunningHubAPIKey() = %v, want test-key", kit.GetRunningHubAPIKey())
	}
	if kit.GetRunningHubTimeout() != 600 {
		t.Errorf("GetRunningHubTimeout() = %v, want 600", kit.GetRunningHubTimeout())
	}
}

func TestIsRunninghubWorkflowId(t *testing.T) {
	kit := NewComfyKit()

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"pure numeric", "12345", true},
		{"empty", "", false},
		{"with letters", "abc123", false},
		{"with special chars", "123-456", false},
		{"large number", "1234567890", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := kit.isRunninghubWorkflowId(tt.input)
			if result != tt.expected {
				t.Errorf("isRunninghubWorkflowId(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsUrl(t *testing.T) {
	kit := NewComfyKit()

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"http url", "http://example.com", true},
		{"https url", "https://example.com", true},
		{"local file", "/path/to/file.json", false},
		{"relative path", "workflow.json", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := kit.isUrl(tt.input)
			if result != tt.expected {
				t.Errorf("isUrl(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestClose(t *testing.T) {
	kit := NewComfyKit()
	err := kit.Close()
	if err != nil {
		t.Errorf("Close() error = %v", err)
	}
}

package comfyui

import (
	"testing"
)

func TestWorkflowParser_ParseTitle(t *testing.T) {
	parser := NewWorkflowParser()

	tests := []struct {
		name     string
		title    string
		expected []string
	}{
		{
			name:     "basic param",
			title:    "Prompt, $prompt",
			expected: []string{"$prompt"},
		},
		{
			name:     "required param",
			title:    "Prompt, $prompt!",
			expected: []string{"$prompt!"},
		},
		{
			name:     "upload param",
			title:    "Image, $~image",
			expected: []string{"$~image"},
		},
		{
			name:     "multiple params",
			title:    "Prompt, $prompt!, $seed, $steps",
			expected: []string{"$prompt!", "$seed", "$steps"},
		},
		{
			name:     "no params",
			title:    "Simple Prompt",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.ParseTitle(tt.title)
			if len(result) != len(tt.expected) {
				t.Fatalf("ParseTitle(%q) returned %d params, want %d", tt.title, len(result), len(tt.expected))
			}
			for i, param := range result {
				if param != tt.expected[i] {
					t.Errorf("Param[%d] = %v, want %v", i, param, tt.expected[i])
				}
			}
		})
	}
}

func TestWorkflowParser_ParseParamMarker(t *testing.T) {
	parser := NewWorkflowParser()

	tests := []struct {
		name          string
		marker        string
		wantName      string
		wantRequired  bool
		wantNeedUpload bool
	}{
		{"basic", "$prompt", "prompt", false, false},
		{"required", "$prompt!", "prompt", true, false},
		{"upload", "$~image", "image", false, true},
		{"required upload", "$~image!", "image", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.ParseParamMarker(tt.marker)
			if result == nil {
				t.Fatal("ParseParamMarker returned nil")
			}
			if result["name"] != tt.wantName {
				t.Errorf("name = %v, want %v", result["name"], tt.wantName)
			}
			if result["required"] != tt.wantRequired {
				t.Errorf("required = %v, want %v", result["required"], tt.wantRequired)
			}
			if result["upload"] != tt.wantNeedUpload {
				t.Errorf("upload = %v, want %v", result["upload"], tt.wantNeedUpload)
			}
		})
	}
}

func TestWorkflowParser_InferTypeFromValue(t *testing.T) {
	parser := NewWorkflowParser()

	tests := []struct {
		name     string
		value    interface{}
		expected string
	}{
		{"bool", true, "bool"},
		{"int", int(42), "int"},
		{"float", float64(3.14), "float"},
		{"string", "hello", "str"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.InferTypeFromValue(tt.value)
			if result != tt.expected {
				t.Errorf("InferTypeFromValue(%v) = %v, want %v", tt.value, result, tt.expected)
			}
		})
	}
}

func TestWorkflowParser_IsKnownOutputNode(t *testing.T) {
	parser := NewWorkflowParser()

	tests := []struct {
		name     string
		classType string
		expected bool
	}{
		{"SaveImage", "SaveImage", true},
		{"SaveVideo", "SaveVideo", true},
		{"SaveAudio", "SaveAudio", true},
		{"VHS_SaveVideo", "VHS_SaveVideo", true},
		{"VHS_SaveAudio", "VHS_SaveAudio", true},
		{"CLIPTextEncode", "CLIPTextEncode", false},
		{"KSampler", "KSampler", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.IsKnownOutputNode(tt.classType)
			if result != tt.expected {
				t.Errorf("IsKnownOutputNode(%q) = %v, want %v", tt.classType, result, tt.expected)
			}
		})
	}
}

func TestIsRunningHubWorkflow(t *testing.T) {
	result := IsRunningHubWorkflow("nonexistent.json")
	if result != false {
		t.Errorf("IsRunningHubWorkflow() for non-existent file = %v, want false", result)
	}
}

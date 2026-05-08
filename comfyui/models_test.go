package comfyui

import (
	"testing"
)

func TestNewExecuteResult(t *testing.T) {
	result := NewExecuteResult("completed")
	if result.Status != "completed" {
		t.Errorf("Status = %v, want completed", result.Status)
	}
	if result.Images == nil {
		t.Error("Images should not be nil")
	}
	if result.ImagesByVar == nil {
		t.Error("ImagesByVar should not be nil")
	}
}

func TestExecuteResultErrorResult(t *testing.T) {
	result := NewExecuteResult("completed").ErrorResult("test error")
	if result.Status != "error" {
		t.Errorf("Status = %v, want error", result.Status)
	}
	if result.Message != "test error" {
		t.Errorf("Message = %v, want test error", result.Message)
	}
}

func TestExecuteResultTimeoutResult(t *testing.T) {
	result := NewExecuteResult("completed").TimeoutResult(30.5)
	if result.Status != "timeout" {
		t.Errorf("Status = %v, want timeout", result.Status)
	}
	if result.Duration != 30.5 {
		t.Errorf("Duration = %v, want 30.5", result.Duration)
	}
}

func TestWorkflowParam(t *testing.T) {
	param := WorkflowParam{
		Name:      "prompt",
		Type:      "str",
		Required:  true,
		Default:   "default value",
		NeedUpload: false,
	}

	if param.Name != "prompt" {
		t.Errorf("Name = %v, want prompt", param.Name)
	}
	if !param.Required {
		t.Error("Required should be true")
	}
}

func TestWorkflowMetadata(t *testing.T) {
	metadata := WorkflowMetadata{
		Title:    "Test Workflow",
		WorkflowID: "12345",
		IsRunningHub: true,
	}

	if metadata.Title != "Test Workflow" {
		t.Errorf("Title = %v, want Test Workflow", metadata.Title)
	}
	if metadata.IsRunningHub != true {
		t.Error("IsRunningHub should be true")
	}
}

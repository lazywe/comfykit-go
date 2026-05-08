package utils

import (
	"fmt"
	"regexp"
	"strconv"
)

func IsRunningHubWorkflowID(workflow string) bool {
	return regexp.MustCompile(`^\d+$`).MatchString(workflow)
}

func IsRunningHubWorkflowFile(filePath string) bool {
	var workflowJSON map[string]interface{}
	if err := ReadJSONFile(filePath, &workflowJSON); err != nil {
		return false
	}

	if _, ok := workflowJSON["_source"]; ok {
		source, _ := workflowJSON["_source"].(string)
		return source == "runninghub"
	}

	if metadata, ok := workflowJSON["__metadata__"]; ok {
		if metaMap, ok := metadata.(map[string]interface{}); ok {
			if source, ok := metaMap["source"].(string); ok {
				return source == "runninghub"
			}
		}
	}

	return false
}

func GetWorkflowIDFromFile(filePath string) (string, error) {
	var workflowJSON map[string]interface{}
	if err := ReadJSONFile(filePath, &workflowJSON); err != nil {
		return "", err
	}

	if workflowID, ok := workflowJSON["workflow_id"]; ok {
		return fmt.Sprintf("%v", workflowID), nil
	}

	if metadata, ok := workflowJSON["__metadata__"]; ok {
		if metaMap, ok := metadata.(map[string]interface{}); ok {
			if workflowID, ok := metaMap["workflow_id"].(string); ok {
				return workflowID, nil
			}
		}
	}

	return "", fmt.Errorf("workflow_id not found in workflow file")
}

func ParseTaskID(taskID string) (int64, error) {
	return strconv.ParseInt(taskID, 10, 64)
}

func FormatTaskID(taskID int64) string {
	return fmt.Sprintf("%d", taskID)
}

func ValidateAPIKey(apiKey string) bool {
	if apiKey == "" {
		return false
	}
	return len(apiKey) >= 16
}

func GetDefaultInstanceType() string {
	return "standard"
}

func GetInstanceTypes() []string {
	return []string{"standard", "plus", "pro", "gpu"}
}

func IsValidInstanceType(instanceType string) bool {
	for _, t := range GetInstanceTypes() {
		if t == instanceType {
			return true
		}
	}
	return false
}

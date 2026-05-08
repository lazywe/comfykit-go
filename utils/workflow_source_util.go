package utils

import (
	"fmt"
	"strings"
)

type WorkflowSourceType string

const (
	WorkflowSourceTypeFile      WorkflowSourceType = "file"
	WorkflowSourceTypeURL       WorkflowSourceType = "url"
	WorkflowSourceTypeRunningHub WorkflowSourceType = "runninghub"
)

func DetectWorkflowSource(workflow string) WorkflowSourceType {
	if IsRunningHubWorkflowID(workflow) {
		return WorkflowSourceTypeRunningHub
	}

	if IsURL(workflow) {
		return WorkflowSourceTypeURL
	}

	if FileExists(workflow) {
		return WorkflowSourceTypeFile
	}

	if strings.Contains(workflow, "/") || strings.Contains(workflow, "\\") {
		return WorkflowSourceTypeFile
	}

	return WorkflowSourceTypeFile
}

func ResolveWorkflowPath(workflow string) (string, error) {
	sourceType := DetectWorkflowSource(workflow)

	switch sourceType {
	case WorkflowSourceTypeURL:
		return DownloadFileToTemp(workflow)

	case WorkflowSourceTypeRunningHub:
		return "", fmt.Errorf("cannot resolve RunningHub workflow ID to file path")

	case WorkflowSourceTypeFile:
		if FileExists(workflow) {
			return workflow, nil
		}

		searchPaths := []string{
			workflow,
			"./workflows/" + workflow,
			"./workflows/" + workflow + ".json",
			"workflows/" + workflow,
			"workflows/" + workflow + ".json",
			GetHomeDir() + "/.comfykit/workflows/" + workflow,
			GetHomeDir() + "/.comfykit/workflows/" + workflow + ".json",
		}

		for _, path := range searchPaths {
			if FileExists(path) {
				return path, nil
			}
		}

		return "", fmt.Errorf("workflow file not found: %s", workflow)

	default:
		return "", fmt.Errorf("unknown workflow source type")
	}
}

func ListAvailableWorkflows() ([]string, error) {
	var workflows []string

	searchDirs := []string{
		"./workflows",
		"workflows",
		GetHomeDir() + "/.comfykit/workflows",
	}

	for _, dir := range searchDirs {
		if !DirExists(dir) {
			continue
		}

		files, err := ListDir(dir)
		if err != nil {
			continue
		}

		for _, file := range files {
			if strings.HasSuffix(file, ".json") {
				workflows = append(workflows, file[:len(file)-5])
			}
		}
	}

	return workflows, nil
}

func GetWorkflowMetadata(workflowPath string) (map[string]interface{}, error) {
	var workflowJSON map[string]interface{}
	if err := ReadJSONFile(workflowPath, &workflowJSON); err != nil {
		return nil, err
	}

	metadata := map[string]interface{}{}

	if meta, ok := workflowJSON["__metadata__"]; ok {
		if metaMap, ok := meta.(map[string]interface{}); ok {
			for k, v := range metaMap {
				metadata[k] = v
			}
		}
	}

	if title, ok := workflowJSON["title"]; ok {
		metadata["title"] = title
	}

	return metadata, nil
}

func SaveWorkflow(workflowJSON map[string]interface{}, filePath string) error {
	return WriteJSONFile(filePath, workflowJSON)
}

func CreateWorkflowBackup(workflowPath string) (string, error) {
	backupPath := workflowPath + ".bak"
	if err := CopyFile(workflowPath, backupPath); err != nil {
		return "", err
	}
	return backupPath, nil
}

func ValidateWorkflowFile(workflowPath string) error {
	if !FileExists(workflowPath) {
		return fmt.Errorf("workflow file does not exist: %s", workflowPath)
	}

	var workflowJSON map[string]interface{}
	if err := ReadJSONFile(workflowPath, &workflowJSON); err != nil {
		return fmt.Errorf("invalid JSON in workflow file: %w", err)
	}

	if _, ok := workflowJSON["nodes"]; !ok {
		if len(workflowJSON) == 0 {
			return fmt.Errorf("workflow file is empty")
		}
	}

	return nil
}

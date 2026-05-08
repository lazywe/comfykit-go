package comfyui

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type RunningHubExecutor struct {
	*ComfyUIExecutor
	timeout     int
	retryCount  int
	instanceType string
	client      *RunningHubClient
}

func NewRunningHubExecutor(baseURL, apiKey string, timeout, retryCount int, instanceType string) *RunningHubExecutor {
	return &RunningHubExecutor{
		ComfyUIExecutor: NewComfyUIExecutor(baseURL, apiKey, ""),
		timeout:         timeout,
		retryCount:      retryCount,
		instanceType:    instanceType,
		client:          NewRunningHubClient(apiKey, baseURL, timeout, retryCount, instanceType),
	}
}

func (e *RunningHubExecutor) Close() error {
	return e.client.Close()
}

func (e *RunningHubExecutor) ExecuteByID(workflowID string, params map[string]interface{}) (*ExecuteResult, error) {
	startTime := time.Now()

	workflowJSON, err := e.client.GetWorkflowJSON(workflowID)
	if err != nil {
		return NewExecuteResult("error").ErrorResult("Failed to get workflow JSON: "+err.Error()), nil
	}

	workflowJSON, seedChanges := randomizeSeedInWorkflow(workflowJSON)

	parser := NewWorkflowParser()
	metadata := parser.ParseWorkflow(workflowJSON, "workflow_"+workflowID)
	if metadata == nil {
		return NewExecuteResult("error").ErrorResult("Failed to parse workflow metadata"), nil
	}

	metadata.WorkflowID = workflowID
	metadata.IsRunningHub = true

	nodeInfoList := e.ConvertParamsToNodeInfoList(metadata, params, seedChanges)
	outputID2Var := extractOutputNodes(metadata)

	taskData, err := e.client.CreateTask(workflowID, nodeInfoList)
	if err != nil {
		return NewExecuteResult("error").ErrorResult("Failed to create task: "+err.Error()), nil
	}

	taskID, ok := taskData["taskId"].(string)
	if !ok || taskID == "" {
		return NewExecuteResult("error").ErrorResult("Failed to get task ID"), nil
	}

	result, err := e.WaitForTaskCompletion(taskID, outputID2Var)
	if err != nil {
		return result, err
	}

	result.Duration = time.Since(startTime).Seconds()
	return result, nil
}

func (e *RunningHubExecutor) ExecuteWorkflow(workflowFile string, params map[string]interface{}) (*ExecuteResult, error) {
	startTime := time.Now()

	if _, err := os.Stat(workflowFile); os.IsNotExist(err) {
		return NewExecuteResult("error").ErrorResult("Workflow file does not exist: "+workflowFile), nil
	}

	workflowJSON, err := loadWorkflowJSON(workflowFile)
	if err != nil {
		return NewExecuteResult("error").ErrorResult("Failed to load workflow: "+err.Error()), nil
	}

	workflowJSON, seedChanges := randomizeSeedInWorkflow(workflowJSON)

	parser := NewWorkflowParser()
	metadata := parser.ParseWorkflow(workflowJSON, "")
	if metadata == nil {
		return NewExecuteResult("error").ErrorResult("Cannot parse workflow metadata"), nil
	}

	workflowID := metadata.WorkflowID
	if workflowID == "" {
		return NewExecuteResult("error").ErrorResult("RunningHub workflow_id not found in metadata"), nil
	}

	nodeInfoList := e.ConvertParamsToNodeInfoList(metadata, params, seedChanges)

	taskData, err := e.client.CreateTask(workflowID, nodeInfoList)
	if err != nil {
		return NewExecuteResult("error").ErrorResult("Failed to create task: "+err.Error()), nil
	}

	taskID, ok := taskData["taskId"].(string)
	if !ok || taskID == "" {
		return NewExecuteResult("error").ErrorResult("Failed to get task ID"), nil
	}

	outputID2Var := extractOutputNodes(metadata)
	result, err := e.WaitForTaskCompletion(taskID, outputID2Var)
	if err != nil {
		return result, err
	}

	result.Duration = time.Since(startTime).Seconds()
	return result, nil
}

func (e *RunningHubExecutor) ConvertParamsToNodeInfoList(metadata *WorkflowMetadata, params map[string]interface{}, seedChanges map[string]int64) []map[string]interface{} {
	nodeInfoList := []map[string]interface{}{}

	for _, mapping := range metadata.MappingInfo.ParamMappings {
		if val, ok := params[mapping.ParamName]; ok {
			nodeClassType := mapping.NodeClassType
			needUpload := mapping.NeedUpload

			if needUpload || IsMediaUploadNode(nodeClassType) {
				val = e.HandleRunningHubMediaUpload(val)
			}

			nodeInfo := map[string]interface{}{
				"nodeId":    mapping.NodeID,
				"fieldName": mapping.InputField,
				"fieldValue": val,
			}
			nodeInfoList = append(nodeInfoList, nodeInfo)
		}
	}

	for nodeID, seedValue := range seedChanges {
		alreadySet := false
		for _, ni := range nodeInfoList {
			if ni["nodeId"] == nodeID && ni["fieldName"] == "seed" {
				alreadySet = true
				break
			}
		}
		if !alreadySet {
			nodeInfoList = append(nodeInfoList, map[string]interface{}{
				"nodeId":    nodeID,
				"fieldName": "seed",
				"fieldValue": seedValue,
			})
		}
	}

	return nodeInfoList
}

func (e *RunningHubExecutor) HandleRunningHubMediaUpload(paramValue interface{}) interface{} {
	if strVal, ok := paramValue.(string); ok {
		if strings.HasPrefix(strVal, "http://") || strings.HasPrefix(strVal, "https://") {
			mediaValue, err := e.UploadMediaFromURL(strVal)
			if err == nil {
				return mediaValue
			}
		} else if _, err := os.Stat(strVal); err == nil {
			uploadedName, err := e.client.UploadFile(strVal)
			if err == nil {
				return uploadedName
			}
		}
	}
	return paramValue
}

func (e *RunningHubExecutor) UploadMediaFromURL(mediaURL string) (string, error) {
	tempPath, err := downloadFile(mediaURL)
	if err != nil {
		return "", err
	}
	defer os.Remove(tempPath)
	return e.client.UploadFile(tempPath)
}

func downloadFile(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	tempFile, err := os.CreateTemp("", "workflow_*.json")
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	buf := make([]byte, 4096)
	for {
		n, err := resp.Body.Read(buf)
		if n == 0 || err != nil {
			break
		}
		tempFile.Write(buf[:n])
	}

	return tempFile.Name(), nil
}

func (e *RunningHubExecutor) WaitForTaskCompletion(taskID string, outputID2Var map[string]string) (*ExecuteResult, error) {
	maxWaitTime := e.timeout
	checkInterval := 2
	startTime := time.Now()

	for {
		elapsedTime := time.Since(startTime).Seconds()
		if maxWaitTime > 0 && elapsedTime >= float64(maxWaitTime) {
			return NewExecuteResult("error").ErrorResult("RunningHub task timeout after "+fmt.Sprintf("%d", maxWaitTime)+" seconds"), nil
		}

		statusInfo, err := e.client.QueryTaskStatus(taskID)
		if err != nil {
			time.Sleep(time.Duration(checkInterval) * time.Second)
			continue
		}

		taskStatus, _ := statusInfo["status"].(string)
		statusMsg, _ := statusInfo["msg"].(string)

		switch taskStatus {
		case "SUCCESS":
			resultData, err := e.client.QueryTaskResult(taskID)
			if err != nil {
				return NewExecuteResult("error").ErrorResult("Failed to get task result: "+err.Error()), nil
			}
			return e.ProcessTaskResult(taskID, resultData, outputID2Var), nil

		case "FAILED":
			errorMsg := "RunningHub task " + taskID + " failed"
			if statusMsg != "" {
				errorMsg += ": " + statusMsg
			}
			return NewExecuteResult("error").ErrorResult(errorMsg), nil

		case "QUEUED", "RUNNING":
			time.Sleep(time.Duration(checkInterval) * time.Second)
			continue
		}
	}
}

func (e *RunningHubExecutor) ProcessTaskResult(taskID string, resultData interface{}, outputID2Var map[string]string) *ExecuteResult {
	result := NewExecuteResult("completed")
	result.PromptID = taskID

	outputID2Images := map[string][]string{}
	outputID2Videos := map[string][]string{}
	outputID2Audios := map[string][]string{}
	outputID2Texts := map[string][]string{}

	if list, ok := resultData.([]interface{}); ok {
		for idx, item := range list {
			if itemMap, ok := item.(map[string]interface{}); ok {
				fileURL, _ := itemMap["fileUrl"].(string)
				fileType, _ := itemMap["fileType"].(string)
				nodeID, _ := itemMap["nodeId"].(string)
				if nodeID == "" {
					nodeID = fmt.Sprintf("%d", idx)
				}

				if fileURL != "" {
					fileType = strings.ToLower(fileType)
					switch {
					case strings.Contains(fileType, "image") ||
						strings.Contains(fileType, "png") ||
						strings.Contains(fileType, "jpg") ||
						strings.Contains(fileType, "jpeg") ||
						strings.Contains(fileType, "gif") ||
						strings.Contains(fileType, "webp"):
						outputID2Images[nodeID] = append(outputID2Images[nodeID], fileURL)

					case strings.Contains(fileType, "video") ||
						strings.Contains(fileType, "mp4") ||
						strings.Contains(fileType, "avi") ||
						strings.Contains(fileType, "mov") ||
						strings.Contains(fileType, "mkv"):
						outputID2Videos[nodeID] = append(outputID2Videos[nodeID], fileURL)

					case strings.Contains(fileType, "audio") ||
						strings.Contains(fileType, "mp3") ||
						strings.Contains(fileType, "wav") ||
						strings.Contains(fileType, "flac"):
						outputID2Audios[nodeID] = append(outputID2Audios[nodeID], fileURL)

					case strings.Contains(fileType, "text") ||
						strings.Contains(fileType, "txt") ||
						strings.Contains(fileType, "json") ||
						strings.Contains(fileType, "xml"):
						textContent := e.DownloadTextFromURL(fileURL)
						if textContent != "" {
							outputID2Texts[nodeID] = append(outputID2Texts[nodeID], textContent)
						}
					}
				}
			}
		}
	}

	if len(outputID2Images) > 0 {
		result.ImagesByVar = mapOutputsByVar(outputID2Var, outputID2Images)
		result.Images = extendFlatListFromDict(result.ImagesByVar)
	}
	if len(outputID2Videos) > 0 {
		result.VideosByVar = mapOutputsByVar(outputID2Var, outputID2Videos)
		result.Videos = extendFlatListFromDict(result.VideosByVar)
	}
	if len(outputID2Audios) > 0 {
		result.AudiosByVar = mapOutputsByVar(outputID2Var, outputID2Audios)
		result.Audios = extendFlatListFromDict(result.AudiosByVar)
	}
	if len(outputID2Texts) > 0 {
		result.TextsByVar = mapOutputsByVar(outputID2Var, outputID2Texts)
		result.Texts = extendFlatListFromDict(result.TextsByVar)
	}

	result.Outputs = map[string]interface{}{"raw_data": resultData}
	return result
}

func (e *RunningHubExecutor) DownloadTextFromURL(textURL string) string {
	resp, err := http.Get(textURL)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ""
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return string(body)
}

package comfyui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type RunningHubClient struct {
	apiKey      string
	baseURL     string
	timeout     int
	retryCount  int
	instanceType string
	client      *http.Client
}

func NewRunningHubClient(apiKey, baseURL string, timeout, retryCount int, instanceType string) *RunningHubClient {
	if baseURL == "" {
		baseURL = "https://www.runninghub.ai"
	}
	if timeout <= 0 {
		timeout = 300
	}
	if retryCount <= 0 {
		retryCount = 3
	}

	return &RunningHubClient{
		apiKey:       apiKey,
		baseURL:      trimTrailingSlash(baseURL),
		timeout:      timeout,
		retryCount:   retryCount,
		instanceType: instanceType,
		client: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
	}
}

func (c *RunningHubClient) Close() error {
	return nil
}

func (c *RunningHubClient) makeRequest(method, endpoint string, data map[string]interface{}, files map[string]fileInfo) (map[string]interface{}, error) {
	url := c.baseURL + endpoint

	var req *http.Request
	var err error

	if files != nil {
		var buf bytes.Buffer
		writer := multipart.NewWriter(&buf)

		if data != nil {
			for key, value := range data {
				writer.WriteField(key, fmt.Sprintf("%v", value))
			}
		}

		for key, fi := range files {
			part, err := writer.CreateFormFile(key, fi.filename)
			if err != nil {
				return nil, err
			}
			part.Write(fi.content)
		}
		writer.Close()

		req, err = http.NewRequest(method, url, &buf)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())
	} else {
		var jsonData []byte
		if data != nil {
			jsonData, err = json.Marshal(data)
			if err != nil {
				return nil, err
			}
		}

		req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonData))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
	}

	var lastErr error
	for attempt := 0; attempt <= c.retryCount; attempt++ {
		resp, err := c.client.Do(req)
		if err != nil {
			lastErr = err
			if attempt < c.retryCount {
				time.Sleep(time.Duration(1<<attempt) * time.Second)
				continue
			}
			return nil, lastErr
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			var result map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				return nil, err
			}

			if code, ok := result["code"].(float64); ok && code == 0 {
				return result, nil
			} else {
				msg := result["msg"].(string)
				return nil, fmt.Errorf("RunningHub API error: %s", msg)
			}
		}

		body, _ := io.ReadAll(resp.Body)
		lastErr = fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
		if attempt < c.retryCount {
			time.Sleep(time.Duration(1<<attempt) * time.Second)
		}
	}

	return nil, lastErr
}

type fileInfo struct {
	content  []byte
	filename string
}

func (c *RunningHubClient) GetWorkflowJSON(workflowID string) (map[string]interface{}, error) {
	data := map[string]interface{}{
		"apiKey":      c.apiKey,
		"workflowId": workflowID,
	}

	result, err := c.makeRequest("POST", "/api/openapi/getJsonApiFormat", data, nil)
	if err != nil {
		return nil, err
	}

	dataMap, ok := result["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("No workflow JSON found in response")
	}

	promptStr, ok := dataMap["prompt"].(string)
	if !ok || promptStr == "" {
		return nil, fmt.Errorf("No workflow JSON found in response")
	}

	var workflowJSON map[string]interface{}
	if err := json.Unmarshal([]byte(promptStr), &workflowJSON); err != nil {
		return nil, err
	}

	return workflowJSON, nil
}

func (c *RunningHubClient) UploadFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	data := map[string]interface{}{
		"apiKey": c.apiKey,
	}

	files := map[string]fileInfo{
		"file": {
			content:  content,
			filename: filepath.Base(filePath),
		},
	}

	result, err := c.makeRequest("POST", "/task/openapi/upload", data, files)
	if err != nil {
		return "", err
	}

	dataMap, ok := result["data"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("No data in response")
	}

	if fileName, ok := dataMap["fileName"].(string); ok && fileName != "" {
		return fileName, nil
	}

	if fileURL, ok := dataMap["url"].(string); ok && fileURL != "" {
		return fileURL, nil
	}

	return "", fmt.Errorf("Neither fileName nor URL found in upload response")
}

func (c *RunningHubClient) CreateTask(workflowID string, nodeInfoList []map[string]interface{}) (map[string]interface{}, error) {
	data := map[string]interface{}{
		"apiKey":      c.apiKey,
		"workflowId": workflowID,
	}

	if len(nodeInfoList) > 0 {
		data["nodeInfoList"] = nodeInfoList
	}

	if c.instanceType != "" {
		data["instanceType"] = c.instanceType
	}

	result, err := c.makeRequest("POST", "/task/openapi/create", data, nil)
	if err != nil {
		return nil, err
	}

	return result["data"].(map[string]interface{}), nil
}

func (c *RunningHubClient) QueryTaskStatus(taskID string) (map[string]interface{}, error) {
	data := map[string]interface{}{
		"apiKey": c.apiKey,
		"taskId": taskID,
	}

	result, err := c.makeRequest("POST", "/task/openapi/status", data, nil)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"status": result["data"],
		"msg":    result["msg"],
		"code":   result["code"],
	}, nil
}

func (c *RunningHubClient) QueryTaskResult(taskID string) (interface{}, error) {
	data := map[string]interface{}{
		"apiKey": c.apiKey,
		"taskId": taskID,
	}

	result, err := c.makeRequest("POST", "/task/openapi/outputs", data, nil)
	if err != nil {
		return nil, err
	}

	return result["data"], nil
}

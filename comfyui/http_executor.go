package comfyui

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type HTTPExecutor struct {
	*ComfyUIExecutor
	client *http.Client
}

func NewHTTPExecutor(baseURL, apiKey, cookies string) *HTTPExecutor {
	return &HTTPExecutor{
		ComfyUIExecutor: NewComfyUIExecutor(baseURL, apiKey, cookies),
		client:          &http.Client{Timeout: 300 * time.Second},
	}
}

func (e *HTTPExecutor) Close() {
}

func (e *HTTPExecutor) ExecuteWorkflow(workflowFile string, params map[string]interface{}) (*ExecuteResult, error) {
	if _, err := os.Stat(workflowFile); os.IsNotExist(err) {
		return NewExecuteResult("error").ErrorResult("Workflow file does not exist: " + workflowFile), nil
	}

	metadata := e.getWorkflowMetadata(workflowFile)
	if metadata == nil {
		return NewExecuteResult("error").ErrorResult("Cannot parse workflow metadata"), nil
	}

	workflowData, err := loadWorkflowJSON(workflowFile)
	if err != nil {
		return NewExecuteResult("error").ErrorResult("Failed to load workflow: " + err.Error()), nil
	}

	if params != nil {
		workflowData = e.applyParamsToWorkflow(workflowData, metadata, params)
	} else {
		workflowData = e.applyParamsToWorkflow(workflowData, metadata, map[string]interface{}{})
	}

	workflowData, _ = randomizeSeedInWorkflow(workflowData)
	outputID2Var := extractOutputNodes(metadata)

	clientID := generateUUID()
	promptExtParams := map[string]interface{}{}
	if e.APIKey != "" {
		promptExtParams["extra_data"] = map[string]interface{}{
			"api_key_comfy_org": e.APIKey,
		}
	}

	promptID, err := e.QueuePrompt(workflowData, clientID, promptExtParams)
	if err != nil {
		return NewExecuteResult("error").ErrorResult("Submit workflow failed: " + err.Error()), nil
	}

	return e.WaitForResults(promptID, clientID, nil, outputID2Var)
}

func (e *HTTPExecutor) QueuePrompt(workflow map[string]interface{}, clientID string, extParams map[string]interface{}) (string, error) {
	promptData := map[string]interface{}{
		"prompt":    workflow,
		"client_id": clientID,
	}

	for k, v := range extParams {
		promptData[k] = v
	}

	data, err := json.Marshal(promptData)
	if err != nil {
		return "", err
	}

	url := e.BaseURL + "/prompt"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	e.setAuthHeaders(req)

	resp, err := e.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if promptID, ok := result["prompt_id"].(string); ok && promptID != "" {
		return promptID, nil
	}

	return "", fmt.Errorf("Get prompt_id failed: %v", result)
}

func (e *HTTPExecutor) WaitForResults(promptID, clientID string, timeout *int, outputID2Var map[string]string) (*ExecuteResult, error) {
	startTime := time.Now()
	result := NewExecuteResult("processing")
	result.PromptID = promptID

	for {
		if timeout != nil && *timeout > 0 {
			duration := time.Since(startTime).Seconds()
			if duration > float64(*timeout) {
				return result.TimeoutResult(duration), nil
			}
		}

		url := e.BaseURL + "/history/" + promptID
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		e.setAuthHeaders(req)

		resp, err := e.client.Do(req)
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			time.Sleep(1 * time.Second)
			continue
		}

		var historyData map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&historyData); err != nil {
			resp.Body.Close()
			time.Sleep(1 * time.Second)
			continue
		}
		resp.Body.Close()

		promptHistory, ok := historyData[promptID]
		if !ok {
			time.Sleep(1 * time.Second)
			continue
		}

		historyMap, ok := promptHistory.(map[string]interface{})
		if !ok {
			time.Sleep(1 * time.Second)
			continue
		}

		if status, ok := historyMap["status"].(map[string]interface{}); ok {
			if statusStr, ok := status["status_str"].(string); ok && statusStr == "error" {
				result.Status = "error"
				if messages, ok := status["messages"].([]interface{}); ok {
					var errors []string
					for _, msg := range messages {
						if msgMap, ok := msg.([]interface{}); ok && len(msgMap) >= 2 {
							if body, ok := msgMap[1].(map[string]interface{}); ok {
								if excMsg, ok := body["exception_message"].(string); ok {
									errors = append(errors, excMsg)
								}
							}
						}
					}
					result.Message = strings.Join(errors, "\n")
				} else {
					result.Message = "Unknown error"
				}
				result.Duration = time.Since(startTime).Seconds()
				return result, nil
			}
		}

		if outputs, ok := historyMap["outputs"]; ok {
			result.Outputs = outputs.(map[string]interface{})
			result.Status = "completed"

			e.processOutputs(result, outputs.(map[string]interface{}), outputID2Var, e.BaseURL)
			result.Duration = time.Since(startTime).Seconds()
			return result, nil
		}

		time.Sleep(1 * time.Second)
	}
}

func (e *HTTPExecutor) processOutputs(result *ExecuteResult, outputs map[string]interface{}, outputID2Var map[string]string, baseURL string) {
	outputID2Images := map[string][]string{}
	outputID2Videos := map[string][]string{}
	outputID2Audios := map[string][]string{}
	outputID2Texts := map[string][]string{}

	for nodeID, nodeOutput := range outputs {
		nodeOutMap, ok := nodeOutput.(map[string]interface{})
		if !ok {
			continue
		}

		images, videos, audios := splitMediaBySuffix(nodeOutMap, baseURL)
		if len(images) > 0 {
			outputID2Images[nodeID] = images
		}
		if len(videos) > 0 {
			outputID2Videos[nodeID] = videos
		}
		if len(audios) > 0 {
			outputID2Audios[nodeID] = audios
		}

		if texts, ok := nodeOutMap["text"]; ok {
			var textList []string
			switch v := texts.(type) {
			case string:
				textList = []string{v}
			case []interface{}:
				for _, t := range v {
					textList = append(textList, fmt.Sprintf("%v", t))
				}
			default:
				textList = []string{fmt.Sprintf("%v", texts)}
			}
			outputID2Texts[nodeID] = textList
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
}

func (e *HTTPExecutor) getWorkflowMetadata(workflowFile string) *WorkflowMetadata {
	parser := NewWorkflowParser()
	return parser.ParseWorkflowFile(workflowFile, "")
}

func (e *HTTPExecutor) applyParamsToWorkflow(workflowData map[string]interface{}, metadata *WorkflowMetadata, params map[string]interface{}) map[string]interface{} {
	for _, mapping := range metadata.MappingInfo.ParamMappings {
		if val, ok := params[mapping.ParamName]; ok {
			e.applyParamMapping(workflowData, mapping, val)
		} else {
			if paramInfo, ok := metadata.Params[mapping.ParamName]; ok {
				if paramInfo.Default != nil {
					e.applyParamMapping(workflowData, mapping, paramInfo.Default)
				} else if paramInfo.Required {
					panic(fmt.Sprintf("Required parameter '%s' is missing", mapping.ParamName))
				}
			}
		}
	}
	return workflowData
}

func (e *HTTPExecutor) applyParamMapping(workflowData map[string]interface{}, mapping WorkflowParamMapping, paramValue interface{}) {
	nodeData, ok := workflowData[mapping.NodeID].(map[string]interface{})
	if !ok {
		return
	}

	if _, ok := nodeData["inputs"]; !ok {
		nodeData["inputs"] = map[string]interface{}{}
	}

	inputs := nodeData["inputs"].(map[string]interface{})

	if mapping.NeedUpload || IsMediaUploadNode(mapping.NodeClassType) {
		e.handleMediaUpload(nodeData, mapping.InputField, paramValue)
	} else {
		inputs[mapping.InputField] = paramValue
	}
}

func (e *HTTPExecutor) handleMediaUpload(nodeData map[string]interface{}, inputField string, paramValue interface{}) {
	if _, ok := nodeData["inputs"]; !ok {
		nodeData["inputs"] = map[string]interface{}{}
	}
	inputs := nodeData["inputs"].(map[string]interface{})

	if strVal, ok := paramValue.(string); ok {
		if strings.HasPrefix(strVal, "http://") || strings.HasPrefix(strVal, "https://") {
			mediaValue, err := e.uploadMediaFromSource(strVal)
			if err == nil {
				inputs[inputField] = mediaValue
			}
		} else if _, err := os.Stat(strVal); err == nil {
			uploadedName, err := e.uploadMedia(strVal)
			if err == nil {
				inputs[inputField] = uploadedName
			}
		}
	} else {
		inputs[inputField] = paramValue
	}
}

func (e *HTTPExecutor) uploadMediaFromSource(mediaURL string) (string, error) {
	resp, err := http.Get(mediaURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	return e.uploadMediaFromReader(resp.Body, mediaURL)
}

func (e *HTTPExecutor) uploadMediaFromReader(reader io.Reader, fileName string) (string, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("image", filepath.Base(fileName))
	if err != nil {
		return "", err
	}
	io.Copy(part, reader)
	writer.Close()

	url := e.BaseURL + "/upload/image"
	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	e.setAuthHeaders(req)

	resp, err := e.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if name, ok := result["name"].(string); ok {
		return name, nil
	}
	return "", fmt.Errorf("No name in response")
}

func (e *HTTPExecutor) uploadMedia(mediaPath string) (string, error) {
	file, err := os.Open(mediaPath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	return e.uploadMediaFromReader(file, mediaPath)
}

func (e *HTTPExecutor) setAuthHeaders(req *http.Request) {
	if e.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+e.APIKey)
	}
	if e.Cookies != "" {
		req.Header.Set("Cookie", e.Cookies)
	}
}

func generateUUID() string {
	b := make([]byte, 16)
	rand.Read(b)
	b[6] = (b[6] & 0x0F) | 0x40
	b[8] = (b[8] & 0x3F) | 0x80
	return hex.EncodeToString(b)
}

func loadWorkflowJSON(filePath string) (map[string]interface{}, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

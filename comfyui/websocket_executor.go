package comfyui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

type WebSocketExecutor struct {
	*ComfyUIExecutor
	httpBaseURL   string
	wsBaseURL     string
	conn          *websocket.Conn
	httpExecutor  *HTTPExecutor
}

func NewWebSocketExecutor(baseURL, apiKey, cookies string) *WebSocketExecutor {
	return &WebSocketExecutor{
		ComfyUIExecutor: NewComfyUIExecutor(baseURL, apiKey, cookies),
		httpExecutor:    NewHTTPExecutor(baseURL, apiKey, cookies),
	}
}

func (e *WebSocketExecutor) ParseWsUrl() {
	if e.httpBaseURL != "" && e.wsBaseURL != "" {
		return
	}
	parsed, _ := url.Parse(e.BaseURL)
	wsScheme := "ws"
	httpScheme := "http"
	if parsed.Scheme == "https" {
		wsScheme = "wss"
		httpScheme = "https"
	}

	e.httpBaseURL = httpScheme + "://" + parsed.Host + parsed.Path
	e.wsBaseURL = wsScheme + "://" + parsed.Host + parsed.Path + "/ws"
}

func (e *WebSocketExecutor) Close() {
	if e.conn != nil {
		e.conn.Close()
	}
}

func (e *WebSocketExecutor) ExecuteWorkflow(workflowFile string, params map[string]interface{}) (*ExecuteResult, error) {
	startTime := time.Now()

	if _, err := os.Stat(workflowFile); os.IsNotExist(err) {
		return NewExecuteResult("error").ErrorResult("Workflow file does not exist: "+workflowFile), nil
	}

	metadata := e.getWorkflowMetadata(workflowFile)
	if metadata == nil {
		return NewExecuteResult("error").ErrorResult("Cannot parse workflow metadata"), nil
	}

	workflowData, err := loadWorkflowJSON(workflowFile)
	if err != nil {
		return NewExecuteResult("error").ErrorResult("Failed to load workflow: "+err.Error()), nil
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

	timeout := 30 * 60
	e.ParseWsUrl()
	wsURL := e.wsBaseURL + "?clientId=" + clientID

	collectedOutputs := map[string]interface{}{}
	var promptID string

	headers := map[string][]string{}
	if e.Cookies != "" {
		headers["Cookie"] = []string{e.Cookies}
	}

	c, _, err := websocket.DefaultDialer.Dial(wsURL, headers)
	if err != nil {
		return NewExecuteResult("error").ErrorResult("WebSocket connection failed: "+err.Error()), nil
	}
	defer c.Close()

	promptID, err = e.QueuePrompt(workflowData, clientID, promptExtParams)
	if err != nil {
		return NewExecuteResult("error").ErrorResult("Submit workflow failed: "+err.Error()), nil
	}

	c.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Second))

	for {
		elapsed := time.Since(startTime).Seconds()
		if elapsed > float64(timeout) {
			return NewExecuteResult("timeout").TimeoutResult(elapsed), nil
		}

		_, message, err := c.ReadMessage()
		if err != nil {
			return NewExecuteResult("error").ErrorResult("WebSocket read error: "+err.Error()), nil
		}

		var msg map[string]interface{}
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

		msgType, _ := msg["type"].(string)
		data, _ := msg["data"].(map[string]interface{})

		if msgType == "executed" {
			nodeID, _ := data["node"].(string)
			output, _ := data["output"].(map[string]interface{})

			if output != nil && nodeID != "" {
				hasMedia := false
				if _, ok := output["images"]; ok {
					hasMedia = true
				}
				if _, ok := output["gifs"]; ok {
					hasMedia = true
				}
				if _, ok := output["audio"]; ok {
					hasMedia = true
				}
				if _, ok := output["text"]; ok {
					hasMedia = true
				}

				if hasMedia {
					collectedOutputs[nodeID] = output
				}
			}
		} else if msgType == "execution_error" {
			errorMsg, _ := data["exception_message"].(string)
			if errorMsg == "" {
				errorMsg = "Unknown error"
			}
			return NewExecuteResult("error").ErrorResult(errorMsg), nil
		}

		if msgType == "executing" {
			if data["node"] == nil {
				if promptIDVal, _ := data["prompt_id"].(string); promptIDVal == promptID {
					duration := time.Since(startTime).Seconds()

					if len(collectedOutputs) > 0 {
						result := e.BuildResultFromCollectedOutputs(collectedOutputs, promptID, outputID2Var)
						result.Duration = duration
						return result, nil
					} else {
						return NewExecuteResult("error").ErrorResult("WebSocket did not collect any outputs"), nil
					}
				}
			}
		}
	}
}

func (e *WebSocketExecutor) QueuePrompt(workflow map[string]interface{}, clientID string, extParams map[string]interface{}) (string, error) {
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
	if e.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+e.APIKey)
	}
	if e.Cookies != "" {
		req.Header.Set("Cookie", e.Cookies)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioReadAll(resp.Body)
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

func ioReadAll(r io.Reader) ([]byte, error) {
	var buf []byte
	buf = make([]byte, 0, 1024)
	for {
		tmp := make([]byte, 1024)
		n, err := r.Read(tmp)
		if n > 0 {
			buf = append(buf, tmp[:n]...)
		}
		if err != nil {
			if err == io.EOF {
				return buf, nil
			}
			return buf, err
		}
	}
}

func (e *WebSocketExecutor) BuildResultFromCollectedOutputs(collectedOutputs map[string]interface{}, promptID string, outputID2Var map[string]string) *ExecuteResult {
	result := NewExecuteResult("completed")
	result.PromptID = promptID
	result.Outputs = collectedOutputs

	outputID2Images := map[string][]string{}
	outputID2Videos := map[string][]string{}
	outputID2Audios := map[string][]string{}
	outputID2Texts := map[string][]string{}

	for nodeID, output := range collectedOutputs {
		nodeOutMap, ok := output.(map[string]interface{})
		if !ok {
			continue
		}

		e.ParseWsUrl()
		images, videos, audios := splitMediaBySuffix(nodeOutMap, e.httpBaseURL)
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

	if len(result.Images) == 0 && len(result.Videos) == 0 && len(result.Audios) == 0 && len(result.Texts) == 0 {
		result.Status = "error"
		result.Message = "No outputs found"
	}

	return result
}

func (e *WebSocketExecutor) getWorkflowMetadata(workflowFile string) *WorkflowMetadata {
	parser := NewWorkflowParser()
	return parser.ParseWorkflowFile(workflowFile, "")
}

func (e *WebSocketExecutor) applyParamsToWorkflow(workflowData map[string]interface{}, metadata *WorkflowMetadata, params map[string]interface{}) map[string]interface{} {
	return e.httpExecutor.applyParamsToWorkflow(workflowData, metadata, params)
}

func (e *WebSocketExecutor) applyParamMapping(workflowData map[string]interface{}, mapping WorkflowParamMapping, paramValue interface{}) {
	e.httpExecutor.applyParamMapping(workflowData, mapping, paramValue)
}

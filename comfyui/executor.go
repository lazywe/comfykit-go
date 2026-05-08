package comfyui

import (
	"math/rand"
	"time"
)

type Executor interface {
	ExecuteWorkflow(workflowFile string, params map[string]interface{}) (*ExecuteResult, error)
}

type ComfyUIExecutor struct {
	BaseURL string
	APIKey  string
	Cookies string
}

func NewComfyUIExecutor(baseURL, apiKey, cookies string) *ComfyUIExecutor {
	if baseURL == "" {
		baseURL = "http://127.0.0.1:8188"
	}
	return &ComfyUIExecutor{
		BaseURL: trimTrailingSlash(baseURL),
		APIKey:  apiKey,
		Cookies: cookies,
	}
}

func trimTrailingSlash(s string) string {
	for len(s) > 0 && s[len(s)-1] == '/' {
		s = s[:len(s)-1]
	}
	return s
}

var mediaUploadNodeTypes = map[string]bool{
	"LoadImage":           true,
	"LoadAudio":           true,
	"LoadVideo":           true,
	"VHS_LoadAudioUpload": true,
	"VHS_LoadVideo":       true,
}

func IsMediaUploadNode(nodeType string) bool {
	return mediaUploadNodeTypes[nodeType]
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func generate63BitSeed() int64 {
	return rand.Int63()
}

func randomizeSeedInWorkflow(workflowData map[string]interface{}) (map[string]interface{}, map[string]int64) {
	changed := map[string]int64{}

	for nodeID, nodeVal := range workflowData {
		node, ok := nodeVal.(map[string]interface{})
		if !ok {
			continue
		}

		inputs, ok := node["inputs"].(map[string]interface{})
		if !ok {
			continue
		}

		if seedVal, ok := inputs["seed"]; ok {
			isZero := false
			switch v := seedVal.(type) {
			case int:
				isZero = v == 0
			case int64:
				isZero = v == 0
			case string:
				isZero = v == "0"
			}

			if isZero {
				newSeed := generate63BitSeed()
				inputs["seed"] = newSeed
				changed[nodeID] = newSeed
			}
		}
	}

	return workflowData, changed
}

func extractOutputNodes(metadata *WorkflowMetadata) map[string]string {
	outputID2Var := map[string]string{}
	for _, mapping := range metadata.MappingInfo.OutputMappings {
		outputID2Var[mapping.NodeID] = mapping.OutputVar
	}
	return outputID2Var
}

func mapOutputsByVar(outputID2Var map[string]string, outputID2Media map[string][]string) map[string][]string {
	result := map[string][]string{}
	for nodeID, media := range outputID2Media {
		varName := outputID2Var[nodeID]
		if varName == "" {
			varName = nodeID
		}
		result[varName] = media
	}
	return result
}

func extendFlatListFromDict(mediaDict map[string][]string) []string {
	var flat []string
	for _, items := range mediaDict {
		flat = append(flat, items...)
	}
	return flat
}

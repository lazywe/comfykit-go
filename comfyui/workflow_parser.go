package comfyui

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type WorkflowParser struct{}

func NewWorkflowParser() *WorkflowParser {
	return &WorkflowParser{}
}

func (p *WorkflowParser) ParseTitle(title string) []string {
	if title == "" {
		return []string{}
	}

	parts := strings.Split(title, ",")
	var markers []string
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "$") {
			markers = append(markers, part)
		}
	}
	return markers
}

func (p *WorkflowParser) ParseParamMarker(marker string) map[string]interface{} {
	if marker == "" || !strings.HasPrefix(marker, "$") {
		return nil
	}

	marker = marker[1:]

	if idx := strings.Index(marker, ":"); idx != -1 {
		marker = marker[:idx]
	}

	required := strings.HasSuffix(marker, "!")
	if required {
		marker = marker[:len(marker)-1]
	}

	upload := strings.HasPrefix(marker, "~")
	if upload {
		marker = marker[1:]
	}

	var paramName, fieldName string
	if idx := strings.Index(marker, "."); idx != -1 {
		paramName = marker[:idx]
		fieldPart := marker[idx+1:]

		if strings.HasPrefix(fieldPart, "~") {
			upload = true
			fieldName = fieldPart[1:]
		} else {
			fieldName = fieldPart
		}
	} else {
		paramName = marker
		fieldName = marker
	}

	return map[string]interface{}{
		"name":     paramName,
		"field":    fieldName,
		"upload":   upload,
		"required": required,
	}
}

func (p *WorkflowParser) InferTypeFromValue(value interface{}) string {
	switch value.(type) {
	case bool:
		return "bool"
	case int, int64, float64:
		if _, ok := value.(float64); ok {
			return "float"
		}
		return "int"
	default:
		return "str"
	}
}

func (p *WorkflowParser) ExtractFieldValue(nodeData map[string]interface{}, fieldName string) interface{} {
	inputs, ok := nodeData["inputs"].(map[string]interface{})
	if !ok {
		return nil
	}

	val, ok := inputs[fieldName]
	if !ok {
		return nil
	}

	if _, ok := val.([]interface{}); ok {
		return nil
	}

	return val
}

func (p *WorkflowParser) ParseOutputMarker(title string) string {
	if !strings.HasPrefix(title, "$output.") {
		return ""
	}
	return title[8:]
}

func (p *WorkflowParser) IsKnownOutputNode(classType string) bool {
	knownTypes := map[string]bool{
		"SaveImage":     true,
		"SaveVideo":     true,
		"SaveAudio":     true,
		"VHS_SaveVideo": true,
		"VHS_SaveAudio": true,
	}
	return knownTypes[classType]
}

func (p *WorkflowParser) ParseNode(nodeID string, nodeData map[string]interface{}) ([]WorkflowParam, []WorkflowParamMapping, *WorkflowOutputMapping) {
	if _, ok := nodeData["_meta"]; !ok {
		return nil, nil, nil
	}

	meta, _ := nodeData["_meta"].(map[string]interface{})
	title, _ := meta["title"].(string)
	classType, _ := nodeData["class_type"].(string)

	if outputVar := p.ParseOutputMarker(title); outputVar != "" {
		outputMapping := &WorkflowOutputMapping{
			NodeID:    nodeID,
			OutputVar: outputVar,
		}
		return nil, nil, outputMapping
	}

	if p.IsKnownOutputNode(classType) {
		outputMapping := &WorkflowOutputMapping{
			NodeID:    nodeID,
			OutputVar: nodeID,
		}
		return nil, nil, outputMapping
	}

	paramMarkers := p.ParseTitle(title)
	if len(paramMarkers) == 0 {
		return nil, nil, nil
	}

	var params []WorkflowParam
	var mappings []WorkflowParamMapping

	for _, marker := range paramMarkers {
		paramInfo := p.ParseParamMarker(marker)
		if paramInfo == nil {
			continue
		}

		paramName := paramInfo["name"].(string)
		fieldName := paramInfo["field"].(string)
		isRequired := paramInfo["required"].(bool)
		needUpload := paramInfo["upload"].(bool)

		defaultValue := p.ExtractFieldValue(nodeData, fieldName)
		paramType := "str"
		if defaultValue != nil {
			paramType = p.InferTypeFromValue(defaultValue)
		}

		if !isRequired && defaultValue == nil {
			panic("Parameter `" + paramName + "` has no default value but not marked as required (node " + nodeID + ")")
		}

		param := WorkflowParam{
			Name:       paramName,
			Type:       paramType,
			Required:   isRequired,
			NeedUpload: needUpload,
		}
		if !isRequired {
			param.Default = defaultValue
		}

		paramMapping := WorkflowParamMapping{
			ParamName:     paramName,
			NodeID:        nodeID,
			InputField:    fieldName,
			NodeClassType: classType,
			NeedUpload:    needUpload,
		}

		params = append(params, param)
		mappings = append(mappings, paramMapping)
	}

	return params, mappings, nil
}

func (p *WorkflowParser) ParseWorkflow(workflowData map[string]interface{}, title string) *WorkflowMetadata {
	params := map[string]WorkflowParam{}
	paramMappings := []WorkflowParamMapping{}
	outputMappings := []WorkflowOutputMapping{}

	for nodeID, nodeVal := range workflowData {
		nodeData, ok := nodeVal.(map[string]interface{})
		if !ok {
			continue
		}

		paramsList, mappingsList, outputMapping := p.ParseNode(nodeID, nodeData)

		for _, param := range paramsList {
			if _, exists := params[param.Name]; exists {
				continue
			}
			params[param.Name] = param
		}

		paramMappings = append(paramMappings, mappingsList...)

		if outputMapping != nil {
			outputMappings = append(outputMappings, *outputMapping)
		}
	}

	mappingInfo := WorkflowMappingInfo{
		ParamMappings:  paramMappings,
		OutputMappings: outputMappings,
	}

	return &WorkflowMetadata{
		Title:       title,
		Params:      params,
		MappingInfo: mappingInfo,
	}
}

func (p *WorkflowParser) ParseWorkflowFile(filePath string, toolName string) *WorkflowMetadata {
	file, err := os.Open(filePath)
	if err != nil {
		return nil
	}
	defer file.Close()

	var workflowData map[string]interface{}
	if err := json.NewDecoder(file).Decode(&workflowData); err != nil {
		return nil
	}

	title := toolName
	if title == "" {
		title = filepath.Base(filePath)
		if idx := strings.LastIndex(title, "."); idx != -1 {
			title = title[:idx]
		}
	}

	return p.ParseWorkflow(workflowData, title)
}

func splitMediaBySuffix(nodeOutput map[string]interface{}, baseURL string) ([]string, []string, []string) {
	imageExts := map[string]bool{".png": true, ".jpg": true, ".jpeg": true, ".webp": true, ".bmp": true, ".tiff": true}
	videoExts := map[string]bool{".mp4": true, ".mov": true, ".avi": true, ".webm": true, ".gif": true}
	audioExts := map[string]bool{".mp3": true, ".wav": true, ".flac": true, ".ogg": true, ".aac": true, ".m4a": true, ".wma": true, ".opus": true}

	var images, videos, audios []string

	for _, mediaKey := range []string{"images", "gifs", "audio"} {
		if mediaList, ok := nodeOutput[mediaKey].([]interface{}); ok {
			for _, media := range mediaList {
				if mediaMap, ok := media.(map[string]interface{}); ok {
					filename, _ := mediaMap["filename"].(string)
					subfolder, _ := mediaMap["subfolder"].(string)
					mediaType, _ := mediaMap["type"].(string)

					url := baseURL + "/view?filename=" + filename
					if subfolder != "" {
						url += "&subfolder=" + subfolder
					}
					if mediaType != "" {
						url += "&type=" + mediaType
					}

					ext := strings.ToLower(filepath.Ext(filename))
					if imageExts[ext] {
						images = append(images, url)
					} else if videoExts[ext] {
						videos = append(videos, url)
					} else if audioExts[ext] {
						audios = append(audios, url)
					}
				}
			}
		}
	}

	return images, videos, audios
}

func IsRunningHubWorkflow(workflowFile string) bool {
	data := getWorkflowSourceData(workflowFile)
	if data != nil && data["_source"] == "runninghub" {
		return true
	}
	return false
}

func getWorkflowSourceData(workflowFile string) map[string]interface{} {
	file, err := os.Open(workflowFile)
	if err != nil {
		return nil
	}
	defer file.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return nil
	}
	return data
}

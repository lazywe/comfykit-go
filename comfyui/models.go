package comfyui

import "encoding/json"

type ExecuteResult struct {
	Status       string                 `json:"status"`
	PromptID     string                 `json:"prompt_id,omitempty"`
	Duration     float64                `json:"duration,omitempty"`
	Images       []string               `json:"images"`
	ImagesByVar  map[string][]string    `json:"images_by_var"`
	Videos       []string               `json:"videos"`
	VideosByVar  map[string][]string    `json:"videos_by_var"`
	Audios       []string               `json:"audios"`
	AudiosByVar  map[string][]string    `json:"audios_by_var"`
	Texts        []string               `json:"texts"`
	TextsByVar   map[string][]string    `json:"texts_by_var"`
	Outputs      map[string]interface{} `json:"outputs,omitempty"`
	Message      string                 `json:"msg,omitempty"`
}

func NewExecuteResult(status string) *ExecuteResult {
	return &ExecuteResult{
		Status:      status,
		Images:      []string{},
		ImagesByVar: map[string][]string{},
		Videos:      []string{},
		VideosByVar: map[string][]string{},
		Audios:      []string{},
		AudiosByVar: map[string][]string{},
		Texts:       []string{},
		TextsByVar:  map[string][]string{},
	}
}

func (r *ExecuteResult) ErrorResult(msg string) *ExecuteResult {
	r.Status = "error"
	r.Message = msg
	return r
}

func (r *ExecuteResult) TimeoutResult(duration float64) *ExecuteResult {
	r.Status = "timeout"
	r.Duration = duration
	return r
}

func (r *ExecuteResult) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}

type WorkflowParam struct {
	Name       string      `json:"name"`
	Type       string      `json:"type"`
	Required   bool        `json:"required"`
	Default    interface{} `json:"default,omitempty"`
	NeedUpload bool        `json:"need_upload"`
}

type WorkflowParamMapping struct {
	ParamName     string `json:"param_name"`
	NodeID        string `json:"node_id"`
	InputField    string `json:"input_field"`
	NodeClassType string `json:"node_class_type"`
	NeedUpload    bool   `json:"need_upload"`
}

type WorkflowOutputMapping struct {
	NodeID    string `json:"node_id"`
	OutputVar string `json:"output_var"`
}

type WorkflowMappingInfo struct {
	ParamMappings []WorkflowParamMapping `json:"param_mappings"`
	OutputMappings []WorkflowOutputMapping `json:"output_mappings"`
}

type WorkflowMetadata struct {
	Title        string              `json:"title"`
	Params       map[string]WorkflowParam `json:"params"`
	MappingInfo  WorkflowMappingInfo `json:"mapping_info"`
	WorkflowID   string              `json:"workflow_id,omitempty"`
	IsRunningHub bool                `json:"is_runninghub"`
}

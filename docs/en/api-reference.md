# API Reference

## ComfyKit

### Constructor

```go
func NewComfyKit(opts ...ComfyKitOption) *ComfyKit
```

Create a new ComfyKit instance with optional configuration.

### Options

```go
// ComfyUI configuration
func WithComfyUIBaseURL(url string) ComfyKitOption
func WithExecutorType(t ExecutorType) ComfyKitOption
func WithAPIKey(key string) ComfyKitOption
func WithCookies(cookies string) ComfyKitOption

// RunningHub configuration
func WithRunningHubBaseURL(url string) ComfyKitOption
func WithRunningHubAPIKey(key string) ComfyKitOption
func WithRunningHubTimeout(timeout int) ComfyKitOption
func WithRunningHubRetry(retry int) ComfyKitOption
func WithRunningHubInstance(instance string) ComfyKitOption
```

### Methods

```go
// Execute a workflow
func (k *ComfyKit) Execute(workflow string, params map[string]interface{}) (*ExecuteResult, error)

// Execute workflow from JSON
func (k *ComfyKit) ExecuteJSON(workflowJSON map[string]interface{}, params map[string]interface{}) (*ExecuteResult, error)

// Close all resources
func (k *ComfyKit) Close() error

// Get configuration
func (k *ComfyKit) GetComfyUIBaseURL() string
func (k *ComfyKit) GetExecutorType() ExecutorType
func (k *ComfyKit) GetAPIKey() string
func (k *ComfyKit) GetRunningHubBaseURL() string
func (k *ComfyKit) GetRunningHubAPIKey() string
func (k *ComfyKit) GetRunningHubTimeout() int
```

## ExecuteResult

```go
type ExecuteResult struct {
    Status       string                 // Execution status: "completed", "error", "timeout", "processing"
    PromptID     string                 // Task/prompt ID
    Duration     float64                // Execution duration in seconds
    
    Images       []string               // Image URLs
    ImagesByVar  map[string][]string    // Images grouped by variable name
    
    Audios       []string               // Audio URLs
    AudiosByVar  map[string][]string    // Audios grouped by variable name
    
    Videos       []string               // Video URLs
    VideosByVar  map[string][]string    // Videos grouped by variable name
    
    Texts        []string               // Text contents
    TextsByVar   map[string][]string    // Texts grouped by variable name
    
    Outputs      map[string]interface{} // Raw output data
    Message      string                 // Error message (if any)
}
```

### Methods

```go
func NewExecuteResult(status string) *ExecuteResult
func (r *ExecuteResult) ErrorResult(msg string) *ExecuteResult
func (r *ExecuteResult) TimeoutResult(duration float64) *ExecuteResult
```

## Executor Types

```go
type ExecutorType string

const (
    ExecutorTypeHTTP      ExecutorType = "http"
    ExecutorTypeWebSocket ExecutorType = "websocket"
)
```

## Executor Interface

```go
type Executor interface {
    ExecuteWorkflow(workflowFile string, params map[string]interface{}) (*ExecuteResult, error)
}
```

## Workflow Metadata

```go
type WorkflowMetadata struct {
    Title      string
    Params     map[string]WorkflowParam
    MappingInfo WorkflowMappingInfo
    WorkflowID string
    IsRunningHub bool
}

type WorkflowParam struct {
    Name      string
    Type      string
    Required  bool
    Default   interface{}
    NeedUpload bool
}
```

## Error Handling

All methods return an error when something goes wrong:

```go
result, err := kit.Execute("workflow.json", params)
if err != nil {
    // Handle connection errors, network issues, etc.
    fmt.Printf("Error: %v\n", err)
    return
}

if result.Status == "error" {
    // Handle workflow execution errors
    fmt.Printf("Workflow error: %s\n", result.Message)
}
```

# API 参�?

## ComfyKit

### 构造函�?

```go
func NewComfyKit(opts ...ComfyKitOption) *ComfyKit
```

创建一个新�?ComfyKit 实例，支持可选配置�?

### 选项

```go
// ComfyUI 配置
func WithComfyUIBaseURL(url string) ComfyKitOption
func WithExecutorType(t ExecutorType) ComfyKitOption
func WithAPIKey(key string) ComfyKitOption
func WithCookies(cookies string) ComfyKitOption

// RunningHub 配置
func WithRunningHubBaseURL(url string) ComfyKitOption
func WithRunningHubAPIKey(key string) ComfyKitOption
func WithRunningHubTimeout(timeout int) ComfyKitOption
func WithRunningHubRetry(retry int) ComfyKitOption
func WithRunningHubInstance(instance string) ComfyKitOption
```

### 方法

```go
// 执行工作�?
func (k *ComfyKit) Execute(workflow string, params map[string]interface{}) (*ExecuteResult, error)

// �?JSON 执行工作�?
func (k *ComfyKit) ExecuteJSON(workflowJSON map[string]interface{}, params map[string]interface{}) (*ExecuteResult, error)

// 关闭所有资�?
func (k *ComfyKit) Close() error

// 获取配置
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
    Status       string                 // 执行状�? "completed", "error", "timeout", "processing"
    PromptID     string                 // 任务/提示�?ID
    Duration     float64                // 执行时长（秒�?
    
    Images       []string               // 图片 URL 列表
    ImagesByVar  map[string][]string    // 按变量名分组的图�?
    
    Audios       []string               // 音频 URL 列表
    AudiosByVar  map[string][]string    // 按变量名分组的音�?
    
    Videos       []string               // 视频 URL 列表
    VideosByVar  map[string][]string    // 按变量名分组的视�?
    
    Texts        []string               // 文本内容列表
    TextsByVar   map[string][]string    // 按变量名分组的文�?
    
    Outputs      map[string]interface{} // 原始输出数据
    Message      string                 // 错误消息（如有）
}
```

### 方法

```go
func NewExecuteResult(status string) *ExecuteResult
func (r *ExecuteResult) ErrorResult(msg string) *ExecuteResult
func (r *ExecuteResult) TimeoutResult(duration float64) *ExecuteResult
```

## 执行器类�?

```go
type ExecutorType string

const (
    ExecutorTypeHTTP      ExecutorType = "http"
    ExecutorTypeWebSocket ExecutorType = "websocket"
)
```

## 执行器接�?

```go
type Executor interface {
    ExecuteWorkflow(workflowFile string, params map[string]interface{}) (*ExecuteResult, error)
}
```

## 工作流元数据

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

## 错误处理

所有方法在出错时返回错误：

```go
result, err := kit.Execute("workflow.json", params)
if err != nil {
    // 处理连接错误、网络问题等
    fmt.Printf("错误: %v\n", err)
    return
}

if result.Status == "error" {
    // 处理工作流执行错�?
    fmt.Printf("工作流错�? %s\n", result.Message)
}
```

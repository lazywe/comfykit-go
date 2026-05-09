# API 参考

## ComfyKit

### NewComfyKit

创建一个新的 ComfyKit 实例。

```go
func NewComfyKit(opts ...ComfyKitOption) *ComfyKit
```

**参数**：
- `opts` - 可变参数，用于配置 ComfyKit 实例

**返回值**：
- `*ComfyKit` - ComfyKit 实例

### Execute

执行工作流。

```go
func (k *ComfyKit) Execute(workflow string, params map[string]interface{}) (*ExecuteResult, error)
```

**参数**：
- `workflow` - 工作流路径、URL 或 RunningHub ID
- `params` - 参数映射

**返回值**：
- `*ExecuteResult` - 执行结果
- `error` - 错误信息

### ExecuteJSON

从 JSON 数据执行工作流。

```go
func (k *ComfyKit) ExecuteJSON(workflowJSON map[string]interface{}, params map[string]interface{}) (*ExecuteResult, error)
```

**参数**：
- `workflowJSON` - 工作流 JSON 数据
- `params` - 参数映射

**返回值**：
- `*ExecuteResult` - 执行结果
- `error` - 错误信息

### ExecuteAsyncByID

异步创建 RunningHub 任务，不等待执行完成，立即返回任务 ID。

```go
func (k *ComfyKit) ExecuteAsyncByID(workflowID string, params map[string]interface{}) (string, map[string]string, error)
```

**参数**：
- `workflowID` - RunningHub 工作流 ID
- `params` - 参数映射

**返回值**：
- `string` - 任务 ID
- `map[string]string` - 输出节点 ID 到变量名的映射
- `error` - 错误信息

### GetTaskCompletion

检查任务状态，当任务完成时返回结果。

```go
func (k *ComfyKit) GetTaskCompletion(taskID string, outputID2Var map[string]string) (*ExecuteResult, bool, error)
```

**参数**：
- `taskID` - 任务 ID
- `outputID2Var` - 输出节点 ID 到变量名的映射（从 ExecuteAsyncByID 获取）

**返回值**：
- `*ExecuteResult` - 执行结果（仅在任务完成时有效）
- `bool` - 是否完成
- `error` - 错误信息

### Close

清理资源。

```go
func (k *ComfyKit) Close() error
```

**返回值**：
- `error` - 错误信息

## ExecuteResult

执行结果结构体。

### 字段

| 字段 | 类型 | 描述 |
|------|------|------|
| `Status` | `string` | 执行状态 (completed/error/timeout) |
| `PromptID` | `string` | 提示词/任务 ID |
| `Duration` | `float64` | 执行时长（秒） |
| `Images` | `[]string` | 图片 URL 列表 |
| `ImagesByVar` | `map[string][]string` | 按变量名分组的图片 |
| `Videos` | `[]string` | 视频 URL 列表 |
| `Audios` | `[]string` | 音频 URL 列表 |
| `Texts` | `[]string` | 文本输出列表 |
| `Message` | `string` | 失败时的错误消息 |

## ComfyKitOption

配置选项函数。

### WithComfyUIBaseURL

设置 ComfyUI 基础 URL。

```go
func WithComfyUIBaseURL(url string) ComfyKitOption
```

### WithExecutorType

设置执行器类型。

```go
func WithExecutorType(t ExecutorType) ComfyKitOption
```

### WithAPIKey

设置 ComfyUI API 密钥。

```go
func WithAPIKey(key string) ComfyKitOption
```

### WithCookies

设置 ComfyUI cookies。

```go
func WithCookies(cookies string) ComfyKitOption
```

### WithRunningHubAPIKey

设置 RunningHub API 密钥。

```go
func WithRunningHubAPIKey(key string) ComfyKitOption
```

### WithRunningHubTimeout

设置 RunningHub 超时时间。

```go
func WithRunningHubTimeout(seconds int) ComfyKitOption
```

### WithRunningHubInstance

设置 RunningHub 实例类型。

```go
func WithRunningHubInstance(instance string) ComfyKitOption
```

## ExecutorType

执行器类型枚举。

```go
type ExecutorType string

const (
    ExecutorTypeHTTP      ExecutorType = "http"
    ExecutorTypeWebSocket ExecutorType = "websocket"
)
```

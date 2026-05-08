# 配置

ComfyKit-Go 支持多种配置方式，优先级：构造函数参�?> 环境变量 > 默认值�?
## 配置选项

### 本地 ComfyUI 选项

| 参数 | 环境变量 | 默认�?| 描述 |
|------|---------|--------|------|
| `comfyui_url` | `COMFYUI_BASE_URL` | `http://127.0.0.1:8188` | ComfyUI 服务器地址 |
| `executor_type` | `COMFYUI_EXECUTOR_TYPE` | `http` | 执行器类型：`http` �?`websocket` |
| `api_key` | `COMFYUI_API_KEY` | - | ComfyUI API 密钥 |
| `cookies` | `COMFYUI_COOKIES` | - | ComfyUI cookies |

### RunningHub 选项

| 参数 | 环境变量 | 默认�?| 描述 |
|------|---------|--------|------|
| `runninghub_url` | `RUNNINGHUB_BASE_URL` | `https://www.runninghub.ai` | RunningHub API 地址 |
| `runninghub_api_key` | `RUNNINGHUB_API_KEY` | - | RunningHub API 密钥 |
| `runninghub_timeout` | `RUNNINGHUB_TIMEOUT` | 0（无超时�?| 任务超时时间（秒�?|
| `runninghub_retry_count` | `RUNNINGHUB_RETRY_COUNT` | 3 | API 重试次数 |
| `runninghub_instance_type` | `RUNNINGHUB_INSTANCE_TYPE` | - | 实例类型 |

## 配置方法

### 方法 1：构造函数参�?
```go
kit := comfykit.NewComfyKit(
    comfykit.WithComfyUIBaseURL("http://my-server:8188"),
    comfykit.WithExecutorType(comfykit.ExecutorTypeWebSocket),
    comfykit.WithAPIKey("my-api-key"),
    comfykit.WithRunningHubAPIKey("rh-api-key"),
)
```

### 方法 2：环境变�?
```bash
export COMFYUI_BASE_URL=http://my-server:8188
export COMFYUI_API_KEY=my-api-key
export RUNNINGHUB_API_KEY=rh-api-key
```

```go
kit := comfykit.NewComfyKit()  // 自动从环境变量读�?```

### 方法 3：默认配�?
```go
kit := comfykit.NewComfyKit()
// 使用默认�?// comfyui_url: http://127.0.0.1:8188
// executor_type: http
```

## 配置优先�?
1. **构造函数参�?*（最高优先级�?2. **环境变量**
3. **默认�?*（最低优先级�?
## 多个实例

您可以创建多个独立配置的实例�?
```go
kitLocal := comfykit.NewComfyKit(comfykit.WithComfyUIBaseURL("http://127.0.0.1:8188"))
kitRemote := comfykit.NewComfyKit(comfykit.WithComfyUIBaseURL("http://remote-server:8188"))
kitCloud := comfykit.NewComfyKit(comfykit.WithRunningHubAPIKey("my-key"))
```

## 辅助函数

```go
kit := comfykit.NewComfyKit()

// 获取当前配置
fmt.Println(kit.GetComfyUIBaseURL())    // 获取 ComfyUI URL
fmt.Println(kit.GetExecutorType())       // 获取执行器类�?fmt.Println(kit.GetAPIKey())             // 获取 API 密钥
fmt.Println(kit.GetRunningHubAPIKey())   // 获取 RunningHub API 密钥
fmt.Println(kit.GetRunningHubTimeout())  // 获取超时时间
```

# 配置

## 配置选项

ComfyKit-Go 提供多种配置选项，可以通过函数参数或环境变量来配置。

### 本地 ComfyUI 配置

```go
kit := comfykit.NewComfyKit(
    comfykit.WithComfyUIBaseURL("http://localhost:8188"),
    comfykit.WithExecutorType(comfykit.ExecutorTypeHTTP),
    comfykit.WithAPIKey("your-api-key"),
    comfykit.WithCookies("session=abc123"),
)
```

### RunningHub 云端配置

```go
kit := comfykit.NewComfyKit(
    comfykit.WithRunningHubAPIKey("your-runninghub-api-key"),
    comfykit.WithRunningHubTimeout(300),
    comfykit.WithRunningHubInstance("plus"),
)
```

## 环境变量

以下环境变量可以用来配置 ComfyKit-Go：

| 变量 | 描述 | 默认值 |
|------|------|--------|
| `COMFYUI_BASE_URL` | ComfyUI 服务器地址 | `http://127.0.0.1:8188` |
| `COMFYUI_EXECUTOR_TYPE` | 执行类型 (http/websocket) | `http` |
| `COMFYUI_API_KEY` | ComfyUI API 密钥 | - |
| `COMFYUI_COOKIES` | ComfyUI cookies | - |
| `RUNNINGHUB_BASE_URL` | RunningHub API 地址 | `https://www.runninghub.ai` |
| `RUNNINGHUB_API_KEY` | RunningHub API 密钥 | - |
| `RUNNINGHUB_TIMEOUT` | 任务超时时间（秒） | 0（无限制） |
| `RUNNINGHUB_RETRY_COUNT` | API 重试次数 | 3 |
| `RUNNINGHUB_INSTANCE_TYPE` | 实例类型 (plus) | - |

## 配置优先级

配置的优先级从高到低：

1. 函数参数（如 `WithComfyUIBaseURL`）
2. 环境变量
3. 默认值

## 示例

### 使用环境变量配置

```bash
# Linux/macOS
export COMFYUI_BASE_URL=http://localhost:8188
export RUNNINGHUB_API_KEY=your-api-key

# Windows PowerShell
$env:COMFYUI_BASE_URL='http://localhost:8188'
$env:RUNNINGHUB_API_KEY='your-api-key'
```

### 使用代码配置

```go
kit := comfykit.NewComfyKit(
    comfykit.WithComfyUIBaseURL("http://localhost:8188"),
    comfykit.WithRunningHubAPIKey("your-api-key"),
    comfykit.WithRunningHubTimeout(300),
)
```

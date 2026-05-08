# Configuration

ComfyKit-Go can be configured in multiple ways with priority: constructor parameters > environment variables > defaults.

## Configuration Options

### Local ComfyUI Options

| Parameter | Environment Variable | Default | Description |
|-----------|---------------------|---------|-------------|
| `comfyui_url` | `COMFYUI_BASE_URL` | `http://127.0.0.1:8188` | ComfyUI server URL |
| `executor_type` | `COMFYUI_EXECUTOR_TYPE` | `http` | Executor type: `http` or `websocket` |
| `api_key` | `COMFYUI_API_KEY` | - | ComfyUI API key |
| `cookies` | `COMFYUI_COOKIES` | - | ComfyUI cookies |

### RunningHub Options

| Parameter | Environment Variable | Default | Description |
|-----------|---------------------|---------|-------------|
| `runninghub_url` | `RUNNINGHUB_BASE_URL` | `https://www.runninghub.ai` | RunningHub API URL |
| `runninghub_api_key` | `RUNNINGHUB_API_KEY` | - | RunningHub API key |
| `runninghub_timeout` | `RUNNINGHUB_TIMEOUT` | 0 (no timeout) | Task timeout in seconds |
| `runninghub_retry_count` | `RUNNINGHUB_RETRY_COUNT` | 3 | API retry count |
| `runninghub_instance_type` | `RUNNINGHUB_INSTANCE_TYPE` | - | Instance type |

## Configuration Methods

### Method 1: Constructor Parameters

```go
kit := comfykit.NewComfyKit(
    comfykit.WithComfyUIBaseURL("http://my-server:8188"),
    comfykit.WithExecutorType(comfykit.ExecutorTypeWebSocket),
    comfykit.WithAPIKey("my-api-key"),
    comfykit.WithRunningHubAPIKey("rh-api-key"),
)
```

### Method 2: Environment Variables

```bash
export COMFYUI_BASE_URL=http://my-server:8188
export COMFYUI_API_KEY=my-api-key
export RUNNINGHUB_API_KEY=rh-api-key
```

```go
kit := comfykit.NewComfyKit()  // Automatically reads from env vars
```

### Method 3: Default Configuration

```go
kit := comfykit.NewComfyKit()
// Uses default values
// comfyui_url: http://127.0.0.1:8188
// executor_type: http
```

## Configuration Priority

1. **Constructor parameters** (highest priority)
2. **Environment variables**
3. **Default values** (lowest priority)

## Multiple Instances

You can create multiple independent instances with different configurations:

```go
kitLocal := comfykit.NewComfyKit(comfykit.WithComfyUIBaseURL("http://127.0.0.1:8188"))
kitRemote := comfykit.NewComfyKit(comfykit.WithComfyUIBaseURL("http://remote-server:8188"))
kitCloud := comfykit.NewComfyKit(comfykit.WithRunningHubAPIKey("my-key"))
```

## Helper Functions

```go
kit := comfykit.NewComfyKit()

// Get current configuration
fmt.Println(kit.GetComfyUIBaseURL())    // Get ComfyUI URL
fmt.Println(kit.GetExecutorType())       // Get executor type
fmt.Println(kit.GetAPIKey())             // Get API key
fmt.Println(kit.GetRunningHubAPIKey())   // Get RunningHub API key
fmt.Println(kit.GetRunningHubTimeout())  // Get timeout
```

# Local Execution

Execute workflows on your local ComfyUI server.

## Prerequisites

1. Install and run ComfyUI:
   ```bash
   git clone https://github.com/comfyanonymous/ComfyUI.git
   cd ComfyUI
   python main.py
   ```

2. ComfyUI will be available at `http://127.0.0.1:8188`

## Basic Local Execution

```go
kit := comfykit.NewComfyKit()

result, err := kit.Execute("workflows/my_workflow.json", nil)
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

fmt.Printf("Status: %s\n", result.Status)
```

## Custom ComfyUI URL

```go
kit := comfykit.NewComfyKit(
    comfykit.WithComfyUIBaseURL("http://my-comfyui-server:8188"),
)

result, err := kit.Execute("workflows/my_workflow.json", nil)
```

## Executor Types

ComfyKit supports two executor types for local execution:

### HTTP Executor (Default)

```go
kit := comfykit.NewComfyKit(
    comfykit.WithExecutorType(comfykit.ExecutorTypeHTTP),
)
```

### WebSocket Executor

```go
kit := comfykit.NewComfyKit(
    comfykit.WithExecutorType(comfykit.ExecutorTypeWebSocket),
)
```

**When to use WebSocket:**
- Real-time progress updates
- Better for long-running workflows
- Lower latency for status updates

## Authentication

### API Key

```go
kit := comfykit.NewComfyKit(
    comfykit.WithAPIKey("my-comfyui-api-key"),
)
```

### Cookies

```go
kit := comfykit.NewComfyKit(
    comfykit.WithCookies("session=abc123; token=xyz"),
)
```

## Workflow Parameters

```go
params := map[string]interface{}{
    "prompt": "a beautiful sunset",
    "seed":   42,
    "steps":  25,
    "cfg":    7.0,
}

result, err := kit.Execute("workflows/t2i.json", params)
```

## Media Uploads

ComfyKit automatically handles media uploads for image/audio/video inputs:

```go
params := map[string]interface{}{
    "input_image": "/path/to/image.png",      // Local file
    "reference_image": "https://example.com/img.jpg",  // URL
}

result, err := kit.Execute("workflows/img2img.json", params)
```

## Local Workflow Examples

### Text-to-Image

```go
result, err := kit.Execute(
    "workflows/t2i.json",
    map[string]interface{}{
        "prompt": "a fantasy landscape with castles",
        "seed":   12345,
    },
)
```

### Image-to-Image

```go
result, err := kit.Execute(
    "workflows/img2img.json",
    map[string]interface{}{
        "input_image": "/path/to/input.jpg",
        "prompt":      "make this look like a painting",
        "denoise":     0.7,
    },
)
```

## Troubleshooting

### Connection Refused

```bash
# Make sure ComfyUI is running
python main.py --listen 0.0.0.0 --port 8188
```

### Authentication Errors

```go
kit := comfykit.NewComfyKit(
    comfykit.WithAPIKey("your-api-key"),
)
```

### Timeout Issues

```go
// WebSocket executor may help with long-running workflows
kit := comfykit.NewComfyKit(
    comfykit.WithExecutorType(comfykit.ExecutorTypeWebSocket),
)
```

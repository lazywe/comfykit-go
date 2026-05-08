# 本地执行

在本�?ComfyUI 服务器上执行工作流�?

## 前提条件

1. 安装并运�?ComfyUI�?
   ```bash
   git clone https://github.com/comfyanonymous/ComfyUI.git
   cd ComfyUI
   python main.py
   ```

2. ComfyUI 将在 `http://127.0.0.1:8188` 可用

## 基本本地执行

```go
kit := comfykit.NewComfyKit()

result, err := kit.Execute("workflows/my_workflow.json", nil)
if err != nil {
    fmt.Printf("错误: %v\n", err)
    return
}

fmt.Printf("状�? %s\n", result.Status)
```

## 自定�?ComfyUI URL

```go
kit := comfykit.NewComfyKit(
    comfykit.WithComfyUIBaseURL("http://my-comfyui-server:8188"),
)

result, err := kit.Execute("workflows/my_workflow.json", nil)
```

## 执行器类�?

ComfyKit 支持两种本地执行器类型：

### HTTP 执行器（默认�?

```go
kit := comfykit.NewComfyKit(
    comfykit.WithExecutorType(comfykit.ExecutorTypeHTTP),
)
```

### WebSocket 执行�?

```go
kit := comfykit.NewComfyKit(
    comfykit.WithExecutorType(comfykit.ExecutorTypeWebSocket),
)
```

**何时使用 WebSocket�?*
- 实时进度更新
- 长时间运行的工作�?
- 更低的状态更新延�?

## 认证

### API 密钥

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

## 工作流参�?

```go
params := map[string]interface{}{
    "prompt": "美丽的日�?,
    "seed":   42,
    "steps":  25,
    "cfg":    7.0,
}

result, err := kit.Execute("workflows/t2i.json", params)
```

## 媒体上传

ComfyKit 自动处理图像/音频/视频输入的媒体上传：

```go
params := map[string]interface{}{
    "input_image": "/path/to/image.png",      // 本地文件
    "reference_image": "https://example.com/img.jpg",  // URL
}

result, err := kit.Execute("workflows/img2img.json", params)
```

## 本地工作流示�?

### 文本转图�?

```go
result, err := kit.Execute(
    "workflows/t2i.json",
    map[string]interface{}{
        "prompt": "带有城堡的幻想风�?,
        "seed":   12345,
    },
)
```

### 图像转图�?

```go
result, err := kit.Execute(
    "workflows/img2img.json",
    map[string]interface{}{
        "input_image": "/path/to/input.jpg",
        "prompt":      "让它看起来像一幅画",
        "denoise":     0.7,
    },
)
```

## 故障排除

### 连接被拒�?

```bash
# 确保 ComfyUI 正在运行
python main.py --listen 0.0.0.0 --port 8188
```

### 认证错误

```go
kit := comfykit.NewComfyKit(
    comfykit.WithAPIKey("your-api-key"),
)
```

### 超时问题

```go
// WebSocket 执行器可能有助于长时间运行的工作�?
kit := comfykit.NewComfyKit(
    comfykit.WithExecutorType(comfykit.ExecutorTypeWebSocket),
)
```

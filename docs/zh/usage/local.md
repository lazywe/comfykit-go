# 本地执行

## 前提条件

在使用本地执行功能之前，需要：

1. 安装并运行 ComfyUI 服务器
2. 确保 ComfyUI 服务器可以访问（默认地址：`http://127.0.0.1:8188`）

## 基本用法

### 执行本地工作流文件

```go
import "github.com/lazywe/comfykit-go"

kit := comfykit.NewComfyKit()

result, err := kit.Execute("workflow.json", map[string]interface{}{
    "prompt": "一只可爱的猫",
    "seed":   42,
})
if err != nil {
    fmt.Printf("错误: %v\n", err)
    return
}
```

## 配置 ComfyUI 服务器

### 设置自定义地址

```go
kit := comfykit.NewComfyKit(
    comfykit.WithComfyUIBaseURL("http://localhost:8188"),
)
```

### 使用环境变量

```bash
# Linux/macOS
export COMFYUI_BASE_URL=http://localhost:8188

# Windows PowerShell
$env:COMFYUI_BASE_URL='http://localhost:8188'
```

## 执行器类型

### HTTP 执行器

```go
kit := comfykit.NewComfyKit(
    comfykit.WithExecutorType(comfykit.ExecutorTypeHTTP),
)
```

### WebSocket 执行器

```go
kit := comfykit.NewComfyKit(
    comfykit.WithExecutorType(comfykit.ExecutorTypeWebSocket),
)
```

## 认证配置

### API Key

```go
kit := comfykit.NewComfyKit(
    comfykit.WithAPIKey("your-api-key"),
)
```

### Cookies

```go
kit := comfykit.NewComfyKit(
    comfykit.WithCookies("session=abc123"),
)
```

## 完整示例

```go
package main

import (
    "fmt"
    "github.com/lazywe/comfykit-go"
)

func main() {
    // 创建 ComfyKit 实例，配置本地 ComfyUI 服务器
    kit := comfykit.NewComfyKit(
        comfykit.WithComfyUIBaseURL("http://localhost:8188"),
        comfykit.WithExecutorType(comfykit.ExecutorTypeHTTP),
        comfykit.WithAPIKey("your-api-key"),
    )
    
    // 定义参数
    params := map[string]interface{}{
        "prompt": "一只可爱的猫",
        "seed":   42,
        "steps":  20,
    }
    
    // 执行工作流
    result, err := kit.Execute("workflows/t2i_by_local_flux.json", params)
    if err != nil {
        fmt.Printf("执行错误: %v\n", err)
        return
    }
    
    // 处理结果
    fmt.Printf("状态: %s\n", result.Status)
    fmt.Printf("执行时长: %.2f 秒\n", result.Duration)
    fmt.Printf("生成图片: %d 张\n", len(result.Images))
    
    for i, imageURL := range result.Images {
        fmt.Printf("图片 %d: %s\n", i+1, imageURL)
    }
}
```

## 故障排除

### ComfyUI 服务器未运行

如果收到 "connection refused" 错误，请确保 ComfyUI 服务器正在运行：

```bash
# 启动 ComfyUI
python main.py
```

### 权限问题

如果收到权限错误，请检查 API Key 或 Cookies 是否正确配置。

### 工作流文件不存在

确保工作流文件路径正确：

```go
// 使用绝对路径
result, err := kit.Execute("/path/to/workflow.json", params)

// 或使用相对路径（相对于当前工作目录）
result, err := kit.Execute("workflows/workflow.json", params)
```

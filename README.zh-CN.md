# ComfyKit-Go

[![](https://img.shields.io/badge/lang-en-blue)](README.md) [![](https://img.shields.io/badge/lang-zh--CN-red)](README.zh-CN.md)
[![Experimental](https://img.shields.io/badge/status-experimental-yellow)](README.zh-CN.md)

**⚠️ 实验性产�?*: 这是 ComfyKit-python 的实验�?Go 语言移植版本。API 可能会在没有通知的情况下更改�?
ComfyKit-Go 是一个用于执�?ComfyUI 工作流的 Go SDK。它提供了一个简单、符�?Go 语言习惯�?API，用于在本地�?RunningHub 云平台上运行 ComfyUI 工作流�?
本项目是 [ComfyKit-python](https://github.com/runninghubai/comfykit-python) �?**Go 语言版本**，保持相同的目录结构、函数名称和 API 设计�?
## 特�?
- **本地执行**: 通过 HTTP �?WebSocket 在本�?ComfyUI 服务器上执行工作�?- **云端执行**: �?RunningHub 云平台上运行工作�?- **自动检�?*: 自动检测工作流类型（本地文件、URL、RunningHub ID�?- **参数映射**: 支持 ComfyKit DSL 参数语法
- **媒体处理**: 内置图片、音频和视频的媒体上传支�?- **跨平�?*: 支持 Windows、macOS �?Linux

## 安装

```bash
go get github.com/lazywe/comfykit-go
```

## 快速开�?
```go
package main

import (
    "fmt"
    "github.com/lazywe/comfykit-go"
)

func main() {
    // 使用默认配置创建 ComfyKit 实例
    kit := comfykit.NewComfyKit()
    
    // 执行本地工作�?    result, err := kit.Execute("workflow.json", map[string]interface{}{
        "prompt": "一只可爱的�?,
        "seed":   42,
    })
    if err != nil {
        fmt.Printf("错误: %v\n", err)
        return
    }
    
    fmt.Printf("状�? %s\n", result.Status)
    fmt.Printf("图片: %v\n", result.Images)
}
```

## 配置

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

## 执行模式

### 本地工作流文�?
```go
result, err := kit.Execute("workflow.json", map[string]interface{}{
    "prompt": "美丽的日�?,
})
```

### RunningHub 工作�?ID

```go
result, err := kit.Execute("12345", map[string]interface{}{
    "prompt": "美丽的日�?,
})
```

### 远程工作�?URL

```go
result, err := kit.Execute("https://example.com/workflow.json", map[string]interface{}{
    "prompt": "美丽的日�?,
})
```

## 工作�?DSL

ComfyKit 支持在工作流标题中定义参数的简�?DSL�?
- `$param` - 基本参数
- `$param!` - 必填参数
- `$~param` - 需要媒体上传的参数
- `$param.field` - 映射到特定字段的参数
- `$output.varname` - 输出变量�?
示例节点标题：`"Prompt, $prompt!, $seed"`

## 环境变量

| 变量 | 描述 | 默认�?|
|------|------|--------|
| `COMFYUI_BASE_URL` | ComfyUI 服务器地址 | `http://127.0.0.1:8188` |
| `COMFYUI_EXECUTOR_TYPE` | 执行类型 (http/websocket) | `http` |
| `COMFYUI_API_KEY` | ComfyUI API 密钥 | - |
| `COMFYUI_COOKIES` | ComfyUI cookies | - |
| `RUNNINGHUB_BASE_URL` | RunningHub API 地址 | `https://www.runninghub.ai` |
| `RUNNINGHUB_API_KEY` | RunningHub API 密钥 | - |
| `RUNNINGHUB_TIMEOUT` | 任务超时时间（秒�?| 0（无限制�?|
| `RUNNINGHUB_RETRY_COUNT` | API 重试次数 | 3 |
| `RUNNINGHUB_INSTANCE_TYPE` | 实例类型 (plus) | - |

## API 参�?
### ComfyKit

- `NewComfyKit(opts ...ComfyKitOption) *ComfyKit` - 创建新的 ComfyKit 实例
- `Execute(workflow string, params map[string]interface{}) (*ExecuteResult, error)` - 执行工作�?- `ExecuteJSON(workflowJSON map[string]interface{}, params map[string]interface{}) (*ExecuteResult, error)` - �?JSON 执行工作�?- `Close() error` - 清理资源

### ExecuteResult

- `Status` - 执行状�?(completed/error/timeout)
- `PromptID` - 提示�?任务 ID
- `Duration` - 执行时长（秒�?- `Images` - 图片 URL 列表
- `ImagesByVar` - 按变量名分组的图�?- `Videos` - 视频 URL 列表
- `Audios` - 音频 URL 列表
- `Texts` - 文本输出列表
- `Message` - 失败时的错误消息

## 致谢

本项目受�?[ComfyKit-python](https://github.com/runninghubai/comfykit-python) 的启发并基于其构建。我们要感谢原项目的出色设计和实现�?
## 许可�?
MIT License

# ComfyKit-Go

ComfyKit-Go 是一个用于执行 ComfyUI 工作流的 Go SDK。它提供了简单、符合 Go 语言习惯的 API，用于在本地和 RunningHub 云平台上运行 ComfyUI 工作流。

## 特性

- **本地执行**: 通过 HTTP 或 WebSocket 在本地 ComfyUI 服务器上执行工作流
- **云端执行**: 在 RunningHub 云平台上运行工作流
- **自动检测**: 自动检测工作流类型（本地文件、URL、RunningHub ID）
- **参数映射**: 支持 ComfyKit DSL 参数语法
- **媒体处理**: 内置图片、音频和视频的媒体上传支持
- **跨平台**: 支持 Windows、macOS 和 Linux

## 快速开始

```go
package main

import (
    "fmt"
    "github.com/lazywe/comfykit-go"
)

func main() {
    kit := comfykit.NewComfyKit()
    result, err := kit.Execute("workflow.json", map[string]interface{}{
        "prompt": "一只可爱的猫",
        "seed":   42,
    })
    if err != nil {
        fmt.Printf("错误: %v\n", err)
        return
    }
    fmt.Printf("状态: %s\n", result.Status)
}
```

## 目录

- [安装](installation.md)
- [快速开始](quick-start.md)
- [配置](configuration.md)
- [API 参考](api-reference.md)
- [示例](examples.md)
- [DSL 语法](dsl/overview.md)
- [开发指南](development.md)
- [贡献指南](contributing.md)
- [许可证](license.md)

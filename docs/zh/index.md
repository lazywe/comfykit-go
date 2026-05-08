# ComfyKit-Go

ComfyKit-Go 是一个用于执�?ComfyUI 工作流的 Go SDK，支持本地（HTTP/WebSocket）和云端（RunningHub）执行�?
## 特�?
- **统一 API**: 本地和云端执行使用单一接口
- **多种执行�?*: 支持 HTTP �?WebSocket 执行�?- **云端集成**: 无缝集成 RunningHub 云平�?- **参数映射**: 使用 DSL 语法轻松注入工作流参�?- **自动检�?*: 自动检测工作流源类型（文件/URL/RunningHub ID�?
## 快速开�?
```go
package main

import (
    "fmt"
    "github.com/lazywe/comfykit-go"
)

func main() {
    kit := comfykit.NewComfyKit()
    
    result, err := kit.Execute(
        "workflows/t2i.json",
        map[string]interface{}{"prompt": "美丽的日�?},
    )
    
    if err != nil {
        fmt.Printf("错误: %v\n", err)
        return
    }
    
    fmt.Printf("生成�?%d 张图片\n", len(result.Images))
}
```

## 文档

- [快速开始](quick-start.md) - 快速入�?- [安装](installation.md) - 安装说明
- [使用方法](usage/basic.md) - 基本使用示例
- [API 参考](api-reference.md) - 完整 API 文档

## 许可�?
MIT License

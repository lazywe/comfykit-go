# 快速开始

## 第一个程序

创建你的第一个 ComfyKit-Go 程序：

```go
package main

import (
    "fmt"
    "github.com/lazywe/comfykit-go"
)

func main() {
    // 创建 ComfyKit 实例
    kit := comfykit.NewComfyKit()
    
    // 执行工作流
    result, err := kit.Execute("workflow.json", map[string]interface{}{
        "prompt": "一只可爱的猫",
        "seed":   42,
    })
    if err != nil {
        fmt.Printf("错误: %v\n", err)
        return
    }
    
    // 处理结果
    fmt.Printf("状态: %s\n", result.Status)
    fmt.Printf("图片数量: %d\n", len(result.Images))
}
```

## 运行程序

```bash
go run main.go
```

## 工作流文件

工作流文件是一个 JSON 文件，包含 ComfyUI 工作流的定义。你可以在 ComfyUI 中导出工作流，然后使用 ComfyKit-Go 来执行它。

## 参数映射

ComfyKit-Go 支持简单的 DSL 语法来定义工作流参数：

- `$param` - 基本参数
- `$param!` - 必填参数
- `$~param` - 需要媒体上传的参数

示例工作流节点标题：`"Prompt, $prompt!, $seed"`

## 下一步

- 查看 [配置](configuration.md) 文档了解更多配置选项
- 查看 [API 参考](api-reference.md) 了解完整的 API
- 查看 [示例](examples.md) 了解更多使用示例

# 快速开�?
几分钟内开始使�?ComfyKit-Go�?
## 安装

```bash
go get github.com/lazywe/comfykit-go
```

## 基本用法

```go
package main

import (
    "fmt"
    "github.com/lazywe/comfykit-go"
)

func main() {
    // 创建 ComfyKit 实例
    kit := comfykit.NewComfyKit()
    
    // 执行工作�?    result, err := kit.Execute(
        "workflows/t2i.json",
        map[string]interface{}{
            "prompt": "海上美丽的日�?,
        },
    )
    
    if err != nil {
        fmt.Printf("错误: %v\n", err)
        return
    }
    
    // 检查结�?    fmt.Printf("状�? %s\n", result.Status)
    fmt.Printf("生成�?%d 张图片\n", len(result.Images))
    
    for i, img := range result.Images {
        fmt.Printf("图片 %d: %s\n", i+1, img)
    }
}
```

## 核心概念

### 工作流源

ComfyKit 自动检测工作流源：

1. **RunningHub ID**: 纯数字字符串（如 "12345"�?2. **URL**: �?`http://` �?`https://` 开�?3. **文件路径**: 本地文件路径

### 参数

在工作流 JSON 中使�?`$param` 语法定义参数�?
```json
{
  "text": "$prompt",
  "seed": "$seed"
}
```

执行时传递值：

```go
params := map[string]interface{}{
    "prompt": "我的提示�?,
    "seed":   12345,
}
```

## 下一�?
- [配置](configuration.md) - 了解配置选项
- [本地执行](usage/local.md) - 在本�?ComfyUI 上执行工作流
- [云端执行](usage/cloud.md) - �?RunningHub 上执行工作流

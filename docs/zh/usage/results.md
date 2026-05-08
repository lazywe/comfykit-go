# 结果处理

## 执行结果结构

执行工作流后，会返回一个 `ExecuteResult` 结构体：

```go
type ExecuteResult struct {
    Status      string
    PromptID    string
    Duration    float64
    Images      []string
    ImagesByVar map[string][]string
    Videos      []string
    Audios      []string
    Texts       []string
    Message     string
}
```

## 检查状态

### 状态类型

| 状态 | 描述 |
|------|------|
| `completed` | 执行成功 |
| `error` | 执行失败 |
| `timeout` | 执行超时 |

### 检查状态代码

```go
result, err := kit.Execute("workflow.json", params)
if err != nil {
    fmt.Printf("执行错误: %v\n", err)
    return
}

switch result.Status {
case "completed":
    fmt.Println("执行成功")
case "error":
    fmt.Printf("执行失败: %s\n", result.Message)
case "timeout":
    fmt.Println("执行超时")
}
```

## 获取输出

### 获取所有图片

```go
for _, imageURL := range result.Images {
    fmt.Printf("图片: %s\n", imageURL)
}
```

### 按变量名获取图片

```go
// 获取特定变量的图片
if images, ok := result.ImagesByVar["result"]; ok {
    for _, imageURL := range images {
        fmt.Printf("结果图片: %s\n", imageURL)
    }
}
```

### 获取视频

```go
for _, videoURL := range result.Videos {
    fmt.Printf("视频: %s\n", videoURL)
}
```

### 获取音频

```go
for _, audioURL := range result.Audios {
    fmt.Printf("音频: %s\n", audioURL)
}
```

### 获取文本

```go
for _, text := range result.Texts {
    fmt.Printf("文本: %s\n", text)
}
```

## 完整示例

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
        fmt.Printf("执行错误: %v\n", err)
        return
    }
    
    // 检查状态
    if result.Status != "completed" {
        fmt.Printf("执行失败，状态: %s\n", result.Status)
        if result.Message != "" {
            fmt.Printf("错误消息: %s\n", result.Message)
        }
        return
    }
    
    // 输出执行信息
    fmt.Printf("执行成功!\n")
    fmt.Printf("任务 ID: %s\n", result.PromptID)
    fmt.Printf("执行时长: %.2f 秒\n", result.Duration)
    
    // 输出图片
    fmt.Printf("\n生成的图片:\n")
    for i, imageURL := range result.Images {
        fmt.Printf("  %d. %s\n", i+1, imageURL)
    }
    
    // 输出按变量名分组的图片
    if len(result.ImagesByVar) > 0 {
        fmt.Printf("\n按变量分组的图片:\n")
        for varName, images := range result.ImagesByVar {
            fmt.Printf("  %s:\n", varName)
            for i, imageURL := range images {
                fmt.Printf("    %d. %s\n", i+1, imageURL)
            }
        }
    }
    
    // 输出视频
    if len(result.Videos) > 0 {
        fmt.Printf("\n生成的视频:\n")
        for i, videoURL := range result.Videos {
            fmt.Printf("  %d. %s\n", i+1, videoURL)
        }
    }
    
    // 输出音频
    if len(result.Audios) > 0 {
        fmt.Printf("\n生成的音频:\n")
        for i, audioURL := range result.Audios {
            fmt.Printf("  %d. %s\n", i+1, audioURL)
        }
    }
    
    // 输出文本
    if len(result.Texts) > 0 {
        fmt.Printf("\n生成的文本:\n")
        for i, text := range result.Texts {
            fmt.Printf("  %d. %s\n", i+1, text)
        }
    }
}
```

## 错误处理

### 网络错误

```go
result, err := kit.Execute("workflow.json", params)
if err != nil {
    if strings.Contains(err.Error(), "connection refused") {
        fmt.Println("无法连接到 ComfyUI 服务器")
    } else {
        fmt.Printf("网络错误: %v\n", err)
    }
    return
}
```

### 工作流错误

```go
if result.Status == "error" {
    fmt.Printf("工作流执行错误: %s\n", result.Message)
    
    // 根据错误消息进行处理
    if strings.Contains(result.Message, "参数缺失") {
        fmt.Println("请检查必填参数")
    }
}
```

## 超时处理

```go
// 设置超时时间
kit := comfykit.NewComfyKit(
    comfykit.WithRunningHubTimeout(300),
)

result, err := kit.Execute("workflow.json", params)
if err != nil {
    fmt.Printf("错误: %v\n", err)
    return
}

if result.Status == "timeout" {
    fmt.Println("执行超时，请增加超时时间或优化工作流")
}
```

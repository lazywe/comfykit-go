# 基础用法

## 创建 ComfyKit 实例

```go
import "github.com/lazywe/comfykit-go"

kit := comfykit.NewComfyKit()
```

## 执行工作流

### 从文件执行

```go
result, err := kit.Execute("workflow.json", map[string]interface{}{
    "prompt": "一只可爱的猫",
    "seed":   42,
})
if err != nil {
    fmt.Printf("错误: %v\n", err)
    return
}
```

### 从 URL 执行

```go
result, err := kit.Execute("https://example.com/workflow.json", map[string]interface{}{
    "prompt": "美丽的日落",
})
```

### 从 RunningHub ID 执行

```go
result, err := kit.Execute("12345", map[string]interface{}{
    "prompt": "科幻城市",
})
```

## 处理结果

### 检查状态

```go
if result.Status == "completed" {
    fmt.Println("执行成功")
} else if result.Status == "error" {
    fmt.Printf("执行失败: %s\n", result.Message)
} else if result.Status == "timeout" {
    fmt.Println("执行超时")
}
```

### 获取输出

```go
// 获取图片
for _, imageURL := range result.Images {
    fmt.Printf("图片: %s\n", imageURL)
}

// 获取视频
for _, videoURL := range result.Videos {
    fmt.Printf("视频: %s\n", videoURL)
}

// 获取文本
for _, text := range result.Texts {
    fmt.Printf("文本: %s\n", text)
}
```

### 按变量名获取输出

```go
if images, ok := result.ImagesByVar["result"]; ok {
    for _, imageURL := range images {
        fmt.Printf("结果图片: %s\n", imageURL)
    }
}
```

## 参数映射

### 基本参数

```go
params := map[string]interface{}{
    "prompt": "一只可爱的猫",
    "seed":   42,
    "steps":  30,
}
```

### 媒体参数

```go
params := map[string]interface{}{
    "image": "path/to/image.png",
}
```

### 必填参数

```go
// 如果工作流包含 $prompt!，则必须提供 prompt 参数
params := map[string]interface{}{
    "prompt": "必须提供的提示词",
}
```

## 配置选项

### 自定义配置

```go
kit := comfykit.NewComfyKit(
    comfykit.WithComfyUIBaseURL("http://localhost:8188"),
    comfykit.WithRunningHubAPIKey("your-api-key"),
)
```

### 使用环境变量

```bash
export COMFYUI_BASE_URL=http://localhost:8188
export RUNNINGHUB_API_KEY=your-api-key
```

## 完整示例

```go
package main

import (
    "fmt"
    "github.com/lazywe/comfykit-go"
)

func main() {
    // 创建 ComfyKit 实例
    kit := comfykit.NewComfyKit()
    
    // 定义参数
    params := map[string]interface{}{
        "prompt": "一只可爱的猫",
        "seed":   42,
        "steps":  20,
    }
    
    // 执行工作流
    result, err := kit.Execute("workflow.json", params)
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

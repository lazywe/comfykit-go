# 处理结果

学习如何处理执行结果�?
## 结果结构

```go
type ExecuteResult struct {
    Status       string                 // 执行状�?    PromptID     string                 // 任务 ID
    Duration     float64                // 执行时长（秒�?    
    Images       []string               // 所有图�?URL
    ImagesByVar  map[string][]string    // 按变量分组的图片
    
    Audios       []string               // 所有音�?URL
    AudiosByVar  map[string][]string    // 按变量分组的音频
    
    Videos       []string               // 所有视�?URL
    VideosByVar  map[string][]string    // 按变量分组的视频
    
    Texts        []string               // 所有文本输�?    TextsByVar   map[string][]string    // 按变量分组的文本
    
    Outputs      map[string]interface{} // 原始输出数据
    Message      string                 // 错误消息
}
```

## 状态�?
| 状�?| 描述 |
|------|------|
| `completed` | 执行成功完成 |
| `error` | 执行失败 |
| `timeout` | 执行超时 |
| `processing` | 执行进行�?|

## 访问结果

### 基本访问

```go
result, err := kit.Execute("workflow.json", params)
if err != nil {
    fmt.Printf("错误: %v\n", err)
    return
}

// 检查状�?if result.Status != "completed" {
    fmt.Printf("执行失败: %s\n", result.Message)
    return
}

// 访问所有图�?for i, img := range result.Images {
    fmt.Printf("图片 %d: %s\n", i+1, img)
}
```

### 按变量分�?
```go
// 按变量名访问图片
for varName, images := range result.ImagesByVar {
    fmt.Printf("\n变量: %s\n", varName)
    for i, img := range images {
        fmt.Printf("  %d. %s\n", i+1, img)
    }
}
```

### 不同媒体类型

```go
// 图片
fmt.Printf("图片: %d\n", len(result.Images))

// 音频
fmt.Printf("音频: %d\n", len(result.Audios))

// 视频
fmt.Printf("视频: %d\n", len(result.Videos))

// 文本
fmt.Printf("文本: %d\n", len(result.Texts))
```

### 原始输出

```go
// 访问原始输出数据
if result.Outputs != nil {
    for nodeID, output := range result.Outputs {
        fmt.Printf("节点 %s: %v\n", nodeID, output)
    }
}
```

## 错误处理

### 连接错误

```go
result, err := kit.Execute("workflow.json", params)
if err != nil {
    // 网络错误、认证问题等
    log.Printf("连接错误: %v", err)
    return
}
```

### 执行错误

```go
if result.Status == "error" {
    // 工作流执行失�?    log.Printf("执行错误: %s", result.Message)
    return
}
```

### 超时错误

```go
if result.Status == "timeout" {
    log.Printf("执行超时，耗时 %.2f �?, result.Duration)
    return
}
```

## 结果示例

### 单图片输�?
```go
result, err := kit.Execute("workflows/t2i.json", params)
if err != nil || result.Status != "completed" {
    return
}

if len(result.Images) > 0 {
    fmt.Printf("生成: %s\n", result.Images[0])
}
```

### 多个输出

```go
result, err := kit.Execute("workflows/multi_output.json", params)
if err != nil || result.Status != "completed" {
    return
}

// 按变量名访问输出
if images, ok := result.ImagesByVar["result"]; ok {
    for _, img := range images {
        fmt.Println(img)
    }
}

if images, ok := result.ImagesByVar["preview"]; ok {
    for _, img := range images {
        fmt.Println(img)
    }
}
```

### 文本输出

```go
result, err := kit.Execute("workflows/text_gen.json", params)
if err != nil || result.Status != "completed" {
    return
}

for _, text := range result.Texts {
    fmt.Println(text)
}
```

## 高级结果处理

### 下载图片

```go
import (
    "io"
    "net/http"
    "os"
)

func downloadImage(url, path string) error {
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    out, err := os.Create(path)
    if err != nil {
        return err
    }
    defer out.Close()

    _, err = io.Copy(out, resp.Body)
    return err
}

// 使用
for i, imgURL := range result.Images {
    err := downloadImage(imgURL, fmt.Sprintf("output_%d.png", i))
    if err != nil {
        log.Printf("下载 %s 失败: %v", imgURL, err)
    }
}
```

### 批处�?
```go
prompts := []string{"日落", "�?, "城市"}
for _, prompt := range prompts {
    result, err := kit.Execute("workflows/t2i.json", map[string]interface{}{
        "prompt": prompt,
    })
    if err != nil || result.Status != "completed" {
        continue
    }
    
    // 处理结果
    fmt.Printf("提示�? %s -> %d 张图片\n", prompt, len(result.Images))
}
```

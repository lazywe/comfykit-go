# 基本使用

学习 ComfyKit-Go 的基本用法�?

## 创建 ComfyKit 实例

```go
kit := comfykit.NewComfyKit()
```

使用自定义选项�?

```go
kit := comfykit.NewComfyKit(
    comfykit.WithComfyUIBaseURL("http://my-server:8188"),
    comfykit.WithAPIKey("my-secret-key"),
)
```

## 执行工作�?

### 从文�?

```go
result, err := kit.Execute("workflows/my_workflow.json", nil)
```

### �?URL

```go
result, err := kit.Execute("https://example.com/workflow.json", nil)
```

### �?RunningHub ID

```go
result, err := kit.Execute("12345", nil)  // RunningHub 工作�?ID
```

### 带参�?

```go
params := map[string]interface{}{
    "prompt": "美丽的风�?,
    "seed":   12345,
    "steps":  30,
}

result, err := kit.Execute("workflows/t2i.json", params)
```

### �?JSON

```go
workflowJSON := map[string]interface{}{
    "version": 1,
    "nodes": []interface{}{
        // ... 您的工作流节�?
    },
}

result, err := kit.ExecuteJSON(workflowJSON, params)
```

## 处理结果

```go
result, err := kit.Execute("workflow.json", params)
if err != nil {
    fmt.Printf("错误: %v\n", err)
    return
}

// 检查状�?
switch result.Status {
case "completed":
    fmt.Println("执行成功完成")
case "error":
    fmt.Printf("错误: %s\n", result.Message)
case "timeout":
    fmt.Println("执行超时")
}

// 访问输出
fmt.Printf("生成�?%d 张图片\n", len(result.Images))
fmt.Printf("生成�?%d 个音频\n", len(result.Audios))
fmt.Printf("生成�?%d 个视频\n", len(result.Videos))
fmt.Printf("生成�?%d 个文本\n", len(result.Texts))

// 按变量名访问
for varName, images := range result.ImagesByVar {
    fmt.Printf("%s: %d 张图片\n", varName, len(images))
}

// 获取执行时长
fmt.Printf("时长: %.2f 秒\n", result.Duration)
```

## 清理资源

完成后请关闭 ComfyKit 实例�?

```go
kit := comfykit.NewComfyKit()
defer kit.Close()

// ... 使用 kit ...
```

## 错误处理

```go
result, err := kit.Execute("workflow.json", params)

if err != nil {
    // 网络错误、连接问题等
    log.Fatalf("执行失败: %v", err)
}

if result.Status == "error" {
    // 工作流执行错�?
    log.Printf("工作流错�? %s", result.Message)
}
```

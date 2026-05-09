# 云端执行

## 前提条件

在使用云端执行功能之前，需要：

1. 在 RunningHub 平台注册账号
2. 获取 API Key
3. 设置环境变量或通过代码配置

## 基本用法

### 执行 RunningHub 工作流

```go
import "github.com/lazywe/comfykit-go"

// 设置 API Key
os.Setenv("RUNNINGHUB_API_KEY", "your-api-key")

kit := comfykit.NewComfyKit()

result, err := kit.Execute("12345", map[string]interface{}{
    "prompt": "一只可爱的猫",
    "seed":   42,
})
if err != nil {
    fmt.Printf("错误: %v\n", err)
    return
}
```

## 配置 RunningHub

### 通过代码配置

```go
kit := comfykit.NewComfyKit(
    comfykit.WithRunningHubAPIKey("your-api-key"),
    comfykit.WithRunningHubTimeout(300),
    comfykit.WithRunningHubInstance("plus"),
)
```

### 使用环境变量

```bash
# Linux/macOS
export RUNNINGHUB_API_KEY=your-api-key
export RUNNINGHUB_TIMEOUT=300

# Windows PowerShell
$env:RUNNINGHUB_API_KEY='your-api-key'
$env:RUNNINGHUB_TIMEOUT='300'
```

## 执行模式

### 按工作流 ID 执行

```go
result, err := kit.Execute("12345", map[string]interface{}{
    "prompt": "美丽的日落",
})
```

### 按工作流文件执行

```go
result, err := kit.Execute("workflows/my_runninghub_workflow.json", params)
```

### 自动检测

ComfyKit-Go 会自动检测工作流类型：

- **RunningHub ID**: 纯数字字符串（如 "12345"）
- **URL**: 以 http:// 或 https:// 开头
- **文件路径**: 包含 / 或 \
- **RunningHub 工作流文件**: 包含 `_source: runninghub`

## 完整示例

```go
package main

import (
    "fmt"
    "os"
    "github.com/lazywe/comfykit-go"
)

func main() {
    // 设置 API Key
    apiKey := os.Getenv("RUNNINGHUB_API_KEY")
    if apiKey == "" {
        fmt.Println("请设置 RUNNINGHUB_API_KEY 环境变量")
        return
    }
    
    // 创建 ComfyKit 实例
    kit := comfykit.NewComfyKit(
        comfykit.WithRunningHubAPIKey(apiKey),
        comfykit.WithRunningHubTimeout(300),
    )
    
    // 定义参数
    params := map[string]interface{}{
        "prompt": "一只可爱的猫",
        "seed":   42,
        "steps":  20,
    }
    
    // 执行工作流（使用 RunningHub ID）
    result, err := kit.Execute("12345", params)
    if err != nil {
        fmt.Printf("执行错误: %v\n", err)
        return
    }
    
    // 处理结果
    fmt.Printf("状态: %s\n", result.Status)
    fmt.Printf("执行时长: %.2f 秒\n", result.Duration)
    fmt.Printf("任务 ID: %s\n", result.PromptID)
    fmt.Printf("生成图片: %d 张\n", len(result.Images))
    
    for i, imageURL := range result.Images {
        fmt.Printf("图片 %d: %s\n", i+1, imageURL)
    }
}
```

## 故障排除

### API Key 未设置

确保设置了 `RUNNINGHUB_API_KEY` 环境变量：

```bash
export RUNNINGHUB_API_KEY=your-api-key
```

### 工作流 ID 不存在

检查工作流 ID 是否正确：

```go
result, err := kit.Execute("12345", params)
if err != nil {
    fmt.Printf("工作流不存在或无权访问: %v\n", err)
}
```

### 任务超时

增加超时时间：

```go
kit := comfykit.NewComfyKit(
    comfykit.WithRunningHubTimeout(600), // 10 分钟
)
```

## 异步执行

除了同步等待执行完成，ComfyKit-Go 还支持异步创建任务和轮询获取结果，这样可以更灵活地管理多个任务的执行流程。

### 异步创建任务

```go
// 异步创建任务，立即返回 taskID
taskID, outputID2Var, err := kit.ExecuteAsyncByID("12345", map[string]interface{}{
    "prompt": "一只可爱的猫",
})
if err != nil {
    fmt.Printf("创建任务失败: %v\n", err)
    return
}

fmt.Printf("任务创建成功，ID: %s\n", taskID)
fmt.Printf("输出节点: %v\n", outputID2Var)
```

### 轮询获取结果

```go
// 轮询检查任务状态
checkInterval := 2 * time.Second
maxChecks := 30 // 最多检查30次

for i := 0; i < maxChecks; i++ {
    result, completed, err := kit.GetTaskCompletion(taskID, outputID2Var)
    if err != nil {
        fmt.Printf("检查状态错误: %v\n", err)
        time.Sleep(checkInterval)
        continue
    }

    if completed {
        if result.Status == "completed" {
            fmt.Printf("执行完成，生成了 %d 张图片\n", len(result.Images))
            for _, imageURL := range result.Images {
                fmt.Printf("  - %s\n", imageURL)
            }
        } else {
            fmt.Printf("执行失败: %s - %s\n", result.Status, result.Message)
        }
        break
    }

    fmt.Printf("检查 %d/%d: 任务还在运行中...\n", i+1, maxChecks)
    time.Sleep(checkInterval)
}
```

### 多任务并发管理

```go
// 创建多个任务
tasks := []struct {
    Prompt string
    TaskID string
    Output map[string]string
}{
    {Prompt: "红色的玫瑰"},
    {Prompt: "蓝色的海洋"},
    {Prompt: "绿色的森林"},
}

// 批量创建任务
for i := range tasks {
    taskID, outputID2Var, err := kit.ExecuteAsyncByID("12345", map[string]interface{}{
        "prompt": tasks[i].Prompt,
    })
    if err != nil {
        fmt.Printf("创建任务 %d 失败: %v\n", i, err)
        continue
    }
    tasks[i].TaskID = taskID
    tasks[i].Output = outputID2Var
    fmt.Printf("任务 %d 创建成功: %s\n", i, taskID)
}

// 轮询所有任务
remaining := len(tasks)
for remaining > 0 {
    for i := range tasks {
        if tasks[i].TaskID == "" {
            continue
        }
        
        result, completed, err := kit.GetTaskCompletion(tasks[i].TaskID, tasks[i].Output)
        if err != nil {
            continue
        }
        
        if completed {
            remaining--
            tasks[i].TaskID = "" // 标记为已完成
            if result.Status == "completed" {
                fmt.Printf("任务 %d 完成: %d 张图片\n", i, len(result.Images))
            } else {
                fmt.Printf("任务 %d 失败: %s\n", i, result.Message)
            }
        }
    }
    
    time.Sleep(2 * time.Second)
}
```

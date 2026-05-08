# 云端执行

�?RunningHub 云平台上执行工作流�?
## 前提条件

1. 注册 RunningHub 账户：https://www.runninghub.ai
2. �?RunningHub 控制台获�?API 密钥

## 配置

```go
kit := comfykit.NewComfyKit(
    comfykit.WithRunningHubAPIKey("your-api-key"),
)
```

## �?RunningHub ID 执行

```go
result, err := kit.Execute("12345", nil)
if err != nil {
    fmt.Printf("错误: %v\n", err)
    return
}

fmt.Printf("状�? %s\n", result.Status)
fmt.Printf("生成�?%d 张图片\n", len(result.Images))
```

## 带参数执�?
```go
params := map[string]interface{}{
    "prompt": "未来城市",
    "seed":   12345,
    "steps":  30,
}

result, err := kit.Execute("12345", params)
```

## 执行 RunningHub 工作流文�?
您也可以执行包含 RunningHub 元数据的本地工作流文件：

```json
{
  "_source": "runninghub",
  "workflow_id": "12345",
  "params": {
    "prompt": "$prompt",
    "seed": "$seed"
  }
}
```

```go
result, err := kit.Execute("workflows/my_runninghub_workflow.json", params)
```

## 云端特定选项

### 超时

```go
kit := comfykit.NewComfyKit(
    comfykit.WithRunningHubAPIKey("your-api-key"),
    comfykit.WithRunningHubTimeout(600), // 10 分钟
)
```

### 实例类型

```go
kit := comfykit.NewComfyKit(
    comfykit.WithRunningHubAPIKey("your-api-key"),
    comfykit.WithRunningHubInstance("plus"),
)
```

可用的实例类型：
- `standard` - 默认实例
- `plus` - 增强性能
- `pro` - 高级性能
- `gpu` - GPU 加�?
### 重试次数

```go
kit := comfykit.NewComfyKit(
    comfykit.WithRunningHubAPIKey("your-api-key"),
    comfykit.WithRunningHubRetry(5),
)
```

## 混合执行

您可以使用同一�?ComfyKit 实例进行本地和云端执行：

```go
kit := comfykit.NewComfyKit(
    comfykit.WithRunningHubAPIKey("your-api-key"),
)

// 本地执行
localResult, _ := kit.Execute("workflows/local.json", params)

// 云端执行
cloudResult, _ := kit.Execute("12345", params)
```

## 自动检�?
ComfyKit 自动检测工作流源：

```go
kit := comfykit.NewComfyKit()

// 自动检�?RunningHub ID
result, _ := kit.Execute("12345", params)

// 自动检�?URL
result, _ := kit.Execute("https://example.com/workflow.json", params)

// 自动检测本地文�?result, _ := kit.Execute("workflows/local.json", params)
```

## 故障排除

### API 密钥未设�?
```bash
export RUNNINGHUB_API_KEY="your-api-key"
```

或通过编程方式设置�?
```go
kit := comfykit.NewComfyKit(
    comfykit.WithRunningHubAPIKey("your-api-key"),
)
```

### 任务超时

增加超时时间�?
```go
kit := comfykit.NewComfyKit(
    comfykit.WithRunningHubAPIKey("your-api-key"),
    comfykit.WithRunningHubTimeout(1200), // 20 分钟
)
```

### 工作流未找到

确保工作�?ID 存在于您�?RunningHub 账户中�?
## 云端优势

- 无需维护本地 GPU
- 访问强大�?GPU（NVIDIA A100, RTX 4090�?- 自动扩展
- 按需付费
- 内置工作流管�?
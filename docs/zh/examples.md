# 示例

ComfyKit-Go 提供了多个示例程序来帮助你学习如何使用这个 SDK。

## 示例列表

### 1. 快速入门

文件：`examples/01_quick_start.go`

这个示例展示了最基本的使用方法：

```go
kit := comfykit.NewComfyKit()
result, err := kit.Execute("workflow.json", map[string]interface{}{
    "prompt": "一只可爱的猫",
    "seed":   42,
})
```

### 2. 配置选项

文件：`examples/02_configuration.go`

这个示例展示了如何配置 ComfyKit：

```go
kit := comfykit.NewComfyKit(
    comfykit.WithComfyUIBaseURL("http://localhost:8188"),
    comfykit.WithExecutorType(comfykit.ExecutorTypeHTTP),
)
```

### 3. 本地工作流

文件：`examples/03_local_workflows.go`

这个示例展示了如何执行本地工作流文件。

### 4. RunningHub 云端

文件：`examples/04_runninghub_cloud.go`

这个示例展示了如何在 RunningHub 云平台上执行工作流。

### 5. 高级功能

文件：`examples/05_advanced_features.go`

这个示例展示了高级功能，如批量执行、错误处理和超时处理。

## 运行示例

```bash
# 运行单个示例
go run examples/01_quick_start.go

# 运行所有示例
go run examples/run_all.go
```

## 示例目录结构

```
examples/
├── 01_quick_start.go      # 快速入门
├── 02_configuration.go    # 配置选项
├── 03_local_workflows.go  # 本地工作流
├── 04_runninghub_cloud.go # RunningHub 云端
├── 05_advanced_features.go # 高级功能
├── README.md              # 示例说明
└── run_all.go             # 运行所有示例
```

## 注意事项

- 运行本地工作流需要 ComfyUI 服务器正在运行
- 运行云端工作流需要设置 `RUNNINGHUB_API_KEY` 环境变量
- 所有示例都包含注释说明，建议阅读代码了解详细用法

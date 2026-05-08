# 示例

ComfyKit-Go 包含多个示例帮助您入门�?
## 示例列表

### 01_quick_start.go
- 基本工作流执�?- 简单的 3 行示�?- 结果处理

### 02_configuration.go
- 默认配置
- 自定义参�?- 环境变量
- 配置优先�?- 多个实例

### 03_local_workflows.go
- 本地 ComfyUI 执行
- 自定义参�?- �?JSON 执行
- 不同输出类型
- WebSocket 执行�?
### 04_runninghub_cloud.go
- RunningHub 云端执行
- 带参数的云端执行
- 云端工作流文�?- 自动检�?- 混合本地+云端

### 05_advanced_features.go
- 批处�?- 错误处理
- 超时处理
- 认证
- 执行器类�?
## 运行示例

```bash
# 运行单个示例
go run examples/01_quick_start.go

# 运行所有示�?go run examples/run_all.go
```

## 示例结构

每个示例遵循以下模式�?
```go
/*
示例标题 - 简要描�?
这个示例教您什么的详细说明�?*/

package main

import (
    "fmt"
    "github.com/lazywe/comfykit-go"
)

func main() {
    // 创建 ComfyKit 实例
    kit := comfykit.NewComfyKit()
    
    // 执行工作�?    result, err := kit.Execute("workflow.json", params)
    
    // 处理结果
    if err != nil {
        fmt.Printf("错误: %v\n", err)
        return
    }
    
    fmt.Printf("状�? %s\n", result.Status)
}
```

## 创建自己的示�?
按照以下步骤创建自己的示例：

1. **�?`examples/` 目录创建新文�?*
2. **添加描述性头�?*说明示例的作�?3. **导入 ComfyKit** 和其他必要的�?4. **创建 ComfyKit 实例**并配置适当的选项
5. **执行工作�?*并传递参�?6. **处理结果**并打印有意义的输�?7. **添加注释**解释每个步骤

## 示例提示

1. **从简单开�?*: 从基本执行开�?2. **逐步测试**: 添加更多复杂性之前测试每个部�?3. **处理错误**: 始终优雅处理错误
4. **使用日志**: 添加日志用于调试
5. **文档�?*: 解释示例的作�?
## 贡献示例

如果您有有用的示例，考虑贡献�?
1. 遵循现有示例模式
2. 添加适当的注�?3. 包含错误处理
4. 在注释中添加文档
5. 提交 pull request

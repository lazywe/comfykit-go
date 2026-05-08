# 开�?
学习如何开发和贡献 ComfyKit-Go�?
## 设置开发环�?
```bash
# 克隆仓库
git clone https://github.com/lazywe/comfykit-go.git
cd comfykit-go

# 安装依赖
go mod download

# 构建项目
go build ./...

# 运行测试
go test ./...
```

## 项目结构

```
comfykit-go/
├── comfykit.go          # �?ComfyKit �?├── executor.go          # 执行器接�?├── logger.go            # 日志工具
├── comfyui/             # ComfyUI 相关代码
�?  ├── executor.go      # 基础执行�?�?  ├── http_executor.go # HTTP 执行�?�?  ├── websocket_executor.go # WebSocket 执行�?�?  ├── runninghub_client.go  # RunningHub API 客户�?�?  ├── runninghub_executor.go # RunningHub 执行�?�?  ├── models.go        # 数据模型
�?  └── workflow_parser.go # 工作流解析器
├── utils/               # 工具函数
�?  ├── config_util.go   # 配置工具
�?  ├── file_util.go     # 文件工具
�?  ├── network_util.go  # 网络工具
�?  └── ...
├── examples/            # 示例脚本
├── workflows/           # 示例工作�?└── docs/                # 文档
```

## 编码指南

### 风格
- 遵循 Go 约定（gofmt�?- 函数名使�?`camelCase`
- 类型和导出标识符使用 `PascalCase`
- 内部函数使用 `snake_case`

### 错误处理
- 显式返回错误
- 使用描述性错误消�?- 使用上下文包装错�?- 避免静默失败

### 测试
- 为新功能编写单元测试
- 遵循现有测试模式
- 提交更改前运行测�?- 测试边缘情况

### 文档
- 文档化导出的函数和类�?- 为复杂逻辑添加注释
- 添加功能时更新文�?- 使用清晰简洁的语言

## 添加新功�?
### 步骤 1：计�?- 定义功能需�?- 设计 API
- 考虑向后兼容�?
### 步骤 2：实�?- 编写代码
- 添加测试
- 更新文档

### 步骤 3：测�?- 运行现有测试
- 验证新功�?- 测试边缘情况

### 步骤 4：文档化
- 更新 API 参�?- 必要时添加示�?- 更新文档

## 调试

### 启用调试日志
```go
kit := comfykit.NewComfyKit()
// 使用详细日志
```

### 打印调试信息
```go
result, err := kit.Execute("workflow.json", params)
if err != nil {
    fmt.Printf("错误: %v\n", err)
    return
}

// 打印详细结果
fmt.Printf("结果: %+v\n", result)
```

### 使用日志
```go
import "log"

log.Printf("执行工作�? %s", workflowPath)
log.Printf("参数: %v", params)
log.Printf("结果状�? %s", result.Status)
```

## 测试

### 运行所有测�?```bash
go test ./...
```

### 运行特定测试
```bash
go test ./comfyui/...
go test ./utils/...
```

### 测试覆盖�?```bash
go test -cover ./...
```

## 发布

### 更新版本
- 更新 `go.mod` 中的版本
- 更新文档中的版本

### 创建发布说明
- 文档化更�?- 列出重大变更
- 提供迁移指南

### 发布
- 推送到 GitHub
- 创建发布标签
- 更新文档

## 贡献

### Fork 仓库
1. �?GitHub �?Fork 仓库
2. 克隆您的 fork
3. 创建功能分支

### 进行更改
1. 实现更改
2. 添加测试
3. 更新文档

### 提交 PR
1. 推送更�?2. 创建 pull request
3. 等待审核
4. 处理反馈

## 代码审核指南

### 审核要点
- 代码正确�?- 性能
- 安全�?- 可读�?- 测试覆盖�?- 文档

### 反馈
- 建设性反�?- 解释您的理由
- 建议改进
- 鼓励最佳实�?
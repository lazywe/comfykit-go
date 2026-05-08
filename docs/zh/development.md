# 开发指南

## 环境设置

### 安装依赖

```bash
go mod download
```

### 运行测试

```bash
go test ./...
```

### 构建项目

```bash
go build ./...
```

## 代码结构

```
ComfyKit-Go/
├── comfykit.go              # 主入口
├── comfykit_test.go         # 主模块测试
├── executor.go              # 执行器接口
├── logger.go                # 日志工具
├── go.mod                   # Go 模块配置
├── go.sum                   # 依赖锁定
├── comfyui/                 # ComfyUI 核心模块
│   ├── executor.go          # 执行器接口
│   ├── http_executor.go     # HTTP 执行器
│   ├── websocket_executor.go # WebSocket 执行器
│   ├── runninghub_client.go # RunningHub 客户端
│   ├── runninghub_executor.go # RunningHub 执行器
│   ├── workflow_parser.go   # 工作流解析器
│   └── models.go            # 数据模型
├── utils/                   # 工具函数模块
│   ├── config_util.go       # 配置工具
│   ├── dynamic_util.go      # 动态工具
│   ├── file_util.go         # 文件工具
│   ├── image_util.go        # 图片工具
│   ├── network_util.go      # 网络工具
│   ├── openapi_util.go      # OpenAPI 工具
│   ├── os_util.go           # 操作系统工具
│   ├── process_util.go      # 进程工具
│   ├── runninghub_util.go   # RunningHub 工具
│   └── workflow_source_util.go # 工作流源工具
├── examples/                # 示例代码
├── workflows/               # 工作流文件
└── docs/                    # 文档
```

## 开发流程

### 1. 创建分支

```bash
git checkout -b feature/your-feature-name
```

### 2. 编写代码

遵循现有的代码风格和模式。

### 3. 编写测试

为新功能编写单元测试。

### 4. 运行测试

```bash
go test ./...
```

### 5. 提交代码

```bash
git add .
git commit -m "feat: add your feature"
git push origin feature/your-feature-name
```

## 代码风格

- 使用 `gofmt` 格式化代码
- 使用 `golint` 检查代码风格
- 遵循 Go 语言的命名规范

## 版本发布

### 创建版本标签

```bash
git tag -a v1.0.0 -m "Version 1.0.0"
git push origin v1.0.0
```

## 贡献

请参考 [贡献指南](contributing.md)。

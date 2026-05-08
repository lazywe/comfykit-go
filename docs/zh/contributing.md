# 贡献指南

欢迎贡献代码！请遵循以下指南。

## 代码贡献流程

### 1. Fork 项目

在 GitHub 上 Fork 这个项目到你的账号。

### 2. 克隆项目

```bash
git clone https://github.com/your-username/comfykit-go.git
cd comfykit-go
```

### 3. 创建分支

```bash
git checkout -b feature/your-feature-name
```

### 4. 编写代码

- 遵循现有的代码风格
- 添加单元测试
- 更新文档

### 5. 运行测试

```bash
go test ./...
```

### 6. 提交代码

```bash
git add .
git commit -m "feat: add your feature"
git push origin feature/your-feature-name
```

### 7. 创建 Pull Request

在 GitHub 上创建 Pull Request。

## 代码风格

### Go 代码

- 使用 `gofmt` 格式化代码
- 使用 `golint` 检查代码风格
- 函数名使用 PascalCase
- 变量名使用 camelCase
- 常量名使用 UPPER_CASE

### 提交消息格式

使用 Conventional Commits 格式：

```
feat: add new feature
fix: fix a bug
docs: update documentation
refactor: refactor code
test: add tests
```

## 报告问题

如果你发现了 bug 或者有功能建议，请在 GitHub Issues 上提交。

## 联系方式

如有任何问题，请通过 GitHub Issues 联系我们。

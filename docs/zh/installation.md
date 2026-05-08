# 安装

使用 Go modules 安装 ComfyKit-Go�?
## 前提条件

- Go 1.21 或更高版�?
## 安装

```bash
go get github.com/lazywe/comfykit-go
```

## 导入

```go
import (
    "github.com/lazywe/comfykit-go"
    "github.com/lazywe/comfykit-go/comfyui"
)
```

## 验证安装

创建一个简单的测试文件�?
```go
package main

import (
    "fmt"
    "github.com/lazywe/comfykit-go"
)

func main() {
    kit := comfykit.NewComfyKit()
    fmt.Println("ComfyKit-Go 初始化成功！")
    fmt.Printf("连接�? %s\n", kit.GetComfyUIBaseURL())
}
```

运行它：

```bash
go run main.go
```

## 依赖

ComfyKit-Go 使用以下外部包：

- `github.com/gorilla/websocket` - WebSocket 支持
- 标准库包

## 开�?
要贡献或本地开发：

```bash
git clone https://github.com/lazywe/comfykit-go.git
cd comfykit-go
go build ./...
```

运行测试�?
```bash
go test ./...
```

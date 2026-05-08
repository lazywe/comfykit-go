# 安装

## 要求

- Go 1.22 或更高版本

## 安装方法

使用 `go get` 命令安装：

```bash
go get github.com/lazywe/comfykit-go
```

## 更新

```bash
go get -u github.com/lazywe/comfykit-go
```

## 验证安装

创建一个简单的测试文件来验证安装是否成功：

```go
package main

import (
    "fmt"
    "github.com/lazywe/comfykit-go"
)

func main() {
    kit := comfykit.NewComfyKit()
    fmt.Println("ComfyKit-Go 安装成功!")
}
```

运行测试：

```bash
go run main.go
```

如果输出 "ComfyKit-Go 安装成功!"，则安装成功。

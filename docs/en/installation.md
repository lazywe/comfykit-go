# Installation

Install ComfyKit-Go using Go modules.

## Prerequisites

- Go 1.21 or higher

## Install

```bash
go get github.com/lazywe/comfykit-go
```

## Import

```go
import (
    "github.com/lazywe/comfykit-go"
    "github.com/lazywe/comfykit-go/comfyui"
)
```

## Verify Installation

Create a simple test file:

```go
package main

import (
    "fmt"
    "github.com/lazywe/comfykit-go"
)

func main() {
    kit := comfykit.NewComfyKit()
    fmt.Println("ComfyKit-Go initialized successfully!")
    fmt.Printf("Connected to: %s\n", kit.GetComfyUIBaseURL())
}
```

Run it:

```bash
go run main.go
```

## Dependencies

ComfyKit-Go uses the following external packages:

- `github.com/gorilla/websocket` - WebSocket support
- Standard library packages only

## Development

To contribute or develop locally:

```bash
git clone https://github.com/lazywe/comfykit-go.git
cd comfykit-go
go build ./...
```

Run tests:

```bash
go test ./...
```

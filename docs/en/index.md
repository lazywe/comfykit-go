# ComfyKit-Go

ComfyKit-Go is a Go SDK for executing ComfyUI workflows with both local (HTTP/WebSocket) and cloud (RunningHub) capabilities.

## Features

- **Unified API**: Single interface for local and cloud execution
- **Multiple Executors**: HTTP and WebSocket support for local ComfyUI
- **Cloud Integration**: Seamless RunningHub cloud execution
- **Parameter Mapping**: Easy workflow parameter injection using DSL syntax
- **Auto-detection**: Automatically detect workflow source type (file/URL/RunningHub ID)

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/lazywe/comfykit-go"
)

func main() {
    kit := comfykit.NewComfyKit()
    
    result, err := kit.Execute(
        "workflows/t2i.json",
        map[string]interface{}{"prompt": "a beautiful sunset"},
    )
    
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Printf("Generated %d images\n", len(result.Images))
}
```

## Documentation

- [Quick Start](quick-start.md) - Get started in minutes
- [Installation](installation.md) - Installation instructions
- [Usage](usage/basic.md) - Basic usage examples
- [API Reference](api-reference.md) - Complete API documentation

## License

MIT License

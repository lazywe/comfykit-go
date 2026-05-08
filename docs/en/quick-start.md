# Quick Start

Get started with ComfyKit-Go in minutes!

## Installation

```bash
go get github.com/lazywe/comfykit-go
```

## Basic Usage

```go
package main

import (
    "fmt"
    "github.com/lazywe/comfykit-go"
)

func main() {
    // Create ComfyKit instance
    kit := comfykit.NewComfyKit()
    
    // Execute a workflow
    result, err := kit.Execute(
        "workflows/t2i.json",
        map[string]interface{}{
            "prompt": "a beautiful sunset over the ocean",
        },
    )
    
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    // Check results
    fmt.Printf("Status: %s\n", result.Status)
    fmt.Printf("Generated %d images\n", len(result.Images))
    
    for i, img := range result.Images {
        fmt.Printf("Image %d: %s\n", i+1, img)
    }
}
```

## Key Concepts

### Workflow Sources

ComfyKit automatically detects the workflow source:

1. **RunningHub ID**: Pure numeric string (e.g., "12345")
2. **URL**: Starts with `http://` or `https://`
3. **File Path**: Local file path

### Parameters

Use the `$param` syntax in your workflow JSON to define parameters:

```json
{
  "text": "$prompt",
  "seed": "$seed"
}
```

Then pass values at execution time:

```go
params := map[string]interface{}{
    "prompt": "my prompt",
    "seed":   12345,
}
```

## Next Steps

- [Configuration](configuration.md) - Learn about configuration options
- [Local Execution](usage/local.md) - Execute workflows on local ComfyUI
- [Cloud Execution](usage/cloud.md) - Execute workflows on RunningHub

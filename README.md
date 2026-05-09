# ComfyKit-Go

[![](https://img.shields.io/badge/lang-en-blue)](README.md) [![](https://img.shields.io/badge/lang-zh--CN-red)](README.zh-CN.md)
[![Experimental](https://img.shields.io/badge/status-experimental-yellow)](README.md)

**⚠️ Experimental**: This is an experimental Go port of ComfyKit-python. APIs may change without notice.

ComfyKit-Go is a Go SDK for executing ComfyUI workflows. It provides a simple, idiomatic Go API for running ComfyUI workflows both locally and on the RunningHub cloud platform.

This project is the **Go language port** of [ComfyKit-python](https://https://github.com/puke3615/ComfyKit), maintaining the same directory structure, function names, and API design.

## Features

- **Local Execution**: Execute workflows on your local ComfyUI server via HTTP or WebSocket
- **Cloud Execution**: Run workflows on RunningHub cloud platform
- **Auto-detection**: Automatically detect workflow types (local file, URL, RunningHub ID)
- **Parameter Mapping**: Support for ComfyKit DSL parameter syntax
- **Media Handling**: Built-in media upload support for images, audio, and video
- **Cross-platform**: Works on Windows, macOS, and Linux

## Installation

```bash
go get github.com/lazywe/comfykit-go
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/lazywe/comfykit-go"
)

func main() {
    // Create ComfyKit instance with default configuration
    kit := comfykit.NewComfyKit()
    
    // Execute a local workflow
    result, err := kit.Execute("workflow.json", map[string]interface{}{
        "prompt": "a beautiful cat",
        "seed":   42,
    })
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Printf("Status: %s\n", result.Status)
    fmt.Printf("Images: %v\n", result.Images)
}
```

## Configuration

### Local ComfyUI Configuration

```go
kit := comfykit.NewComfyKit(
    comfykit.WithComfyUIBaseURL("http://localhost:8188"),
    comfykit.WithExecutorType(comfykit.ExecutorTypeHTTP),
    comfykit.WithAPIKey("your-api-key"),
    comfykit.WithCookies("session=abc123"),
)
```

### RunningHub Cloud Configuration

```go
kit := comfykit.NewComfyKit(
    comfykit.WithRunningHubAPIKey("your-runninghub-api-key"),
    comfykit.WithRunningHubTimeout(300),
    comfykit.WithRunningHubInstance("plus"),
)
```

## Execution Modes

### Local Workflow File

```go
result, err := kit.Execute("workflow.json", map[string]interface{}{
    "prompt": "a beautiful sunset",
})
```

### RunningHub Workflow ID

```go
result, err := kit.Execute("12345", map[string]interface{}{
    "prompt": "a beautiful sunset",
})
```

### Remote Workflow URL

```go
result, err := kit.Execute("https://example.com/workflow.json", map[string]interface{}{
    "prompt": "a beautiful sunset",
})
```

### Asynchronous Execution

```go
import "time"

// Create task asynchronously
taskID, outputID2Var, err := kit.ExecuteAsyncByID("12345", map[string]interface{}{
    "prompt": "a cute cat",
})
if err != nil {
    fmt.Printf("Failed to create task: %v\n", err)
    return
}

// Poll task status
for {
    result, completed, err := kit.GetTaskCompletion(taskID, outputID2Var)
    if err != nil {
        fmt.Printf("Status check error: %v\n", err)
        time.Sleep(2 * time.Second)
        continue
    }
    
    if completed {
        fmt.Printf("Task completed: %s\n", result.Status)
        if result.Status == "completed" {
            fmt.Printf("Generated %d images\n", len(result.Images))
        } else {
            fmt.Printf("Task failed: %s\n", result.Message)
        }
        break
    }
    
    fmt.Println("Task still running...")
    time.Sleep(2 * time.Second)
}
```

## Workflow DSL

ComfyKit supports a simple DSL for defining parameters in workflow titles:

- `$param` - Basic parameter
- `$param!` - Required parameter
- `$~param` - Parameter requiring media upload
- `$param.field` - Parameter mapping to specific field
- `$output.varname` - Output variable name

Example node title: `"Prompt, $prompt!, $seed"`

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `COMFYUI_BASE_URL` | ComfyUI server URL | `http://127.0.0.1:8188` |
| `COMFYUI_EXECUTOR_TYPE` | Execution type (http/websocket) | `http` |
| `COMFYUI_API_KEY` | ComfyUI API key | - |
| `COMFYUI_COOKIES` | ComfyUI cookies | - |
| `RUNNINGHUB_BASE_URL` | RunningHub API URL | `https://www.runninghub.ai` |
| `RUNNINGHUB_API_KEY` | RunningHub API key | - |
| `RUNNINGHUB_TIMEOUT` | Task timeout in seconds | 0 (unlimited) |
| `RUNNINGHUB_RETRY_COUNT` | API retry count | 3 |
| `RUNNINGHUB_INSTANCE_TYPE` | Instance type (plus) | - |

## API Reference

### ComfyKit

- `NewComfyKit(opts ...ComfyKitOption) *ComfyKit` - Create new ComfyKit instance
- `Execute(workflow string, params map[string]interface{}) (*ExecuteResult, error)` - Execute workflow
- `ExecuteJSON(workflowJSON map[string]interface{}, params map[string]interface{}) (*ExecuteResult, error)` - Execute workflow from JSON
- `ExecuteAsyncByID(workflowID string, params map[string]interface{}) (string, map[string]string, error)` - Create task asynchronously
- `GetTaskCompletion(taskID string, outputID2Var map[string]string) (*ExecuteResult, bool, error)` - Check task status
- `Close() error` - Cleanup resources

### ExecuteResult

- `Status` - Execution status (completed/error/timeout)
- `PromptID` - Prompt/task ID
- `Duration` - Execution duration in seconds
- `Images` - List of image URLs
- `ImagesByVar` - Images grouped by variable name
- `Videos` - List of video URLs
- `Audios` - List of audio URLs
- `Texts` - List of text outputs
- `Message` - Error message if failed

## Acknowledgments

This project is inspired by and based on [ComfyKit-python](https://https://github.com/puke3615/ComfyKit). We would like to thank the original project for its excellent design and implementation.

## License

MIT License

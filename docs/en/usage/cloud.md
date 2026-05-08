# Cloud Execution

Execute workflows on RunningHub cloud platform.

## Prerequisites

1. Sign up for a RunningHub account: https://www.runninghub.ai
2. Get your API key from the RunningHub dashboard

## Configuration

```go
kit := comfykit.NewComfyKit(
    comfykit.WithRunningHubAPIKey("your-api-key"),
)
```

## Execute by RunningHub ID

```go
result, err := kit.Execute("12345", nil)
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

fmt.Printf("Status: %s\n", result.Status)
fmt.Printf("Generated %d images\n", len(result.Images))
```

## Execute with Parameters

```go
params := map[string]interface{}{
    "prompt": "a futuristic city",
    "seed":   12345,
    "steps":  30,
}

result, err := kit.Execute("12345", params)
```

## Execute RunningHub Workflow File

You can also execute local workflow files that contain RunningHub metadata:

```json
{
  "_source": "runninghub",
  "workflow_id": "12345",
  "params": {
    "prompt": "$prompt",
    "seed": "$seed"
  }
}
```

```go
result, err := kit.Execute("workflows/my_runninghub_workflow.json", params)
```

## Cloud-Specific Options

### Timeout

```go
kit := comfykit.NewComfyKit(
    comfykit.WithRunningHubAPIKey("your-api-key"),
    comfykit.WithRunningHubTimeout(600), // 10 minutes
)
```

### Instance Type

```go
kit := comfykit.NewComfyKit(
    comfykit.WithRunningHubAPIKey("your-api-key"),
    comfykit.WithRunningHubInstance("plus"),
)
```

Available instance types:
- `standard` - Default instance
- `plus` - Enhanced performance
- `pro` - Premium performance
- `gpu` - GPU-accelerated

### Retry Count

```go
kit := comfykit.NewComfyKit(
    comfykit.WithRunningHubAPIKey("your-api-key"),
    comfykit.WithRunningHubRetry(5),
)
```

## Mixed Execution

You can use the same ComfyKit instance for both local and cloud execution:

```go
kit := comfykit.NewComfyKit(
    comfykit.WithRunningHubAPIKey("your-api-key"),
)

// Local execution
localResult, _ := kit.Execute("workflows/local.json", params)

// Cloud execution
cloudResult, _ := kit.Execute("12345", params)
```

## Auto-Detection

ComfyKit automatically detects workflow source:

```go
kit := comfykit.NewComfyKit()

// Auto-detects RunningHub ID
result, _ := kit.Execute("12345", params)

// Auto-detects URL
result, _ := kit.Execute("https://example.com/workflow.json", params)

// Auto-detects local file
result, _ := kit.Execute("workflows/local.json", params)
```

## Troubleshooting

### API Key Not Set

```bash
export RUNNINGHUB_API_KEY="your-api-key"
```

Or set it programmatically:

```go
kit := comfykit.NewComfyKit(
    comfykit.WithRunningHubAPIKey("your-api-key"),
)
```

### Task Timeout

Increase the timeout:

```go
kit := comfykit.NewComfyKit(
    comfykit.WithRunningHubAPIKey("your-api-key"),
    comfykit.WithRunningHubTimeout(1200), // 20 minutes
)
```

### Workflow Not Found

Ensure the workflow ID exists in your RunningHub account.

## Cloud Benefits

- No need to maintain local GPU
- Access to powerful GPUs (NVIDIA A100, RTX 4090)
- Automatic scaling
- Pay-as-you-go pricing
- Built-in workflow management

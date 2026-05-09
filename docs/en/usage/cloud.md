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

## Asynchronous Execution

For more flexible task management, you can use asynchronous methods that don't block execution.

### Create Task Asynchronously

```go
import "time"

// Create a task without waiting for completion
taskID, outputID2Var, err := kit.ExecuteAsyncByID("12345", map[string]interface{}{
    "prompt": "a cute cat",
})
if err != nil {
    fmt.Printf("Failed to create task: %v\n", err)
    return
}

fmt.Printf("Task created: %s\n", taskID)
fmt.Printf("Output nodes: %v\n", outputID2Var)
```

### Poll for Task Completion

```go
import "time"

// Poll for task completion
checkInterval := 2 * time.Second
maxChecks := 30 // Max checks before giving up

for i := 0; i < maxChecks; i++ {
    result, completed, err := kit.GetTaskCompletion(taskID, outputID2Var)
    if err != nil {
        fmt.Printf("Status check error: %v, retrying...\n", err)
        time.Sleep(checkInterval)
        continue
    }
    
    if completed {
        if result.Status == "completed" {
            fmt.Printf("Task completed! Generated %d images\n", len(result.Images))
            for _, img := range result.Images {
                fmt.Printf("  - %s\n", img)
            }
        } else {
            fmt.Printf("Task failed: %s - %s\n", result.Status, result.Message)
        }
        break
    }
    
    fmt.Printf("Check %d/%d: Task still running...\n", i+1, maxChecks)
    time.Sleep(checkInterval)
}
```

### Manage Multiple Tasks

```go
import "time"

// Define multiple tasks
tasks := []struct {
    Prompt string
    TaskID string
    Output map[string]string
}{
    {Prompt: "red roses"},
    {Prompt: "blue ocean"},
    {Prompt: "green forest"},
}

// Create all tasks
for i := range tasks {
    taskID, outputID2Var, err := kit.ExecuteAsyncByID("12345", map[string]interface{}{
        "prompt": tasks[i].Prompt,
    })
    if err != nil {
        fmt.Printf("Failed to create task %d: %v\n", i, err)
        continue
    }
    tasks[i].TaskID = taskID
    tasks[i].Output = outputID2Var
    fmt.Printf("Task %d created: %s\n", i, taskID)
}

// Poll for all tasks to complete
remaining := len(tasks)
for remaining > 0 {
    for i := range tasks {
        if tasks[i].TaskID == "" {
            continue // Already completed
        }
        
        result, completed, err := kit.GetTaskCompletion(tasks[i].TaskID, tasks[i].Output)
        if err != nil {
            continue
        }
        
        if completed {
            remaining--
            tasks[i].TaskID = "" // Mark as completed
            if result.Status == "completed" {
                fmt.Printf("Task %d completed: %d images\n", i, len(result.Images))
            } else {
                fmt.Printf("Task %d failed: %s\n", i, result.Message)
            }
        }
    }
    
    time.Sleep(2 * time.Second)
}
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

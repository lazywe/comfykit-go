# Basic Usage

Learn the basics of using ComfyKit-Go.

## Creating ComfyKit Instance

```go
kit := comfykit.NewComfyKit()
```

With custom options:

```go
kit := comfykit.NewComfyKit(
    comfykit.WithComfyUIBaseURL("http://my-server:8188"),
    comfykit.WithAPIKey("my-secret-key"),
)
```

## Executing Workflows

### From File

```go
result, err := kit.Execute("workflows/my_workflow.json", nil)
```

### From URL

```go
result, err := kit.Execute("https://example.com/workflow.json", nil)
```

### From RunningHub ID

```go
result, err := kit.Execute("12345", nil)  // RunningHub workflow ID
```

### With Parameters

```go
params := map[string]interface{}{
    "prompt": "a beautiful landscape",
    "seed":   12345,
    "steps":  30,
}

result, err := kit.Execute("workflows/t2i.json", params)
```

### From JSON

```go
workflowJSON := map[string]interface{}{
    "version": 1,
    "nodes": []interface{}{
        // ... your workflow nodes
    },
}

result, err := kit.ExecuteJSON(workflowJSON, params)
```

## Handling Results

```go
result, err := kit.Execute("workflow.json", params)
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

// Check status
switch result.Status {
case "completed":
    fmt.Println("Execution completed successfully")
case "error":
    fmt.Printf("Error: %s\n", result.Message)
case "timeout":
    fmt.Println("Execution timed out")
}

// Access outputs
fmt.Printf("Generated %d images\n", len(result.Images))
fmt.Printf("Generated %d audios\n", len(result.Audios))
fmt.Printf("Generated %d videos\n", len(result.Videos))
fmt.Printf("Generated %d texts\n", len(result.Texts))

// Access by variable name
for varName, images := range result.ImagesByVar {
    fmt.Printf("%s: %d images\n", varName, len(images))
}

// Get execution duration
fmt.Printf("Duration: %.2f seconds\n", result.Duration)
```

## Cleaning Up

Always close the ComfyKit instance when done:

```go
kit := comfykit.NewComfyKit()
defer kit.Close()

// ... use kit ...
```

## Error Handling

```go
result, err := kit.Execute("workflow.json", params)

if err != nil {
    // Network errors, connection issues, etc.
    log.Fatalf("Failed to execute: %v", err)
}

if result.Status == "error" {
    // Workflow execution errors
    log.Printf("Workflow error: %s", result.Message)
}
```

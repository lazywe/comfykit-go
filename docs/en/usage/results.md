# Working with Results

Learn how to handle execution results.

## Result Structure

```go
type ExecuteResult struct {
    Status       string                 // Execution status
    PromptID     string                 // Task ID
    Duration     float64                // Execution time in seconds
    
    Images       []string               // All image URLs
    ImagesByVar  map[string][]string    // Images grouped by variable
    
    Audios       []string               // All audio URLs
    AudiosByVar  map[string][]string    // Audios grouped by variable
    
    Videos       []string               // All video URLs
    VideosByVar  map[string][]string    // Videos grouped by variable
    
    Texts        []string               // All text outputs
    TextsByVar   map[string][]string    // Texts grouped by variable
    
    Outputs      map[string]interface{} // Raw output data
    Message      string                 // Error message
}
```

## Status Values

| Status | Description |
|--------|-------------|
| `completed` | Execution completed successfully |
| `error` | Execution failed |
| `timeout` | Execution timed out |
| `processing` | Execution in progress |

## Accessing Results

### Basic Access

```go
result, err := kit.Execute("workflow.json", params)
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

// Check status
if result.Status != "completed" {
    fmt.Printf("Execution failed: %s\n", result.Message)
    return
}

// Access all images
for i, img := range result.Images {
    fmt.Printf("Image %d: %s\n", i+1, img)
}
```

### Grouped by Variable

```go
// Access images by variable name
for varName, images := range result.ImagesByVar {
    fmt.Printf("\nVariable: %s\n", varName)
    for i, img := range images {
        fmt.Printf("  %d. %s\n", i+1, img)
    }
}
```

### Different Media Types

```go
// Images
fmt.Printf("Images: %d\n", len(result.Images))

// Audios
fmt.Printf("Audios: %d\n", len(result.Audios))

// Videos
fmt.Printf("Videos: %d\n", len(result.Videos))

// Texts
fmt.Printf("Texts: %d\n", len(result.Texts))
```

### Raw Outputs

```go
// Access raw output data
if result.Outputs != nil {
    for nodeID, output := range result.Outputs {
        fmt.Printf("Node %s: %v\n", nodeID, output)
    }
}
```

## Error Handling

### Connection Errors

```go
result, err := kit.Execute("workflow.json", params)
if err != nil {
    // Network errors, authentication issues, etc.
    log.Printf("Connection error: %v", err)
    return
}
```

### Execution Errors

```go
if result.Status == "error" {
    // Workflow execution failed
    log.Printf("Execution error: %s", result.Message)
    return
}
```

### Timeout Errors

```go
if result.Status == "timeout" {
    log.Printf("Execution timed out after %.2f seconds", result.Duration)
    return
}
```

## Result Examples

### Single Image Output

```go
result, err := kit.Execute("workflows/t2i.json", params)
if err != nil || result.Status != "completed" {
    return
}

if len(result.Images) > 0 {
    fmt.Printf("Generated: %s\n", result.Images[0])
}
```

### Multiple Outputs

```go
result, err := kit.Execute("workflows/multi_output.json", params)
if err != nil || result.Status != "completed" {
    return
}

// Access outputs by variable name
if images, ok := result.ImagesByVar["result"]; ok {
    for _, img := range images {
        fmt.Println(img)
    }
}

if images, ok := result.ImagesByVar["preview"]; ok {
    for _, img := range images {
        fmt.Println(img)
    }
}
```

### Text Outputs

```go
result, err := kit.Execute("workflows/text_gen.json", params)
if err != nil || result.Status != "completed" {
    return
}

for _, text := range result.Texts {
    fmt.Println(text)
}
```

## Advanced Result Processing

### Download Images

```go
import (
    "io"
    "net/http"
    "os"
)

func downloadImage(url, path string) error {
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    out, err := os.Create(path)
    if err != nil {
        return err
    }
    defer out.Close()

    _, err = io.Copy(out, resp.Body)
    return err
}

// Usage
for i, imgURL := range result.Images {
    err := downloadImage(imgURL, fmt.Sprintf("output_%d.png", i))
    if err != nil {
        log.Printf("Failed to download %s: %v", imgURL, err)
    }
}
```

### Batch Processing

```go
prompts := []string{"sunset", "mountain", "city"}
for _, prompt := range prompts {
    result, err := kit.Execute("workflows/t2i.json", map[string]interface{}{
        "prompt": prompt,
    })
    if err != nil || result.Status != "completed" {
        continue
    }
    
    // Process results
    fmt.Printf("Prompt: %s -> %d images\n", prompt, len(result.Images))
}
```

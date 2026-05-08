# Outputs

Learn how to define and access workflow outputs.

## Output Markers

Use `$output.name` to mark output nodes:

```json
{
  "nodes": [
    {
      "id": "1",
      "type": "SaveImage",
      "inputs": {
        "images": ["2", 0],
        "filename_prefix": "$output.result"
      }
    }
  ]
}
```

## Accessing Outputs

```go
result, err := kit.Execute("workflow.json", params)
if err != nil || result.Status != "completed" {
    return
}

// Access all images
for _, img := range result.Images {
    fmt.Println(img)
}

// Access by variable name
if images, ok := result.ImagesByVar["result"]; ok {
    for _, img := range images {
        fmt.Println(img)
    }
}
```

## Multiple Outputs

Define multiple output variables:

```json
{
  "nodes": [
    {
      "id": "1",
      "type": "SaveImage",
      "inputs": {
        "images": ["2", 0],
        "filename_prefix": "$output.main"
      }
    },
    {
      "id": "3",
      "type": "SaveImage",
      "inputs": {
        "images": ["4", 0],
        "filename_prefix": "$output.preview"
      }
    }
  ]
}
```

Access them separately:

```go
result, err := kit.Execute("workflow.json", params)

// Main output
if images, ok := result.ImagesByVar["main"]; ok {
    fmt.Println("Main images:", len(images))
}

// Preview output
if images, ok := result.ImagesByVar["preview"]; ok {
    fmt.Println("Preview images:", len(images))
}
```

## Output Types

ComfyKit supports multiple output types:

### Images

```json
{
  "type": "SaveImage",
  "inputs": {
    "filename_prefix": "$output.images"
  }
}
```

### Audio

```json
{
  "type": "SaveAudio",
  "inputs": {
    "filename_prefix": "$output.audio"
  }
}
```

### Video

```json
{
  "type": "SaveVideo",
  "inputs": {
    "filename_prefix": "$output.video"
  }
}
```

### Text

```json
{
  "type": "SaveText",
  "inputs": {
    "filename_prefix": "$output.text"
  }
}
```

## Output Structure

```go
type ExecuteResult struct {
    // Flat lists
    Images   []string
    Audios   []string
    Videos   []string
    Texts    []string
    
    // Grouped by variable
    ImagesByVar  map[string][]string
    AudiosByVar  map[string][]string
    VideosByVar  map[string][]string
    TextsByVar   map[string][]string
    
    // Raw output data
    Outputs map[string]interface{}
}
```

## Best Practices

1. **Use Descriptive Names**: Use meaningful variable names
2. **Group Related Outputs**: Use the same variable name for related outputs
3. **Handle All Output Types**: Check for images, audios, videos, and texts
4. **Use Raw Outputs**: Access `Outputs` for custom data
5. **Document Outputs**: Add comments to describe what each output contains

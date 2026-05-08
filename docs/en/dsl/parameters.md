# Parameters

Learn how to define and use parameters in your workflows.

## Basic Parameters

Define parameters using the `$` prefix:

```json
{
  "inputs": {
    "text": "$prompt",
    "seed": "$seed",
    "steps": "$steps"
  }
}
```

Pass values at execution time:

```go
params := map[string]interface{}{
    "prompt": "a beautiful sunset",
    "seed":   12345,
    "steps":  25,
}

result, err := kit.Execute("workflow.json", params)
```

## Required Parameters

Mark parameters as required with `!`:

```json
{
  "inputs": {
    "text": "$prompt!"
  }
}
```

If a required parameter is not provided, ComfyKit will panic:

```go
// This will panic because "prompt" is required
result, err := kit.Execute("workflow.json", nil)
```

## Media Upload Parameters

Mark parameters that require media upload with `~`:

```json
{
  "inputs": {
    "image": "$~input_image"
  }
}
```

ComfyKit will automatically upload media files:

```go
params := map[string]interface{}{
    "input_image": "/path/to/image.png",      // Local file
    // or
    "input_image": "https://example.com/img.jpg",  // URL
}
```

## Required Media Upload

Combine both markers:

```json
{
  "inputs": {
    "image": "$~input_image!"
  }
}
```

## Field Access

Access nested fields using dot notation:

```json
{
  "inputs": {
    "text": "$prompt.text"
  }
}
```

With nested parameters:

```go
params := map[string]interface{}{
    "prompt": map[string]interface{}{
        "text": "a beautiful sunset",
        "style": "cinematic",
    },
}
```

## Supported Parameter Types

| Type | Description | Example |
|------|-------------|---------|
| `str` | String | `"hello"` |
| `int` | Integer | `42` |
| `float` | Float | `3.14` |
| `bool` | Boolean | `true` |
| `image` | Image path/URL | `/path/to/img.png` |
| `audio` | Audio path/URL | `/path/to/audio.mp3` |
| `video` | Video path/URL | `/path/to/video.mp4` |

## Default Values

You can define default values in the workflow metadata:

```json
{
  "__metadata__": {
    "params": {
      "prompt": {
        "type": "str",
        "required": true,
        "default": "a beautiful scene"
      },
      "seed": {
        "type": "int",
        "required": false,
        "default": 0
      }
    }
  }
}
```

## Parameter Validation

ComfyKit validates parameters:

1. **Required Check**: Ensures required parameters are provided
2. **Type Check**: Attempts to convert values to expected types
3. **Media Check**: Validates media files exist or URLs are valid

## Best Practices

1. **Use Descriptive Names**: Use clear parameter names like `prompt`, `seed`, `steps`
2. **Mark Required Params**: Use `!` for required parameters
3. **Use Media Upload**: Use `~` for image/audio/video inputs
4. **Provide Defaults**: Define default values in metadata when possible
5. **Document Parameters**: Add comments or metadata to describe parameters

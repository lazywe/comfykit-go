# DSL Overview

ComfyKit uses a simple DSL (Domain Specific Language) to define parameters and outputs in workflows.

## Syntax

The DSL uses a simple `$` prefix to mark parameters:

```json
{
  "inputs": {
    "text": "$prompt",
    "seed": "$seed",
    "image": "$~input_image"
  }
}
```

## Parameter Types

### Basic Parameter

```json
"text": "$prompt"
```

A basic parameter that accepts any value.

### Required Parameter

```json
"text": "$prompt!"
```

Marks the parameter as required.

### Media Upload Parameter

```json
"image": "$~input_image"
```

Marks the parameter as requiring media upload.

### Required Media Upload

```json
"image": "$~input_image!"
```

Required parameter that requires media upload.

### Field Access

```json
"text": "$prompt.text"
```

Access a nested field in the parameter.

## Output Markers

```json
{
  "inputs": {
    "filename_prefix": "$output.result"
  }
}
```

Marks an output node with a variable name.

## Complete Example

```json
{
  "nodes": [
    {
      "id": "1",
      "type": "CLIPTextEncode",
      "inputs": {
        "text": "$prompt!"  // Required parameter
      }
    },
    {
      "id": "2",
      "type": "LoadImage",
      "inputs": {
        "image": "$~input_image"  // Media upload parameter
      }
    },
    {
      "id": "3",
      "type": "SaveImage",
      "inputs": {
        "images": ["4", 0],
        "filename_prefix": "$output.result"  // Output marker
      }
    }
  ]
}
```

## Benefits

- **Simplicity**: Easy to understand syntax
- **Flexibility**: Supports various parameter types
- **Automatic**: ComfyKit handles parameter injection automatically
- **Validation**: Built-in required parameter validation
- **Media Handling**: Automatic media upload support

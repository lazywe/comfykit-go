# Examples

ComfyKit-Go includes several examples to help you get started.

## Example List

### 01_quick_start.go
- Basic workflow execution
- Simple 3-line example
- Result handling

### 02_configuration.go
- Default configuration
- Custom parameters
- Environment variables
- Configuration priority
- Multiple instances

### 03_local_workflows.go
- Local ComfyUI execution
- Custom parameters
- Execute from JSON
- Different output types
- WebSocket executor

### 04_runninghub_cloud.go
- RunningHub cloud execution
- Cloud with parameters
- Cloud workflow files
- Auto-detection
- Mixed local + cloud

### 05_advanced_features.go
- Batch execution
- Error handling
- Timeout handling
- Authentication
- Executor types

## Running Examples

```bash
# Run a single example
go run examples/01_quick_start.go

# Run all examples
go run examples/run_all.go
```

## Example Structure

Each example follows this pattern:

```go
/*
Example Title - Brief Description

Detailed explanation of what this example teaches.
*/

package main

import (
    "fmt"
    "github.com/lazywe/comfykit-go"
)

func main() {
    // Create ComfyKit instance
    kit := comfykit.NewComfyKit()
    
    // Execute workflow
    result, err := kit.Execute("workflow.json", params)
    
    // Handle result
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Printf("Status: %s\n", result.Status)
}
```

## Creating Your Own Examples

Follow these steps to create your own examples:

1. **Create a new file** in the `examples/` directory
2. **Add a descriptive header** explaining what the example does
3. **Import ComfyKit** and any other required packages
4. **Create a ComfyKit instance** with appropriate configuration
5. **Execute a workflow** with parameters
6. **Handle the result** and print meaningful output
7. **Add comments** explaining each step

## Example Tips

1. **Start Simple**: Begin with basic execution
2. **Test Incrementally**: Test each part before adding more complexity
3. **Handle Errors**: Always handle errors gracefully
4. **Use Logging**: Add logging for debugging
5. **Document**: Explain what your example does

## Contributing Examples

If you have a useful example, consider contributing it:

1. Follow the existing example pattern
2. Add appropriate comments
3. Include error handling
4. Add documentation in the comments
5. Submit a pull request

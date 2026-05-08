# Development

Learn how to develop and contribute to ComfyKit-Go.

## Setup Development Environment

```bash
# Clone the repository
git clone https://github.com/lazywe/comfykit-go.git
cd comfykit-go

# Install dependencies
go mod download

# Build the project
go build ./...

# Run tests
go test ./...
```

## Project Structure

```
comfykit-go/
├── comfykit.go          # Main ComfyKit class
├── executor.go          # Executor interface
├── logger.go            # Logging utilities
├── comfyui/             # ComfyUI related code
�?  ├── executor.go      # Base executor
�?  ├── http_executor.go # HTTP executor
�?  ├── websocket_executor.go # WebSocket executor
�?  ├── runninghub_client.go  # RunningHub API client
�?  ├── runninghub_executor.go # RunningHub executor
�?  ├── models.go        # Data models
�?  └── workflow_parser.go # Workflow parser
├── utils/               # Utility functions
�?  ├── config_util.go   # Configuration utilities
�?  ├── file_util.go     # File utilities
�?  ├── network_util.go  # Network utilities
�?  └── ...
├── examples/            # Example scripts
├── workflows/           # Example workflows
└── docs/                # Documentation
```

## Coding Guidelines

### Style
- Follow Go conventions (gofmt)
- Use `camelCase` for function names
- Use `PascalCase` for types and exported identifiers
- Use `snake_case` for internal functions

### Error Handling
- Return errors explicitly
- Use descriptive error messages
- Wrap errors with context
- Avoid silent failures

### Testing
- Write unit tests for new functionality
- Follow existing test patterns
- Run tests before submitting changes
- Test edge cases

### Documentation
- Document exported functions and types
- Add comments for complex logic
- Update documentation when adding features
- Use clear and concise language

## Adding New Features

### Step 1: Plan
- Define the feature requirements
- Design the API
- Consider backward compatibility

### Step 2: Implement
- Write the code
- Add tests
- Update documentation

### Step 3: Test
- Run existing tests
- Verify new functionality
- Test edge cases

### Step 4: Document
- Update API reference
- Add examples if needed
- Update documentation

## Debugging

### Enable Debug Logging
```go
kit := comfykit.NewComfyKit()
// Use verbose logging
```

### Print Debug Information
```go
result, err := kit.Execute("workflow.json", params)
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

// Print detailed result
fmt.Printf("Result: %+v\n", result)
```

### Use Logging
```go
import "log"

log.Printf("Executing workflow: %s", workflowPath)
log.Printf("Parameters: %v", params)
log.Printf("Result status: %s", result.Status)
```

## Testing

### Run All Tests
```bash
go test ./...
```

### Run Specific Tests
```bash
go test ./comfyui/...
go test ./utils/...
```

### Test Coverage
```bash
go test -cover ./...
```

## Releasing

### Update Version
- Update version in `go.mod`
- Update version in documentation

### Create Release Notes
- Document changes
- List breaking changes
- Provide migration guide

### Publish
- Push to GitHub
- Create release tag
- Update documentation

## Contributing

### Fork the Repository
1. Fork the repository on GitHub
2. Clone your fork
3. Create a feature branch

### Make Changes
1. Implement your changes
2. Add tests
3. Update documentation

### Submit PR
1. Push your changes
2. Create a pull request
3. Wait for review
4. Address feedback

## Code Review Guidelines

### What to Look For
- Code correctness
- Performance
- Security
- Readability
- Test coverage
- Documentation

### Feedback
- Be constructive
- Explain your reasoning
- Suggest improvements
- Encourage best practices

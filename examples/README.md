# ComfyKit-Go Examples

Welcome to ComfyKit-Go examples! These examples serve both as tutorials and functional tests.

## 📚 Learning Path

Follow these examples in order to learn ComfyKit-Go:

1. **[01_quick_start.go](01_quick_start.go)** - Your first ComfyKit-Go program
   - Simple 3-line example
   - Basic workflow execution
   - Result handling

2. **[02_configuration.go](02_configuration.go)** - Configuration options
   - Default configuration
   - Custom parameters
   - Environment variables
   - Configuration priority
   - Multiple executors

3. **[03_local_workflows.go](03_local_workflows.go)** - Local ComfyUI execution
   - Basic workflow execution
   - Custom parameters
   - Execute from dict
   - Different output types
   - Result handling

4. **[04_runninghub_cloud.go](04_runninghub_cloud.go)** - RunningHub cloud execution
   - Cloud execution basics
   - Custom configuration
   - Unified parameters
   - Auto-detection
   - Mixed local + cloud

5. **[05_advanced_features.go](05_advanced_features.go)** - Advanced features
   - Batch execution
   - Error handling
   - Timeout handling
   - Authentication
   - Executor types
   - Result processing

## 🚀 Quick Start

Run a single example:

```bash
go run examples/01_quick_start.go
```

Run all examples (integration test):

```bash
go run examples/run_all.go
```

## 📋 Prerequisites

### For Local Execution (examples 01, 03, 05)

- ComfyUI running at `http://127.0.0.1:8188`
- Or set `COMFYUI_BASE_URL` to your ComfyUI server

### For Cloud Execution (example 04)

- RunningHub API key
- Set environment variable:

  **Linux/macOS (bash/zsh):**
  ```bash
  export RUNNINGHUB_API_KEY='your-api-key-here'
  ```

  **Windows (Command Prompt):**
  ```cmd
  set RUNNINGHUB_API_KEY=your-api-key-here
  ```

  **Windows (PowerShell):**
  ```powershell
  $env:RUNNINGHUB_API_KEY='your-api-key-here'
  ```

## 🧪 Running as Tests

All examples include assertions and can serve as integration tests.

## 💡 Tips

- **Start Simple**: Begin with `01_quick_start.go`
- **Read Code**: Each example is heavily commented
- **Run & Experiment**: Modify parameters and see what happens
- **Check Errors**: Error messages guide you to solutions

## 📚 Next Steps

After completing these examples, check out:

- [README.md](../README.md) - Project documentation
- [workflows/](../workflows/) - Example workflows
- [comfykit/](../comfykit/) - Source code

Happy coding with ComfyKit-Go! 🚀

//go:build ignore

/*
Configuration Example - Different Ways to Configure ComfyKit-Go

Learn how to configure ComfyKit-Go for different scenarios.
Priority: constructor parameters > environment variables > defaults
*/

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/lazywe/comfykit-go"
)

func exampleDefault() {
	fmt.Println("\n=== Example 1: Default Configuration ===")

	kit := comfykit.NewComfyKit()
	fmt.Printf("âœ?ComfyUI URL: %s\n", kit.GetComfyUIBaseURL())
	fmt.Printf("âœ?Executor type: %s\n", kit.GetExecutorType())
}

func exampleCustomParams() {
	fmt.Println("\n=== Example 2: Custom Parameters ===")

	kit := comfykit.NewComfyKit(
		comfykit.WithComfyUIBaseURL("http://my-server:8188"),
		comfykit.WithExecutorType(comfykit.ExecutorTypeWebSocket),
		comfykit.WithAPIKey("my-api-key"),
	)

	fmt.Printf("âœ?ComfyUI URL: %s\n", kit.GetComfyUIBaseURL())
	fmt.Printf("âœ?Executor type: %s\n", kit.GetExecutorType())
	apiKey := kit.GetAPIKey()
	if apiKey != "" {
		fmt.Println("âœ?API Key: ***")
	} else {
		fmt.Println("âœ?API Key: Not set")
	}
}

func exampleEnvVars() {
	fmt.Println("\n=== Example 3: Environment Variables ===")

	// Set env vars (normally done in shell or .env file)
	os.Setenv("COMFYUI_BASE_URL", "http://env-server:8188")
	os.Setenv("COMFYUI_API_KEY", "env-key")

	kit := comfykit.NewComfyKit()
	fmt.Printf("âœ?ComfyUI URL: %s\n", kit.GetComfyUIBaseURL())
	apiKey := kit.GetAPIKey()
	if apiKey != "" {
		fmt.Println("âœ?API Key: ***")
	} else {
		fmt.Println("âœ?API Key: Not set")
	}

	// Clean up
	os.Unsetenv("COMFYUI_BASE_URL")
	os.Unsetenv("COMFYUI_API_KEY")
}

func examplePriority() {
	fmt.Println("\n=== Example 4: Configuration Priority ===")

	// Set env var
	os.Setenv("COMFYUI_BASE_URL", "http://env:8188")

	// Param overrides env var
	kit := comfykit.NewComfyKit(comfykit.WithComfyUIBaseURL("http://param:8188"))
	fmt.Printf("âœ?Param overrides env: %s\n", kit.GetComfyUIBaseURL())

	// Clean up
	os.Unsetenv("COMFYUI_BASE_URL")
}

func exampleRunningHub() {
	fmt.Println("\n=== Example 5: RunningHub Cloud ===")

	kit := comfykit.NewComfyKit(
		comfykit.WithRunningHubAPIKey("rh-key-xxx"),
		comfykit.WithRunningHubTimeout(600),
	)

	fmt.Printf("âœ?RunningHub URL: %s\n", kit.GetRunningHubBaseURL())
	fmt.Printf("âœ?RunningHub timeout: %ds\n", kit.GetRunningHubTimeout())
}

func exampleMultipleInstances() {
	fmt.Println("\n=== Example 6: Multiple Instances ===")

	kitLocal := comfykit.NewComfyKit(comfykit.WithComfyUIBaseURL("http://127.0.0.1:8188"))
	kitRemote := comfykit.NewComfyKit(comfykit.WithComfyUIBaseURL("http://remote:8188"))

	fmt.Printf("âœ?Local: %s\n", kitLocal.GetComfyUIBaseURL())
	fmt.Printf("âœ?Remote: %s\n", kitRemote.GetComfyUIBaseURL())
}

func main() {
	fmt.Println("ðŸ”§ ComfyKit-Go Configuration Examples")
	fmt.Println("=" + strings.Repeat("=", 59))

	exampleDefault()
	exampleCustomParams()
	exampleEnvVars()
	examplePriority()
	exampleRunningHub()
	exampleMultipleInstances()

	fmt.Println("\n" + "=" + strings.Repeat("=", 59))
	fmt.Println("âœ?All configuration examples completed!")
	fmt.Println("\nðŸ“ Key Takeaways:")
	fmt.Println("  1. Simple and direct - no Config struct needed")
	fmt.Println("  2. Priority: params > env vars > defaults")
	fmt.Println("  3. Multiple instances can have independent configs")
	fmt.Println("\nNext: Check out 03_local_workflows.go for workflow execution")
}

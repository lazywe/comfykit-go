//go:build ignore

/*
Advanced Features Example - Advanced ComfyKit-Go Usage

Learn advanced features like batch execution, error handling, and authentication.
*/

package main

import (
	"fmt"
	"strings"

	"github.com/lazywe/comfykit-go"
)

func exampleBatchExecution() {
	fmt.Println("\n=== Example 1: Batch Execution ===")

	kit := comfykit.NewComfyKit()

	prompts := []string{
		"a beautiful sunset",
		"a mountain landscape",
		"a city skyline",
	}

	for i, prompt := range prompts {
		fmt.Printf("\nProcessing prompt %d/%d...\n", i+1, len(prompts))
		result, err := kit.Execute("workflows/t2i_by_local_flux.json", map[string]interface{}{
			"prompt": prompt,
		})
		if err != nil {
			fmt.Printf("‚ö†Ô∏è  Error: %v\n", err)
			continue
		}
		fmt.Printf("‚ú?Generated %d images\n", len(result.Images))
	}
}

func exampleErrorHandling() {
	fmt.Println("\n=== Example 2: Error Handling ===")

	kit := comfykit.NewComfyKit()

	result, err := kit.Execute("non_existent_workflow.json", nil)
	if err != nil {
		fmt.Printf("‚ú?Caught error: %v\n", err)
	} else if result.Status == "error" {
		fmt.Printf("‚ú?Workflow error: %s\n", result.Message)
	}
}

func exampleTimeout() {
	fmt.Println("\n=== Example 3: Timeout Handling ===")

	// Set a short timeout for demonstration
	kit := comfykit.NewComfyKit(
		comfykit.WithRunningHubTimeout(30),
	)
	fmt.Printf("‚ú?Timeout configured: %ds\n", kit.GetRunningHubTimeout())
}

func exampleAuthentication() {
	fmt.Println("\n=== Example 4: Authentication ===")

	// With API key
	_ = comfykit.NewComfyKit(comfykit.WithAPIKey("my-secret-key"))
	fmt.Println("‚ú?Configured with API key")

	// With cookies
	_ = comfykit.NewComfyKit(comfykit.WithCookies("session=abc123"))
	fmt.Println("‚ú?Configured with cookies")
}

func exampleCloseResources() {
	fmt.Println("\n=== Example 5: Resource Management ===")

	kit := comfykit.NewComfyKit()
	defer kit.Close() // Ensure resources are cleaned up

	fmt.Println("‚ú?Kit created with deferred cleanup")
}

func exampleExecutorTypes() {
	fmt.Println("\n=== Example 6: Executor Types ===")

	// HTTP executor (default)
	kitHTTP := comfykit.NewComfyKit(comfykit.WithExecutorType(comfykit.ExecutorTypeHTTP))
	fmt.Printf("‚ú?HTTP executor: %s\n", kitHTTP.GetExecutorType())

	// WebSocket executor
	kitWS := comfykit.NewComfyKit(comfykit.WithExecutorType(comfykit.ExecutorTypeWebSocket))
	fmt.Printf("‚ú?WebSocket executor: %s\n", kitWS.GetExecutorType())
}

func main() {
	fmt.Println("‚ö?ComfyKit-Go Advanced Features")
	fmt.Println("=" + strings.Repeat("=", 59))

	exampleBatchExecution()
	exampleErrorHandling()
	exampleTimeout()
	exampleAuthentication()
	exampleCloseResources()
	exampleExecutorTypes()

	fmt.Println("\n" + "=" + strings.Repeat("=", 59))
	fmt.Println("‚ú?Advanced features examples completed!")
	fmt.Println("\nüìö All examples completed! Check the documentation for more.")
}

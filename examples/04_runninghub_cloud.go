//go:build ignore

/*
RunningHub Cloud Example - Execute Workflows on RunningHub Cloud

Learn how to execute workflows on the RunningHub cloud platform.
*/

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/lazywe/comfykit-go"
)

func exampleBasicCloudExecution() {
	fmt.Println("\n=== Example 1: Basic Cloud Execution ===")

	apiKey := os.Getenv("RUNNINGHUB_API_KEY")
	if apiKey == "" {
		fmt.Println("⚠️  Skipping - RUNNINGHUB_API_KEY not set")
		fmt.Println("   To test, set the environment variable first")
		fmt.Println("   Example (PowerShell): $env:RUNNINGHUB_API_KEY='your-key'")
		return
	}

	kit := comfykit.NewComfyKit(comfykit.WithRunningHubAPIKey(apiKey))

	// Execute by RunningHub workflow ID
	result, err := kit.Execute("12345", nil)
	if err != nil {
		fmt.Printf("⚠️  Execution error: %v\n", err)
		return
	}

	fmt.Printf("�?Status: %s\n", result.Status)
	if result.Status == "error" && result.Message != "" {
		fmt.Printf("⚠️  Error message: %s\n", result.Message)
	} else {
		fmt.Printf("�?Images: %d\n", len(result.Images))
	}
}

func exampleCloudWithParams() {
	fmt.Println("\n=== Example 2: Cloud with Parameters ===")

	apiKey := os.Getenv("RUNNINGHUB_API_KEY")
	if apiKey == "" {
		fmt.Println("⚠️  Skipping - RUNNINGHUB_API_KEY not set")
		fmt.Println("   To test, set the environment variable first")
		return
	}

	kit := comfykit.NewComfyKit(comfykit.WithRunningHubAPIKey(apiKey))

	params := map[string]interface{}{
		"prompt": "a futuristic city skyline",
	}

	result, err := kit.Execute("1988434426705133569", params)
	if err != nil {
		fmt.Printf("⚠️  Execution error: %v\n", err)
		return
	}

	fmt.Printf("�?Status: %s\n", result.Status)
	if result.Status == "error" && result.Message != "" {
		fmt.Printf("⚠️  Error message: %s\n", result.Message)
	}
}

func exampleCloudWorkflowFile() {
	fmt.Println("\n=== Example 3: Cloud Workflow File ===")

	apiKey := os.Getenv("RUNNINGHUB_API_KEY")
	if apiKey == "" {
		fmt.Println("⚠️  Skipping - RUNNINGHUB_API_KEY not set")
		fmt.Println("   To test, set the environment variable first")
		return
	}

	kit := comfykit.NewComfyKit(comfykit.WithRunningHubAPIKey(apiKey))

	// Execute a local workflow file that contains RunningHub metadata
	result, err := kit.Execute("workflows/my_runninghub_workflow.json", nil)
	if err != nil {
		fmt.Printf("⚠️  Execution error: %v\n", err)
		return
	}

	fmt.Printf("�?Status: %s\n", result.Status)
	if result.Status == "error" && result.Message != "" {
		fmt.Printf("⚠️  Error message: %s\n", result.Message)
	}
}

func exampleAutoDetection() {
	fmt.Println("\n=== Example 4: Auto Detection ===")

	// Auto-detects RunningHub ID (numeric string)
	fmt.Println("�?RunningHub ID format: pure numeric")
	fmt.Println("�?URL format: starts with http:// or https://")
	fmt.Println("�?File path: contains / or \\")
	fmt.Println("�?RunningHub workflow file: contains _source: runninghub")
}

func exampleMixedExecution() {
	fmt.Println("\n=== Example 5: Mixed Local + Cloud ===")

	apiKey := os.Getenv("RUNNINGHUB_API_KEY")
	if apiKey == "" {
		fmt.Println("⚠️  Skipping cloud part - RUNNINGHUB_API_KEY not set")
	}

	// Local execution
	fmt.Println("�?Local execution: kit.Execute('workflows/local.json', params)")

	// Cloud execution
	fmt.Println("�?Cloud execution: kit.Execute('12345', params)")
}

func main() {
	fmt.Println("☁️  ComfyKit-Go RunningHub Cloud Examples")
	fmt.Println("=" + strings.Repeat("=", 59))

	exampleBasicCloudExecution()
	exampleCloudWithParams()
	exampleCloudWorkflowFile()
	exampleAutoDetection()
	exampleMixedExecution()

	fmt.Println("\n" + "=" + strings.Repeat("=", 59))
	fmt.Println("�?RunningHub cloud examples completed!")
	fmt.Println("\nNext: Check out 05_advanced_features.go for advanced usage")
}

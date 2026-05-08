//go:build ignore

/*
Local Workflows Example - Execute Workflows on Local ComfyUI

Learn how to execute workflows on your local ComfyUI server.
*/

package main

import (
	"fmt"
	"strings"

	"github.com/lazywe/comfykit-go"
)

func exampleBasicExecution() {
	fmt.Println("\n=== Example 1: Basic Workflow Execution ===")

	kit := comfykit.NewComfyKit()
	result, err := kit.Execute("workflows/t2i_by_local_flux.json", nil)
	if err != nil {
		fmt.Printf("‚ö†Ô∏è  Execution error (expected if ComfyUI not running): %v\n", err)
		return
	}

	fmt.Printf("‚ú?Status: %s\n", result.Status)
	fmt.Printf("‚ú?Images generated: %d\n", len(result.Images))
}

func exampleWithParameters() {
	fmt.Println("\n=== Example 2: With Custom Parameters ===")

	kit := comfykit.NewComfyKit()
	params := map[string]interface{}{
		"prompt": "a beautiful landscape with mountains",
		"seed":   12345,
	}

	result, err := kit.Execute("workflows/t2i_by_local_flux.json", params)
	if err != nil {
		fmt.Printf("‚ö†Ô∏è  Execution error: %v\n", err)
		return
	}

	fmt.Printf("‚ú?Status: %s\n", result.Status)
	fmt.Printf("‚ú?Used prompt: %s\n", params["prompt"])
}

func exampleExecuteJSON() {
	fmt.Println("\n=== Example 3: Execute from JSON ===")

	kit := comfykit.NewComfyKit()

	workflowJSON := map[string]interface{}{
		"version": 1,
		"nodes":   []interface{}{},
	}

	result, err := kit.ExecuteJSON(workflowJSON, nil)
	if err != nil {
		fmt.Printf("‚ö†Ô∏è  Execution error: %v\n", err)
		return
	}

	fmt.Printf("‚ú?Status: %s\n", result.Status)
}

func exampleResultTypes() {
	fmt.Println("\n=== Example 4: Different Output Types ===")

	kit := comfykit.NewComfyKit()
	result, err := kit.Execute("workflows/t2i_by_local_flux.json", nil)
	if err != nil {
		fmt.Printf("‚ö†Ô∏è  Execution error: %v\n", err)
		return
	}

	fmt.Printf("‚ú?Images: %d\n", len(result.Images))
	fmt.Printf("‚ú?Audios: %d\n", len(result.Audios))
	fmt.Printf("‚ú?Videos: %d\n", len(result.Videos))
	fmt.Printf("‚ú?Texts: %d\n", len(result.Texts))

	if len(result.ImagesByVar) > 0 {
		fmt.Println("\nüìä Images by variable:")
		for varName, images := range result.ImagesByVar {
			fmt.Printf("  - %s: %d images\n", varName, len(images))
		}
	}
}

func exampleWebSocketExecutor() {
	fmt.Println("\n=== Example 5: WebSocket Executor ===")

	kit := comfykit.NewComfyKit(comfykit.WithExecutorType(comfykit.ExecutorTypeWebSocket))
	fmt.Printf("‚ú?Using WebSocket executor\n")

	result, err := kit.Execute("workflows/t2i_by_local_flux.json", nil)
	if err != nil {
		fmt.Printf("‚ö†Ô∏è  Execution error: %v\n", err)
		return
	}

	fmt.Printf("‚ú?Status: %s\n", result.Status)
}

func main() {
	fmt.Println("üñ•Ô∏? ComfyKit-Go Local Workflow Examples")
	fmt.Println("=" + strings.Repeat("=", 59))

	exampleBasicExecution()
	exampleWithParameters()
	exampleExecuteJSON()
	exampleResultTypes()
	exampleWebSocketExecutor()

	fmt.Println("\n" + "=" + strings.Repeat("=", 59))
	fmt.Println("‚ú?Local workflow examples completed!")
	fmt.Println("\nNext: Check out 04_runninghub_cloud.go for cloud execution")
}

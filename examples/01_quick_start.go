//go:build ignore

/*
Quick Start Example - Your First ComfyKit-Go Program

This is the simplest way to use ComfyKit-Go. Just a few lines of code!
*/

package main

import (
	"fmt"
	"strings"

	"github.com/lazywe/comfykit-go"
)

func main() {
	fmt.Println("ðŸš€ ComfyKit-Go Quick Start")
	fmt.Println()
	fmt.Println("=" + strings.Repeat("=", 59))

	// Step 1: Create ComfyKit instance (connects to local ComfyUI by default)
	kit := comfykit.NewComfyKit()
	fmt.Printf("âœ?Connected to ComfyUI at %s\n", kit.GetComfyUIBaseURL())

	// Step 2: Execute a workflow
	fmt.Println("\nðŸ“ Executing workflow...")
	result, err := kit.Execute(
		"workflows/t2i_by_local_flux.json",
		map[string]interface{}{"prompt": "a beautiful sunset over the ocean"},
	)
	if err != nil {
		fmt.Printf("â?Error: %v\n", err)
		return
	}

	// Step 3: Check results
	fmt.Printf("\nâœ?Status: %s\n", result.Status)
	if result.Duration > 0 {
		fmt.Printf("âœ?Duration: %.2fs\n", result.Duration)
	}
	fmt.Printf("âœ?Generated %d images\n", len(result.Images))

	if len(result.Images) > 0 {
		fmt.Println("\nðŸ–¼ï¸? Images:")
		for i, img := range result.Images {
			fmt.Printf("   %d. %s\n", i+1, img)
		}
	}

	fmt.Println("\n" + "=" + strings.Repeat("=", 59))
	fmt.Println("âœ?Quick start completed successfully!")
	fmt.Println("\nNext: Check out 02_configuration.go to learn about configuration options")
}

//go:build ignore

/*
Run All Examples - Integration Test Runner

Run all ComfyKit-Go examples to verify functionality.
*/

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	fmt.Println("🧪 ComfyKit-Go Example Runner")
	fmt.Println("=" + strings.Repeat("=", 59))
	fmt.Println()

	examples := []string{
		"01_quick_start.go",
		"02_configuration.go",
		"03_local_workflows.go",
		"04_runninghub_cloud.go",
		"05_advanced_features.go",
	}

	successCount := 0
	failureCount := 0
	skipCount := 0

	// 获取当前源文件所在目录的绝对路径
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Println("�?Failed to get source file path")
		os.Exit(1)
	}
	examplesDir := filepath.Dir(filename)

	for _, example := range examples {
		fmt.Printf("▶️  Running %s...\n", example)

		cmd := exec.Command("go", "run", example)
		cmd.Dir = examplesDir
		output, err := cmd.CombinedOutput()

		if err != nil {
			outputStr := string(output)
			if strings.Contains(outputStr, "RUNNINGHUB_API_KEY not set") {
				fmt.Println("  ⏭️  Skipped (needs RunningHub API key)")
				skipCount++
			} else if strings.Contains(outputStr, "ComfyUI not running") ||
				strings.Contains(outputStr, "connection refused") {
				fmt.Println("  ⏭️  Skipped (ComfyUI not running)")
				skipCount++
			} else {
				fmt.Printf("  �?Failed: %v\n", err)
				fmt.Printf("     Output: %s\n", outputStr)
				failureCount++
			}
		} else {
			fmt.Println("  �?Success")
			successCount++
		}
		fmt.Println()
	}

	fmt.Println("=" + strings.Repeat("=", 59))
	fmt.Printf("📊 Results: %d passed, %d failed, %d skipped\n", successCount, failureCount, skipCount)

	if failureCount > 0 {
		os.Exit(1)
	}
}

//go:build ignore

/*
Asynchronous Execution Example - Create tasks and poll for results

Learn how to asynchronously create tasks and poll for completion using ExecuteAsyncByID and GetTaskCompletion.
*/

package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/lazywe/comfykit-go"
	"github.com/lazywe/comfykit-go/comfyui"
)

func exampleAsyncTaskCreation() {
	fmt.Println("\n=== Example 1: Async Task Creation ===")

	apiKey := os.Getenv("RUNNINGHUB_API_KEY")
	if apiKey == "" {
		fmt.Println("⚠️  Skipping - RUNNINGHUB_API_KEY not set")
		fmt.Println("   To test, set the environment variable first")
		fmt.Println("   Example (PowerShell): $env:RUNNINGHUB_API_KEY='your-key'")
		return
	}

	kit := comfykit.NewComfyKit(comfykit.WithRunningHubAPIKey(apiKey))

	// Create a task asynchronously
	taskID, outputID2Var, err := kit.ExecuteAsyncByID("1988434426705133569", map[string]interface{}{
		"prompt": "a cute cat sitting on a sofa",
	})
	if err != nil {
		fmt.Printf("⚠️  Failed to create task: %v\n", err)
		return
	}

	fmt.Printf("✅ Task created successfully!\n")
	fmt.Printf("   Task ID: %s\n", taskID)
	fmt.Printf("   Output nodes: %d\n", len(outputID2Var))
	for nodeID, varName := range outputID2Var {
		fmt.Printf("     %s: %s\n", nodeID, varName)
	}
}

func examplePollForResults() {
	fmt.Println("\n=== Example 2: Poll for Results ===")

	apiKey := os.Getenv("RUNNINGHUB_API_KEY")
	if apiKey == "" {
		fmt.Println("⚠️  Skipping - RUNNINGHUB_API_KEY not set")
		return
	}

	kit := comfykit.NewComfyKit(comfykit.WithRunningHubAPIKey(apiKey))

	// First, create a task
	fmt.Println("🔧 Creating task...")
	taskID, outputID2Var, err := kit.ExecuteAsyncByID("1988434426705133569", map[string]interface{}{
		"prompt": "a futuristic city at night",
	})
	if err != nil {
		fmt.Printf("⚠️  Failed to create task: %v\n", err)
		return
	}
	fmt.Printf("✅ Task ID: %s\n", taskID)

	// Poll for completion
	fmt.Println("⏳ Polling for completion...")
	checkInterval := 2 * time.Second
	maxChecks := 30 // 30 checks * 2 seconds = 1 minute timeout
	checkCount := 0

	for {
		checkCount++
		if checkCount > maxChecks {
			fmt.Println("⚠️  Timeout - giving up")
			return
		}

		result, completed, err := kit.GetTaskCompletion(taskID, outputID2Var)
		if err != nil {
			fmt.Printf("⚠️  Error checking status: %v, retrying...\n", err)
			time.Sleep(checkInterval)
			continue
		}

		if completed {
			fmt.Printf("✅ Task completed! (took %d checks)\n", checkCount)
			fmt.Printf("   Status: %s\n", result.Status)

			if result.Status == "error" {
				fmt.Printf("   Error message: %s\n", result.Message)
				return
			}

			fmt.Printf("   Images found: %d\n", len(result.Images))
			for i, imageURL := range result.Images {
				fmt.Printf("   Image %d: %s\n", i+1, imageURL)
			}
			return
		}

		// Not completed yet
		fmt.Printf("   Check %02d/%d: task still running...\n", checkCount, maxChecks)
		time.Sleep(checkInterval)
	}
}

func exampleMultipleAsyncTasks() {
	fmt.Println("\n=== Example 3: Multiple Async Tasks ===")

	apiKey := os.Getenv("RUNNINGHUB_API_KEY")
	if apiKey == "" {
		fmt.Println("⚠️  Skipping - RUNNINGHUB_API_KEY not set")
		return
	}

	kit := comfykit.NewComfyKit(comfykit.WithRunningHubAPIKey(apiKey))

	type TaskInfo struct {
		Index        int
		Prompt       string
		TaskID       string
		OutputID2Var map[string]string
		Result       *comfyui.ExecuteResult
		Completed    bool
	}

	// Create multiple tasks
	prompts := []string{
		"a red rose in a vase",
	}

	var tasks []*TaskInfo
	fmt.Println("🔧 Creating", len(prompts), "tasks...")

	for i, prompt := range prompts {
		taskID, outputID2Var, err := kit.ExecuteAsyncByID("1988434426705133569", map[string]interface{}{
			"prompt": prompt,
		})
		if err != nil {
			fmt.Printf("⚠️  Failed to create task %d: %v\n", i, err)
			continue
		}
		tasks = append(tasks, &TaskInfo{
			Index:        i,
			Prompt:       prompt,
			TaskID:       taskID,
			OutputID2Var: outputID2Var,
			Completed:    false,
		})
		fmt.Printf("✅ Task %d: %s (ID: %s)\n", i, prompt[:20]+"...", taskID)
	}

	if len(tasks) == 0 {
		fmt.Println("⚠️  No tasks were created")
		return
	}

	// Poll all tasks until complete
	fmt.Println("⏳ Polling all tasks...")
	checkInterval := 2 * time.Second
	maxChecks := 300
	checkCount := 0
	remaining := len(tasks)

	for remaining > 0 && checkCount < maxChecks {
		checkCount++
		fmt.Printf("\nCheck %02d/%d: %d tasks remaining\n", checkCount, maxChecks, remaining)

		for _, task := range tasks {
			if task.Completed {
				continue
			}

			result, completed, err := kit.GetTaskCompletion(task.TaskID, task.OutputID2Var)
			if err != nil {
				fmt.Printf("   Task %d: error checking (%v)\n", task.Index, err)
				continue
			}

			if completed {
				task.Result = result
				task.Completed = true
				remaining--
				if result.Status == "completed" {
					fmt.Printf("   ✅ Task %d: complete, %d images\n", task.Index, len(result.Images))
				} else {
					fmt.Printf("   ❌ Task %d: %s - %s\n", task.Index, result.Status, result.Message)
				}
			} else {
				fmt.Printf("   ⏳ Task %d: running...\n", task.Index)
			}
		}

		if remaining > 0 {
			time.Sleep(checkInterval)
		}
	}

	// Summary
	fmt.Println("\n=== Summary ===")
	completedCount := 0
	errorCount := 0

	for _, task := range tasks {
		if task.Completed {
			if task.Result.Status == "completed" {
				completedCount++
			} else {
				errorCount++
			}
		}
	}

	fmt.Printf("Total tasks: %d\n", len(tasks))
	fmt.Printf("Successful: %d\n", completedCount)
	fmt.Printf("Errors: %d\n", errorCount)
}

func main() {
	fmt.Println("⚡ ComfyKit-Go Async Execution Examples")
	fmt.Println("=" + strings.Repeat("=", 59))

	// exampleAsyncTaskCreation()
	// examplePollForResults()
	exampleMultipleAsyncTasks()

	fmt.Println("\n" + "=" + strings.Repeat("=", 59))
	fmt.Println("✅ Async execution examples completed!")
}

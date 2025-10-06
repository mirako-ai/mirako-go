package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/mirako-ai/mirako-go/client"
	"github.com/mirako-ai/mirako-go/api"
)

func main() {
	apiKey := os.Getenv("MIRAKO_API_KEY")
	if apiKey == "" {
		log.Fatal("MIRAKO_API_KEY environment variable is required")
	}

	c, err := client.NewClient(
		client.WithAPIKey(apiKey),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	prompt := "A friendly robot with a smiling face"

	resp, err := c.GenerateAvatarAsync(ctx, api.GenerateAvatarAsyncJSONRequestBody{
		Prompt: prompt,
	})
	if err != nil {
		log.Fatalf("Failed to generate avatar: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("Unexpected status code: %d", resp.StatusCode)
	}

	var result api.AsyncGenerateAvatarApiResponseBody
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatalf("Failed to decode response: %v", err)
	}

	if result.Data == nil {
		log.Fatal("No data in response")
	}

	taskID := result.Data.TaskId
	fmt.Printf("Avatar generation started. Task ID: %s\n", taskID)
	fmt.Printf("Initial Status: %s\n", result.Data.Status)
	fmt.Println("\nNote: Use webhooks or poll the task status endpoint to check completion.")
	fmt.Println("The task will be processed asynchronously.")
}

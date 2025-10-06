package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/mirako-ai/mirako-go/client"
	"github.com/mirako-ai/mirako-go/gen"
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

	audioFile := "input.wav"
	audioData, err := os.ReadFile(audioFile)
	if err != nil {
		log.Fatalf("Failed to read audio file: %v", err)
	}

	audioB64 := base64.StdEncoding.EncodeToString(audioData)

	resp, err := c.ConvertSpeechToText(ctx, gen.ConvertSpeechToTextJSONRequestBody{
		Audio: audioB64,
	})
	if err != nil {
		log.Fatalf("Failed to convert speech to text: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("Unexpected status code: %d", resp.StatusCode)
	}

	var result gen.STTApiResponseBody
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatalf("Failed to decode response: %v", err)
	}

	if result.Data == nil {
		log.Fatal("No data in response")
	}

	fmt.Printf("Transcription: %s\n", result.Data.Transcription)
	fmt.Printf("Request ID: %s\n", result.Data.Id)
	if result.Data.InputDuration != nil {
		fmt.Printf("Duration: %.2f seconds\n", *result.Data.InputDuration)
	}
}

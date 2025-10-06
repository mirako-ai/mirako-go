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

	voiceProfileID := "your-voice-profile-id"
	text := "Hello, welcome to Mirako AI!"

	resp, err := c.ConvertTextToSpeech(ctx, gen.ConvertTextToSpeechJSONRequestBody{
		Text:           text,
		VoiceProfileId: voiceProfileID,
		ReturnType:     "b64_audio_str",
	})
	if err != nil {
		log.Fatalf("Failed to convert text to speech: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("Unexpected status code: %d", resp.StatusCode)
	}

	var result gen.TTSApiResponseBody
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatalf("Failed to decode response: %v", err)
	}

	if result.Data == nil || result.Data.B64AudioStr == nil {
		log.Fatal("No audio data in response")
	}

	audioData, err := base64.StdEncoding.DecodeString(*result.Data.B64AudioStr)
	if err != nil {
		log.Fatalf("Failed to decode base64 audio: %v", err)
	}

	if err := os.WriteFile("output.pcm", audioData, 0644); err != nil {
		log.Fatalf("Failed to write audio file: %v", err)
	}

	fmt.Println("Speech generated successfully: output.pcm")
	fmt.Printf("Request ID: %s\n", result.Data.Id)
	if result.Data.OutputDuration != nil {
		fmt.Printf("Duration: %.2f seconds\n", *result.Data.OutputDuration)
	}
}

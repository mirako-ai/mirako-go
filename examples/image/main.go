package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/mirako-ai/mirako-go/api"
	"github.com/mirako-ai/mirako-go/client"
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

	// Example 1: Synchronous Image Generation
	fmt.Println("=== Example 1: Synchronous Image Generation ===")
	generateImageSync(ctx, c)

	fmt.Println()

	// Example 2: Text-Image-to-Image Generation
	fmt.Println("=== Example 2: Text-Image-to-Image Generation ===")
	generateImageWithInputImages(ctx, c)
}

// generateImageSync demonstrates synchronous image generation
func generateImageSync(ctx context.Context, c *client.Client) {
	prompt := "A serene mountain landscape at sunset, with vibrant orange and purple skies"
	aspectRatio := api.GenerateImageApiRequestBodyAspectRatio("16:9")

	resp, err := c.GenerateImage(ctx, api.GenerateImageJSONRequestBody{
		Prompt:      prompt,
		AspectRatio: aspectRatio,
	})
	if err != nil {
		log.Fatalf("Failed to generate image: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("Unexpected status code: %d", resp.StatusCode)
	}

	var result api.GenerateImageApiResponseBody
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatalf("Failed to decode response: %v", err)
	}

	if result.Data == nil || result.Data.Image == nil {
		log.Fatal("No image data in response")
	}

	// Decode and save the image
	imageData, err := base64.StdEncoding.DecodeString(*result.Data.Image)
	if err != nil {
		log.Fatalf("Failed to decode base64 image: %v", err)
	}

	outputFile := "output_sync.jpg"
	if err := os.WriteFile(outputFile, imageData, 0644); err != nil {
		log.Fatalf("Failed to write image file: %v", err)
	}

	fmt.Printf("✓ Image generated successfully: %s\n", outputFile)
	fmt.Printf("  Prompt: %s\n", prompt)
	fmt.Printf("  Aspect Ratio: %s\n", aspectRatio)
}

// generateImageWithInputImages demonstrates text-image-to-image generation
func generateImageWithInputImages(ctx context.Context, c *client.Client) {
	// Read an input image (you need to provide this)
	inputImagePath := os.Getenv("INPUT_IMAGE_PATH")
	if inputImagePath == "" {
		fmt.Println("⚠ Skipping: INPUT_IMAGE_PATH environment variable not set")
		fmt.Println("  Set INPUT_IMAGE_PATH to use text-image-to-image generation")
		return
	}

	imageData, err := os.ReadFile(inputImagePath)
	if err != nil {
		log.Printf("Warning: Failed to read input image: %v", err)
		fmt.Println("⚠ Skipping text-image-to-image generation")
		return
	}

	// Encode image to base64 data URL
	imageB64 := base64.StdEncoding.EncodeToString(imageData)

	// Determine image format based on file extension
	format := "image/jpeg"
	if len(inputImagePath) > 4 && inputImagePath[len(inputImagePath)-4:] == ".png" {
		format = "image/png"
	}

	dataURL := fmt.Sprintf("data:%s;base64,%s", format, imageB64)

	label := "reference"
	labeledImage := api.LabeledImage{
		Data:  dataURL,
		Label: &label,
	}

	prompt := "Transform <reference> into a watercolor painting style with soft pastel colors"
	aspectRatio := api.GenerateImageApiRequestBodyAspectRatio("1:1")
	images := []api.LabeledImage{labeledImage}

	resp, err := c.GenerateImage(ctx, api.GenerateImageJSONRequestBody{
		Prompt:      prompt,
		AspectRatio: aspectRatio,
		Images:      &images,
	})
	if err != nil {
		log.Fatalf("Failed to generate image with input: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("Unexpected status code: %d", resp.StatusCode)
	}

	var result api.GenerateImageApiResponseBody
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatalf("Failed to decode response: %v", err)
	}

	if result.Data == nil || result.Data.Image == nil {
		log.Fatal("No image data in response")
	}

	// Decode and save the image
	outputImageData, err := base64.StdEncoding.DecodeString(*result.Data.Image)
	if err != nil {
		log.Fatalf("Failed to decode base64 image: %v", err)
	}

	outputFile := "output_image_to_image.jpg"
	if err := os.WriteFile(outputFile, outputImageData, 0644); err != nil {
		log.Fatalf("Failed to write image file: %v", err)
	}

	fmt.Printf("✓ Image generated successfully: %s\n", outputFile)
	fmt.Printf("  Input Image: %s\n", inputImagePath)
	fmt.Printf("  Prompt: %s\n", prompt)
	fmt.Printf("  Aspect Ratio: %s\n", aspectRatio)
	fmt.Printf("  Images Provided: %d\n", len(images))
}

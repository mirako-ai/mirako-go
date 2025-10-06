# Mirako Go SDK

Official Go SDK for [Mirako.ai](https://mirako.ai). Generate AI avatars, images, videos, and realistic speech with a simple and intuitive interface.

[![Go Reference](https://pkg.go.dev/badge/github.com/mirako-ai/mirako-go.svg)](https://pkg.go.dev/github.com/mirako-ai/mirako-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/mirako-ai/mirako-go)](https://goreportcard.com/report/github.com/mirako-ai/mirako-go)

## Features

- **Avatar Generation**: Create and build custom avatars
- **Image Generation**: Generate images with AI
- **Video Generation**: Create talking avatar and motion videos
- **Text-to-Speech (TTS)**: Convert text to natural-sounding speech
- **Speech-to-Text (STT)**: Transcribe audio to text
- **Interactive Sessions**: Manage real-time interactive avatar sessions
- **Voice Cloning**: Clone and manage custom voice profiles

## Installation

```bash
go get github.com/mirako-ai/mirako-go
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mirako-ai/mirako-go/client"
	"github.com/mirako-ai/mirako-go/api"
)

func main() {
	c, err := client.NewClient(
		client.WithAPIKey("your-api-key"),
	)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	
	resp, err := c.GetUserAvatarList(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	fmt.Printf("Status: %d\n", resp.StatusCode)
}
```

## Authentication

Get your API key from the [Mirako Dashboard](https://mirako.co/dashboard). Initialize the client with your API key:

```go
client, err := client.NewClient(
	client.WithAPIKey("mi-xxxxxxxxxxxxx"),
)
```

## Configuration Options

### Custom Base URL

```go
client, err := client.NewClient(
	client.WithAPIKey("your-api-key"),
	client.WithBaseURL("https://custom.mirako.co"),
)
```

### Custom HTTP Client

```go
httpClient := &http.Client{
	Timeout: 120 * time.Second,
}

client, err := client.NewClient(
	client.WithAPIKey("your-api-key"),
	client.WithHTTPClient(httpClient),
)
```

### Retry Configuration

```go
retryConfig := client.DefaultRetryConfig()
retryConfig.MaxRetries = 5

client, err := client.NewClient(
	client.WithAPIKey("your-api-key"),
	client.WithRetry(retryConfig),
)
```

### Logging

```go
type MyLogger struct{}

func (l *MyLogger) Logf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

client, err := client.NewClient(
	client.WithAPIKey("your-api-key"),
	client.WithLogger(&MyLogger{}),
)
```

### Tracing

```go
type MyTracer struct{}

func (t *MyTracer) TraceRequest(ctx context.Context, req *http.Request) {
	// Add OpenTelemetry or custom tracing
}

func (t *MyTracer) TraceResponse(ctx context.Context, resp *http.Response) {
	// Trace response
}

client, err := client.NewClient(
	client.WithAPIKey("your-api-key"),
	client.WithTracer(&MyTracer{}),
)
```

## Usage Examples

### Text-to-Speech

```go
ctx := context.Background()

body := api.ConvertTextToSpeechJSONRequestBody{
	Text:            "Hello, welcome to Mirako!",
	VoiceProfileId:  "profile-id",
	ReturnType:      "b64_audio_str",
}

resp, err := client.ConvertTextToSpeech(ctx, body)
if err != nil {
	log.Fatal(err)
}
defer resp.Body.Close()
```

### Speech-to-Text

```go
ctx := context.Background()

audioBase64 := "UklGRiQAAABXQVZFZm10..." // base64 encoded audio

body := api.ConvertSpeechToTextJSONRequestBody{
	Audio: audioBase64,
}

resp, err := client.ConvertSpeechToText(ctx, body)
if err != nil {
	log.Fatal(err)
}
defer resp.Body.Close()
```

### Generate Avatar (Async)

```go
ctx := context.Background()

prompt := "A realistic photo of an Asian girl with long black hair"
body := api.GenerateAvatarAsyncJSONRequestBody{
	Prompt: prompt,
}

resp, err := client.GenerateAvatarAsync(ctx, body)
if err != nil {
	log.Fatal(err)
}
defer resp.Body.Close()

// Parse response to get task_id
// Then poll status with GetAvatarGenerationStatus
```

### Generate Image (Async)

```go
ctx := context.Background()

aspectRatio := api.N169
body := api.GenerateImageAsyncJSONRequestBody{
	Prompt:      "A girl sitting in a cafe smiling",
	AspectRatio: aspectRatio,
}

resp, err := client.GenerateImageAsync(ctx, body)
if err != nil {
	log.Fatal(err)
}
defer resp.Body.Close()
```

### Start Interactive Session

```go
ctx := context.Background()

body := api.StartInteractiveSessionJSONRequestBody{
	AvatarId:       "avatar-id",
	VoiceProfileId: "voice-profile-id",
	LlmModel:       "gemini-2.0-flash",
	Instruction:    "You are a helpful assistant.",
}

resp, err := client.StartInteractiveSession(ctx, body)
if err != nil {
	log.Fatal(err)
}
defer resp.Body.Close()
```

### List Avatars

```go
ctx := context.Background()

resp, err := client.GetUserAvatarList(ctx)
if err != nil {
	log.Fatal(err)
}
defer resp.Body.Close()
```

## Error Handling

The SDK returns standard Go errors. HTTP errors follow the OpenAPI error model:

```go
resp, err := client.GetAvatarById(ctx, "avatar-id")
if err != nil {
	log.Fatal(err)
}

if resp.StatusCode != http.StatusOK {
	body, _ := io.ReadAll(resp.Body)
	log.Printf("Error: %s", body)
}
```

## API Reference

See the [official API documentation](https://mirako.co/docs) for detailed information about request/response schemas and parameters.

All generated types and constants are available in the `api` package:

```go
import "github.com/mirako-ai/mirako-go/api"

// Use generated types
var status api.AsyncTaskStatusStatus = api.AsyncTaskStatusStatusCOMPLETED
```

## Code Generation

This SDK is automatically generated from the official OpenAPI specification using [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen).

To regenerate the client code:

```bash
make generate
```

## Development

### Prerequisites

- Go 1.21 or higher
- oapi-codegen v2

### Build

```bash
make build
```

### Test

```bash
make test
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Links

- [Mirako Website](https://mirako.ai)
- [API Documentation](https://docs.mirako.ai)
- [Dashboard](https://developer.mirako.ai)
- Discord: [Join our Discord](https://discord.gg/Wkxxr54y)

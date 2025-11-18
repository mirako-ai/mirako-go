# Mirako Go SDK Examples

This directory contains example code demonstrating how to use the Mirako Go SDK.

## Running Examples

Each example is in its own directory. To run an example:

```bash
# Set your API key
export MIRAKO_API_KEY=your-api-key-here

# Run the example
cd tts && go run main.go
```

## Available Examples

### Text-to-Speech (TTS)
- **Directory**: `tts/`
- **Description**: Convert text to speech audio
- **Output**: PCM audio file

### Speech-to-Text (STT)
- **Directory**: `stt/`
- **Description**: Transcribe audio to text
- **Input**: WAV audio file

### Avatar Generation
- **Directory**: `avatar/`
- **Description**: Generate an avatar image from a text prompt (async)
- **Output**: Task ID for tracking generation status

### Image Generation
- **Directory**: `image/`
- **Description**: Generate images using text prompts (synchronous) and text-image-to-image generation
- **Features**:
  - Synchronous image generation from text prompts
  - Text-image-to-image generation with labeled input images (up to 5 images)
  - Customizable aspect ratios
- **Output**: JPG image files
- **Environment Variables**:
  - `MIRAKO_API_KEY`: Required for authentication
  - `INPUT_IMAGE_PATH`: Optional, path to input image for text-image-to-image generation

## Prerequisites

- Go 1.21 or later
- Mirako API key (get one at https://mirako.co)

## Notes

- Some operations (like avatar generation) are asynchronous and return a task ID
- Use webhooks or implement polling to check task completion status
- Audio files may need to be in specific formats depending on the API endpoint

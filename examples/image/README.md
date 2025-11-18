# Image Generation Examples

This example demonstrates the Mirako API's image generation capabilities.

## Features

### 1. Synchronous Image Generation
Generate images directly from text prompts with immediate results.

### 2. Text-Image-to-Image Generation
Transform existing images using text prompts and labeled reference images.

## Usage

### Basic Image Generation

```bash
export MIRAKO_API_KEY=your-api-key-here
go run main.go
```

This will:
- Generate an image from a text prompt
- Save the result as `output_sync.jpg`

### Text-Image-to-Image Generation

```bash
export MIRAKO_API_KEY=your-api-key-here
export INPUT_IMAGE_PATH=/path/to/your/image.jpg
go run main.go
```

This will:
- Use your input image as a reference
- Transform it according to the prompt
- Save the result as `output_image_to_image.jpg`

## Input Image Format

Supported formats for input images:
- JPEG (`.jpg`, `.jpeg`)
- PNG (`.png`)

The example automatically detects the format and creates the appropriate data URL.

## Configuration Options

### Aspect Ratios
Available aspect ratios:
- `16:9` - Widescreen (default)
- `1:1` - Square
- `4:3` - Standard
- `3:2` - Classic photo
- `2:3` - Portrait
- `3:4` - Portrait
- `9:16` - Vertical

### Multiple Input Images
You can provide up to 5 labeled images:

```go
images := []api.LabeledImage{
    {Data: dataURL1, Label: &label1},
    {Data: dataURL2, Label: &label2},
}
```

Reference them in your prompt using `<label>`:
```go
prompt := "Combine <style> and <content> into a surreal artwork"
```

## Example Prompts

### Synchronous Generation
```
"A serene mountain landscape at sunset, with vibrant orange and purple skies"
"A futuristic city with flying cars and neon lights"
"A close-up portrait of a smiling golden retriever"
```

### Text-Image-to-Image
```
"Transform <reference> into a watercolor painting style"
"Apply <style> artistic style to <content>"
"Make <photo> look like a vintage 1970s photograph"
```

## API Details

This example uses:
- **Endpoint**: `/v1/image/generate` (synchronous)
- **Request**: `GenerateImageApiRequestBody`
  - `prompt`: Text description (required)
  - `aspect_ratio`: Output dimensions (required)
  - `images`: Array of labeled images (optional, max 5)
  - `seed`: Random seed for reproducibility (optional)
- **Response**: `GenerateImageApiResponseBody`
  - Returns base64-encoded JPG image immediately

## Notes

- This is a **synchronous** API call - it waits for the image to be generated
- For asynchronous generation, use the `/v1/image/async_generate` endpoint
- Generation time varies based on complexity and server load
- Images are returned as base64-encoded JPG format
- Input images must be provided as data URLs with base64 encoding

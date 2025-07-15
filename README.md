## Go Auto1111 Launcher

This program is a Go-based launcher for AUTOMATIC1111's Stable Diffusion WebUI. It loads configuration from a `.env` file, prints the settings, launches the Stable Diffusion Web UI, and sends an API request to create an image with AI-powered prompt generation and enhancement.

### Usage

1. **Configure your `.env` file** in the project root. Example:

    ```env
    SAMPLER_NAME="DPM++ 2M SDE"
    SCHEDULER_NAME="Karras"
    IMAGE_HEIGHT=1024
    IMAGE_WIDTH=1024
    CFG_SCALE=4.5
    STEPS=30
    AUTO1111_BAT="C:\\Users\\user\\auto1111\\run.bat"
    ENVIRONMENT_BAT="C:\\Users\\user\\auto1111\\environment.bat"
    GEMINI_API_KEY="your_gemini_api_key_here"
    ```
    **Note:** On Windows, use double backslashes (`\\`) in all file paths in your `.env` file.

2. **Build and run the program:**
    ```sh
    go run .
    ```
    or
    ```sh
    go build -o auto1111.exe
    ./auto1111.exe
    ```

### Prompt Generation Features

The launcher includes AI-powered prompt generation using Google's Gemini API:

- **Automatic Generation**: Press Enter when prompted to generate a completely new creative prompt
- **Prompt Enhancement**: Type your basic idea (e.g., "a cat") and the AI will enhance it with artistic style, lighting, and quality terms
- **Examples**:
  - Input: "a cat" -> Enhanced: "a cat, digital art, soft lighting, highly detailed, photorealistic"
  - Input: (empty) -> Generated: "A majestic dragon soaring through storm clouds, digital art, cinematic lighting, highly detailed"

### How It Works

The launcher will:
- Load your configuration settings
- Prompt you to enter an image generation prompt (or generate one automatically)
- Enhance your prompt using AI if you provide input, or generate a new one if empty
- Start the AUTOMATIC1111 WebUI if it's not already running
- Generate your image using the txt2img API
- Save the generated image to your specified location

### Requirements

- Go 1.19 or later
- AUTOMATIC1111 Stable Diffusion WebUI installed
- Google Gemini API key (for prompt generation)

### Future Enhancement Ideas

- **Prompt History**: Save and recall previously used prompts
- **Style Presets**: Pre-defined artistic styles (anime, photorealistic, oil painting, etc.)
- **Batch Generation**: Generate multiple images with variations of the same prompt
- **Image-to-Prompt**: Analyze existing images to generate similar prompts
- **Negative Prompt Generation**: AI-generated negative prompts to improve image quality
- **Model Selection**: Support for different Stable Diffusion models and checkpoints
- **GUI Interface**: Web-based or desktop GUI for easier interaction
- **Prompt Templates**: Customizable templates for different image types (portraits, landscapes, etc.)
- **Quality Scoring**: AI evaluation of generated images with suggestions for improvement
- **Social Features**: Share prompts and generated images with the community

---
**Troubleshooting:**
- Ensure all paths in `.env` use double backslashes.
- The batch file will be launched in its own directory so that relative paths work as expected.
- Make sure your GEMINI_API_KEY is valid for prompt generation to work.
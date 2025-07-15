## Go Auto1111 Launcher

This program is a Go-based launcher for AUTOMATIC1111's Stable Diffusion WebUI. It loads configuration from a `.env` file, prints the settings, launches the Stable Diffusion Web UI, and sends an API request to create an image.

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

The launcher will:
- Load your configuration settings
- Prompt you to enter an image generation prompt
- Start the AUTOMATIC1111 WebUI if it's not already running
- Generate your image using the txt2img API
- Save the generated image to your specified location

---
**Troubleshooting:**
- Ensure all paths in `.env` use double backslashes.
- The batch file will be launched in its own directory so that relative paths work as expected.

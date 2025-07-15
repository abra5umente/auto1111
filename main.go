package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	err = godotenv.Load(filepath.Join(exPath, ".env"))
	if err != nil {
		fmt.Println("No .env file found, using system environment variables.")
	}
	sampler := os.Getenv("SAMPLER_NAME")
	scheduler := os.Getenv("SCHEDULER_NAME")
	width := os.Getenv("IMAGE_WIDTH")
	height := os.Getenv("IMAGE_HEIGHT")
	steps := os.Getenv("STEPS")
	cfgScale := os.Getenv("CFG_SCALE")
	fmt.Println("Using the following settings from .env file:")
	fmt.Println("Sampler:", sampler)
	fmt.Println("Scheduler:", scheduler)
	fmt.Println("Steps:", steps)
	fmt.Println("Width:", width)
	fmt.Println("Height:", height)
	fmt.Println("CFG Scale:", cfgScale)

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter your prompt (or press Enter to generate one automatically):")
	userPrompt, _ := reader.ReadString('\n')
	userPrompt = strings.TrimSpace(userPrompt)

	if userPrompt == "" {
		fmt.Println("Generating a creative prompt...")
		generatedPrompt, err := generate_prompt("")
		if err != nil {
			fmt.Println("Failed to generate prompt:", err)
			os.Exit(1)
		}
		userPrompt = generatedPrompt
		fmt.Println("Generated prompt:", userPrompt)
	} else {
		fmt.Println("Enhancing your prompt...")
		enhancedPrompt, err := generate_prompt(userPrompt)
		if err != nil {
			fmt.Printf("Failed to enhance prompt, using original: %v\n", err)
		} else {
			userPrompt = enhancedPrompt
			fmt.Println("Enhanced prompt:", userPrompt)
		}
	}

	// build the payload
	payload := map[string]interface{}{
		"prompt":    userPrompt,
		"sampler":   sampler,
		"scheduler": scheduler,
		"width":     width,
		"height":    height,
		"steps":     steps,
		"cfg_scale": cfgScale,
	}

	var cmd *exec.Cmd
	if !check_auto1111() {
		fmt.Println("auto1111 is not running, starting.")
		cmd = start_auto1111()
		fmt.Println("waiting for auto1111 to be ready...")
		for {
			time.Sleep(2 * time.Second)
			if check_auto1111() {
				break
			}
		}
	}
	// create the image
	if !create_image(payload) {
		fmt.Println("something went wrong making the image, exiting.")
		if cmd != nil {
			stop_auto1111(cmd)
		}
		os.Exit(1)
	}
	if cmd != nil {
		stop_auto1111(cmd)
	}
}


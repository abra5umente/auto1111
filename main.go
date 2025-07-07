package main

import (
	"bufio" // for reading user input
	"fmt" // printing to console
	"os" // environment variables, file operations
	"strings" // string manipulation
	"time" // sleeping, timestamps
	"github.com/joho/godotenv" // loading environment variables from .env file
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file found, using system environment variables.")
	}
	sampler := os.Getenv("SAMPLER_NAME")
	scheduler := os.Getenv("SCHEDULER_NAME")
	width := os.Getenv("IMAGE_WIDTH")
	height := os.Getenv("IMAGE_HEIGHT")
	steps := os.Getenv("STEPS")
	fmt.Println("Using the following settings from .env file:")
	fmt.Println("Sampler:", sampler)
	fmt.Println("Scheduler:", scheduler)
	fmt.Println("Steps:", steps)
	fmt.Println("Width:", width)
	fmt.Println("Height:", height)

	// prompt user for input
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("enter your prompt (or leave empty to let the rock think for you, which will be weird):")
	userPrompt, _ := reader.ReadString('\n')
	userPrompt = strings.TrimSpace(userPrompt)

	// build the payload
	payload := map[string]interface{}{
		"prompt":    userPrompt,
		"sampler":   sampler,
		"scheduler": scheduler,
		"width":     width,
		"height":    height,
		"steps":     steps,
	}

// start the AUTOMATIC1111 server if not running
if !check_auto1111() {
	fmt.Println("auto1111 is not running, starting.")
	start_auto1111()
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
	stop_auto1111()
	os.Exit(1)
}
stop_auto1111()
}
package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func generate_prompt(user_prompt string) (string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY environment variable not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return "", fmt.Errorf("failed to create Gemini client: %v", err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")
	model.SetTemperature(0.8)
	model.SetMaxOutputTokens(100)

	var prompt string
	if user_prompt == "" {
		prompt = `Generate a creative prompt for Stable Diffusion image generation. Include subject, art style, lighting, and quality terms. Return only the prompt text, nothing else.

Examples:
- A majestic dragon soaring through storm clouds, digital art, cinematic lighting, highly detailed
- Cozy coffee shop interior, warm lighting, photorealistic, detailed textures
- Futuristic cityscape at night, neon lights, cyberpunk style, atmospheric lighting

Generate a new creative prompt:`
	} else {
		prompt = fmt.Sprintf(`Enhance this Stable Diffusion prompt by adding artistic style, lighting, and quality terms: "%s"

Return only the enhanced prompt text, nothing else. Keep the original concept but make it more detailed and artistic.

Enhanced prompt:`, user_prompt)
	}

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("Gemini API call failed: %v", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("empty response from Gemini")
	}

	if txt, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
		rawResponse := string(txt)
		result := cleanPromptResponse(rawResponse)
		if result == "" {
			if user_prompt == "" {
				result = "A beautiful landscape at sunset, digital art, cinematic lighting, highly detailed"
			} else {
				result = user_prompt + ", digital art, highly detailed"
			}
		}
		return result, nil
	}

	return "", fmt.Errorf("failed to extract text from Gemini response")
}

func cleanPromptResponse(response string) string {
	// Simply trim whitespace and return the whole response
	return strings.TrimSpace(response)
}
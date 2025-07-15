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
	fmt.Println("DEBUG: Starting generate_prompt function")
	
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY environment variable not set")
	}
	fmt.Printf("DEBUG: API key found (length: %d)\n", len(apiKey))

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	fmt.Println("DEBUG: Creating Gemini client...")
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return "", fmt.Errorf("failed to create Gemini client: %v", err)
	}
	defer client.Close()
	fmt.Println("DEBUG: Client created successfully")

	model := client.GenerativeModel("gemini-1.5-flash")
	model.SetTemperature(0.8)
	model.SetMaxOutputTokens(100)
	fmt.Println("DEBUG: Model configured")

	var prompt string
	if user_prompt == "" {
		fmt.Println("DEBUG: Generating new prompt")
		prompt = `Generate a creative prompt for Stable Diffusion image generation. Include subject, art style, lighting, and quality terms. Return only the prompt text, nothing else.

Examples:
- A majestic dragon soaring through storm clouds, digital art, cinematic lighting, highly detailed
- Cozy coffee shop interior, warm lighting, photorealistic, detailed textures
- Futuristic cityscape at night, neon lights, cyberpunk style, atmospheric lighting

Generate a new creative prompt:`
	} else {
		fmt.Printf("DEBUG: Enhancing existing prompt: %s\n", user_prompt)
		prompt = fmt.Sprintf(`Enhance this Stable Diffusion prompt by adding artistic style, lighting, and quality terms: "%s"

Return only the enhanced prompt text, nothing else. Keep the original concept but make it more detailed and artistic.

Enhanced prompt:`, user_prompt)
	}

	fmt.Println("DEBUG: Calling Gemini API...")
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("Gemini API call failed: %v", err)
	}
	fmt.Printf("DEBUG: API response received, candidates: %d\n", len(resp.Candidates))

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("empty response from Gemini")
	}

	if txt, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
		rawResponse := string(txt)
		fmt.Printf("DEBUG: Raw response: %q\n", rawResponse)
		result := cleanPromptResponse(rawResponse)
		fmt.Printf("DEBUG: Cleaned result: %q\n", result)
		if result == "" {
			if user_prompt == "" {
				result = "A beautiful landscape at sunset, digital art, cinematic lighting, highly detailed"
				fmt.Println("DEBUG: Using fallback prompt due to empty result")
			} else {
				result = user_prompt + ", digital art, highly detailed"
				fmt.Println("DEBUG: Using enhanced fallback due to empty result")
			}
		}
		return result, nil
	}

	return "", fmt.Errorf("failed to extract text from Gemini response")
}

func cleanPromptResponse(response string) string {
	response = strings.TrimSpace(response)
	
	if response == "" {
		return response
	}
	
	response = strings.TrimPrefix(response, "Here's ")
	response = strings.TrimPrefix(response, "Here is ")
	response = strings.TrimPrefix(response, "Sure! ")
	response = strings.TrimPrefix(response, "Sure, ")
	
	if strings.HasPrefix(response, "\"") && strings.HasSuffix(response, "\"") {
		response = response[1 : len(response)-1]
	}
	
	lines := strings.Split(response, "\n")
	var bestLine string
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		lowerLine := strings.ToLower(line)
		if strings.HasPrefix(lowerLine, "here") ||
		   strings.HasPrefix(lowerLine, "sure") ||
		   strings.HasPrefix(lowerLine, "i'll") ||
		   strings.HasPrefix(lowerLine, "let me") ||
		   strings.HasPrefix(lowerLine, "this prompt") ||
		   strings.HasPrefix(line, "**") ||
		   strings.HasPrefix(line, "*") ||
		   strings.Contains(lowerLine, "enhanced prompt:") ||
		   strings.Contains(lowerLine, "new prompt:") {
			continue
		}
		
		if bestLine == "" {
			bestLine = line
		}
	}
	
	if bestLine != "" {
		return bestLine
	}
	
	return response
}
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func loadSettings(path string) (map[string]interface{}, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open settings file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read settings file: %w", err)
	}

	var settings map[string]interface{}
	if err := json.Unmarshal(data, &settings); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return settings, nil
}

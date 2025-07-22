package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// -----------------------------------------------------------------------------
// Key‑lookup helpers that understand **both** lowercase and legacy ALL‑CAPS keys
// -----------------------------------------------------------------------------

var keyAlias = map[string][]string{
	"sampler_name": {"sampler_name", "SAMPLER_NAME"},
	"scheduler":    {"scheduler", "SCHEDULER_NAME"},
	"width":        {"width", "IMAGE_WIDTH"},
	"height":       {"height", "IMAGE_HEIGHT"},
	"steps":        {"steps", "STEPS"},
	"cfg_scale":    {"cfg_scale", "CFG_SCALE"},
}

func lookup(settings map[string]interface{}, canonical string) (interface{}, bool) {
	for _, k := range keyAlias[canonical] {
		if v, ok := settings[k]; ok {
			return v, true
		}
	}
	return nil, false
}

func strFrom(settings map[string]interface{}, canonical, fallback string) string {
	if v, ok := lookup(settings, canonical); ok {
		return fmt.Sprint(v)
	}
	return fallback
}

func intFrom(settings map[string]interface{}, canonical string, fallback int) int {
	if v, ok := lookup(settings, canonical); ok {
		switch vv := v.(type) {
		case float64:
			return int(vv)
		case int:
			return vv
		case json.Number:
			i, _ := vv.Int64()
			return int(i)
		}
	}
	return fallback
}

func floatFrom(settings map[string]interface{}, canonical string, fallback float64) float64 {
	if v, ok := lookup(settings, canonical); ok {
		switch vv := v.(type) {
		case float64:
			return vv
		case json.Number:
			f, _ := vv.Float64()
			return f
		}
	}
	return fallback
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

func firstNonZero(vals ...int) int {
	for _, v := range vals {
		if v != 0 {
			return v
		}
	}
	return 0
}

func firstNonNeg(vals ...float64) float64 {
	for _, v := range vals {
		if v >= 0 {
			return v
		}
	}
	return 0
}

func main() {
	// Locate exe directory (so .env & settings.json live beside it)
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)

	_ = godotenv.Load(filepath.Join(exPath, ".env")) // non‑fatal if missing

	// ---------------- flags ----------------
	settingsPath := flag.String("settings", "settings.json", "Path to settings JSON")
	promptFlag := flag.String("prompt", "", "Prompt text")
	outputFlag := flag.String("output", "", "Output PNG path")

	samplerFlag := flag.String("sampler_name", "", "Override sampler")
	schedulerFlag := flag.String("scheduler", "", "Override scheduler")
	stepsFlag := flag.Int("steps", 0, "Override steps")
	widthFlag := flag.Int("width", 0, "Override width")
	heightFlag := flag.Int("height", 0, "Override height")
	cfgFlag := flag.Float64("cfg_scale", -1.0, "Override CFG scale")

	flag.Parse()

	// ---------------- load settings ----------------
	settings, err := loadSettings(*settingsPath)
	if err != nil {
		fmt.Println("Failed to load settings.json:", err)
		os.Exit(1)
	}

	// ---------------- resolve parameters ----------------
	sampler := firstNonEmpty(*samplerFlag, strFrom(settings, "sampler_name", "Euler"))
	scheduler := firstNonEmpty(*schedulerFlag, strFrom(settings, "scheduler", "default"))
	width := firstNonZero(*widthFlag, intFrom(settings, "width", 512))
	height := firstNonZero(*heightFlag, intFrom(settings, "height", 512))
	steps := firstNonZero(*stepsFlag, intFrom(settings, "steps", 30))
	cfg := firstNonNeg(*cfgFlag, floatFrom(settings, "cfg_scale", 7.0))

	fmt.Printf("[DEBUG] sampler=%s scheduler=%s w=%d h=%d steps=%d cfg=%0.1f\n",
		sampler, scheduler, width, height, steps, cfg)

	// ---------------- prompt ----------------
	var userPrompt string
	if *promptFlag != "" {
		userPrompt = *promptFlag
	} else {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter prompt (blank = auto): ")
		userPrompt, _ = reader.ReadString('\n')
		userPrompt = strings.TrimSpace(userPrompt)
	}

	// Prompt enhancement
	var enhanced string
	if p, err := generate_prompt(userPrompt); err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] prompt enhancement failed: %v\n", err)
		enhanced = userPrompt
	} else {
		fmt.Fprintf(os.Stderr, "[DEBUG] enhanced prompt: %q\n", p)
		enhanced = p
	}
	userPrompt = enhanced

	// ---------------- payload ----------------
	payload := map[string]interface{}{
		"prompt":       userPrompt,
		"sampler_name": sampler,
		"scheduler":    scheduler,
		"width":        width,
		"height":       height,
		"steps":        steps,
		"cfg_scale":    cfg,
	}
	if *outputFlag != "" {
		payload["output"] = *outputFlag
	}

	// ---------------- ensure auto1111 running -------------
	var cmd *exec.Cmd
	if !check_auto1111() {
		fmt.Println("Starting auto1111…")
		cmd = start_auto1111()
		for !check_auto1111() {
			time.Sleep(1 * time.Second)
		}
	}

	// ---------------- generate image ----------------------
	outPath, ok := create_image(payload)
	if !ok {
		fmt.Println("Image generation failed")
		if cmd != nil {
			stop_auto1111(cmd)
		}
		os.Exit(1)
	}
	fmt.Println("Image generated at:", outPath)

	if cmd != nil {
		stop_auto1111(cmd)
	}
}

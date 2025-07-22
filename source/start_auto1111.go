package main

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "time"
    "io"
)

// start_auto1111 launches AUTOMATIC1111 using the BAT file given in the
// AUTO1111_BAT environment variable (or a default beside the exe).
// It propagates the current process environment (including variables
// loaded via godotenv) and streams the child's stdout/stderr so you can
// see exactly what WebUI is doing.
func start_auto1111() *exec.Cmd {
    // ------------------------------------------------------------------
    // 1. Locate the BAT file that sets up the correct Python environment
    // ------------------------------------------------------------------
    a1Bat := os.Getenv("AUTO1111_BAT")
    if a1Bat == "" {
        // Fallback: assume <exeDir>/start_auto1111.bat
        exe, _ := os.Executable()
        a1Bat = filepath.Join(filepath.Dir(exe), "start_auto1111.bat")
    }

    if _, err := os.Stat(a1Bat); err != nil {
        fmt.Println("[ERROR] AUTO1111 BAT not found:", a1Bat)
        return nil
    }

    fmt.Println("[INFO] Launching AUTOMATIC1111 via", a1Bat)

    // ------------------------------------------------------------------
    // 2. Prepare the *cmd /C <bat>* process
    // ------------------------------------------------------------------
    cmd := exec.Command("cmd", "/C", a1Bat)

    // Work dir = BAT's folder so any relative paths inside it resolve.
    cmd.Dir = filepath.Dir(a1Bat)

    // Inherit ALL env vars from the current process (godotenv‑loaded ones too).
    cmd.Env = os.Environ()

    // Stream WebUI's console directly to ours.
    cmd.Stdout = io.Discard
    cmd.Stderr = io.Discard

    // Small debug: confirm ENVIRONMENT_BAT is visible to the child.
    fmt.Println("[DEBUG] ENVIRONMENT_BAT =", os.Getenv("ENVIRONMENT_BAT"))

    // ------------------------------------------------------------------
    // 3. Launch
    // ------------------------------------------------------------------
    if err := cmd.Start(); err != nil {
        fmt.Println("[ERROR] Failed to start auto1111:", err)
        return nil
    }

    // Give the process a head start before we begin polling.
    time.Sleep(2 * time.Second)
    return cmd
}

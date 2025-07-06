package main

import (
	"fmt"           // printing to console
	"os"            // for os.Getenv
	"os/exec"       // starting/stopping .bat files
	"path/filepath" // for getting directory of .bat file

	"github.com/joho/godotenv" // loading environment variables from .env file
)

func start_auto1111() {
	fmt.Println("starting auto1111...")
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file found, using system environment variables.")
	}

	a1_bat := os.Getenv("AUTO1111_BAT")
	fmt.Println("AUTO1111_BAT:", a1_bat)
	if a1_bat == "" {
		fmt.Println("AUTO1111_BAT environment variable is not set.")
		return
	}

	fmt.Println("About to run AUTO1111_BAT...")
	cmd := exec.Command("cmd", "/C", a1_bat)
	cmd.Dir = filepath.Dir(a1_bat)
	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		fmt.Println("Failed to start auto1111:", err)
		os.Exit(1)
	}
	fmt.Println("Started AUTOMATIC1111 server successfully.")
}

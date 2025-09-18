package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {
	// Simply run the HTTP bridge
	cmd := exec.Command("go", "run", "./cmd/http-bridge")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to run HTTP bridge: %v", err)
	}
}

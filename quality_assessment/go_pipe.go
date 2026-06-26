package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func main() {
	audio_in := `01 - Polyphia - Playing God.m4a`
	audio_out := `Playing God.mp3`
	const AUDIO_ROOT = `./examples/audio/`
	// Run ffmpeg with arguments
	cmd := exec.Command("cmd", "/c", "ffmpeg.exe", "-i", AUDIO_ROOT+audio_in, "-b:a", "128k", AUDIO_ROOT+audio_out)

	// Get pipes for stdout, stderr, and stdin
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	stdin, _ := cmd.StdinPipe()

	// Start the command
	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting command:", err)
		return
	}

	// Create goroutines to stream output
	go streamOutput(stdout, "STDOUT")
	go streamOutput(stderr, "STDERR")

	// Monitor user input in case ffmpeg prompts for it
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("Enter input: ")
			input, _ := reader.ReadString('\n')
			_, _ = stdin.Write([]byte(input)) // Send user input to ffmpeg
		}
	}()

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		fmt.Println("Error:", err)
	}
}

// Helper function to stream command output
func streamOutput(pipe io.ReadCloser, label string) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		fmt.Printf("[%s] %s\n", label, scanner.Text())

		// Example: Detecting a yes/no prompt
		if strings.Contains(scanner.Text(), "[y/N]") {
			fmt.Println("Detected prompt! Please enter input.")
		}
	}
}

package installer

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

type DevboxManager struct{}

func NewDevboxManager() *DevboxManager {
	return &DevboxManager{}
}

func (d *DevboxManager) Name() string {
	return "devbox"
}

func (d *DevboxManager) IsAvailable() bool {
	return IsCommandAvailable("devbox")
}

func (d *DevboxManager) Install(packages ...string) error {
	args := append([]string{"global", "add"}, packages...)
	cmd := exec.Command("devbox", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Running: devbox global add %v\n", packages)
	return cmd.Run()
}

// InstallDevbox installs devbox on the system
func InstallDevbox() error {
	fmt.Println("Installing devbox...")

	// Download the installation script using Go's HTTP client
	resp, err := http.Get("https://get.jetify.com/devbox")
	if err != nil {
		return fmt.Errorf("failed to download devbox installer: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download devbox installer: HTTP %d", resp.StatusCode)
	}

	// Read the script content
	script, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read installer script: %w", err)
	}

	// Execute the script with bash
	cmd := exec.Command("bash", "-s")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Provide the script content as stdin
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start bash: %w", err)
	}

	// Write the script to stdin
	if _, err := stdin.Write(script); err != nil {
		return fmt.Errorf("failed to write script to bash: %w", err)
	}
	stdin.Close()

	// Wait for the command to complete
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("devbox installation failed: %w", err)
	}

	return nil
}

// InitDevboxShell initializes devbox shell environment
func InitDevboxShell() error {
	fmt.Println("Initializing devbox shell environment...")
	// This prints the shell initialization commands
	cmd := exec.Command("devbox", "global", "shellenv", "--init-hook")
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	fmt.Println("Add the following to your shell RC file:")
	fmt.Println(string(output))
	return nil
}

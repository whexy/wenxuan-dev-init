package installer

import (
	"fmt"
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
	cmd := exec.Command("bash", "-c", "curl -fsSL https://get.jetify.com/devbox | bash")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
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

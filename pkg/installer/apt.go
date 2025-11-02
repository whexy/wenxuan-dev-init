package installer

import (
	"fmt"
	"os"
	"os/exec"
)

type AptManager struct{}

func NewAptManager() *AptManager {
	return &AptManager{}
}

func (a *AptManager) Name() string {
	return "apt"
}

func (a *AptManager) IsAvailable() bool {
	return IsCommandAvailable("apt-get")
}

func (a *AptManager) Install(packages ...string) error {
	// Separate packages that need special handling
	var standardPackages []string
	var needsGH, needs1Password, needsChezmoi bool

	for _, pkg := range packages {
		switch pkg {
		case "gh":
			needsGH = true
		case "1password-cli":
			needs1Password = true
		case "chezmoi":
			needsChezmoi = true
		default:
			standardPackages = append(standardPackages, pkg)
		}
	}

	// Install standard packages first
	if len(standardPackages) > 0 {
		// Update package list first
		updateCmd := exec.Command("sudo", "apt-get", "update")
		updateCmd.Stdout = os.Stdout
		updateCmd.Stderr = os.Stderr
		fmt.Println("Running: sudo apt-get update")
		if err := updateCmd.Run(); err != nil {
			return fmt.Errorf("failed to update package list: %w", err)
		}

		// Install standard packages
		args := append([]string{"apt-get", "install", "-y"}, standardPackages...)
		cmd := exec.Command("sudo", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		fmt.Printf("Running: sudo apt-get install -y %v\n", standardPackages)
		if err := cmd.Run(); err != nil {
			fmt.Printf("Warning: Some standard packages failed to install: %v\n", err)
		}
	}

	// Install GitHub CLI with repository setup
	if needsGH {
		fmt.Println("\n📦 Installing GitHub CLI (requires repository setup)...")
		if err := InstallGitHubCLI(); err != nil {
			return fmt.Errorf("failed to install GitHub CLI: %w", err)
		}
	}

	// Install 1Password CLI with repository setup
	if needs1Password {
		fmt.Println("\n📦 Installing 1Password CLI (requires repository setup)...")
		if err := Install1PasswordCLI(); err != nil {
			return fmt.Errorf("failed to install 1Password CLI: %w", err)
		}
	}

	// Install chezmoi with official installer
	if needsChezmoi {
		fmt.Println("\n📦 Installing chezmoi (using official installer)...")
		if err := InstallChezmoi(); err != nil {
			return fmt.Errorf("failed to install chezmoi: %w", err)
		}
	}

	return nil
}

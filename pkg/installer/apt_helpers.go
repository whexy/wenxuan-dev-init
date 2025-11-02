package installer

import (
	"fmt"
	"os"
	"os/exec"
)

// ensurePrerequisites makes sure required tools are installed
func ensurePrerequisites() error {
	// Check if curl and gpg are available, install if not
	prereqs := []string{}

	if !IsCommandAvailable("curl") {
		prereqs = append(prereqs, "curl")
	}
	if !IsCommandAvailable("gpg") {
		prereqs = append(prereqs, "gnupg")
	}
	if !IsCommandAvailable("sudo") {
		prereqs = append(prereqs, "sudo")
	}

	if len(prereqs) > 0 {
		fmt.Printf("Installing prerequisites: %v\n", prereqs)
		updateCmd := exec.Command("sudo", "apt-get", "update", "-qq")
		updateCmd.Run() // Ignore errors

		args := append([]string{"apt-get", "install", "-y"}, prereqs...)
		cmd := exec.Command("sudo", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to install prerequisites: %w", err)
		}
	}

	return nil
}

// InstallGitHubCLI installs GitHub CLI with proper repository setup
func InstallGitHubCLI() error {
	// Ensure prerequisites
	if err := ensurePrerequisites(); err != nil {
		return err
	}

	fmt.Println("Setting up GitHub CLI repository...")

	// Add GitHub CLI repository
	cmd1 := exec.Command("bash", "-c", "curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | sudo dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg")
	cmd1.Stdout = os.Stdout
	cmd1.Stderr = os.Stderr
	if err := cmd1.Run(); err != nil {
		return fmt.Errorf("failed to add GitHub CLI keyring: %w", err)
	}

	cmd2 := exec.Command("bash", "-c", "echo \"deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main\" | sudo tee /etc/apt/sources.list.d/github-cli.list")
	cmd2.Stdout = os.Stdout
	cmd2.Stderr = os.Stderr
	if err := cmd2.Run(); err != nil {
		return fmt.Errorf("failed to add GitHub CLI repository: %w", err)
	}

	// Update and install
	updateCmd := exec.Command("sudo", "apt-get", "update")
	updateCmd.Stdout = os.Stdout
	updateCmd.Stderr = os.Stderr
	if err := updateCmd.Run(); err != nil {
		return fmt.Errorf("failed to update after adding repo: %w", err)
	}

	installCmd := exec.Command("sudo", "apt-get", "install", "-y", "gh")
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr
	if err := installCmd.Run(); err != nil {
		return fmt.Errorf("failed to install gh: %w", err)
	}

	fmt.Println("✓ GitHub CLI installed successfully")
	return nil
}

// Install1PasswordCLI installs 1Password CLI with proper repository setup
func Install1PasswordCLI() error {
	// Ensure prerequisites
	if err := ensurePrerequisites(); err != nil {
		return err
	}

	fmt.Println("Setting up 1Password CLI repository...")

	// Add 1Password repository
	cmd1 := exec.Command("bash", "-c", "curl -sS https://downloads.1password.com/linux/keys/1password.asc | sudo gpg --dearmor --output /usr/share/keyrings/1password-archive-keyring.gpg")
	cmd1.Stdout = os.Stdout
	cmd1.Stderr = os.Stderr
	if err := cmd1.Run(); err != nil {
		return fmt.Errorf("failed to add 1Password keyring: %w", err)
	}

	cmd2 := exec.Command("bash", "-c", "echo 'deb [arch=amd64 signed-by=/usr/share/keyrings/1password-archive-keyring.gpg] https://downloads.1password.com/linux/debian/amd64 stable main' | sudo tee /etc/apt/sources.list.d/1password.list")
	cmd2.Stdout = os.Stdout
	cmd2.Stderr = os.Stderr
	if err := cmd2.Run(); err != nil {
		return fmt.Errorf("failed to add 1Password repository: %w", err)
	}

	// Update and install
	updateCmd := exec.Command("sudo", "apt-get", "update")
	updateCmd.Stdout = os.Stdout
	updateCmd.Stderr = os.Stderr
	if err := updateCmd.Run(); err != nil {
		return fmt.Errorf("failed to update after adding repo: %w", err)
	}

	installCmd := exec.Command("sudo", "apt-get", "install", "-y", "1password-cli")
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr
	if err := installCmd.Run(); err != nil {
		return fmt.Errorf("failed to install 1password-cli: %w", err)
	}

	fmt.Println("✓ 1Password CLI installed successfully")
	return nil
}

// InstallChezmoi installs chezmoi using the official installer
func InstallChezmoi() error {
	// Ensure curl is available
	if !IsCommandAvailable("curl") {
		fmt.Println("Installing curl...")
		cmd := exec.Command("sudo", "apt-get", "install", "-y", "curl")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to install curl: %w", err)
		}
	}

	fmt.Println("Installing chezmoi to /usr/local/bin...")

	// Use the official installer with binary install to /usr/local/bin
	cmd := exec.Command("sh", "-c", "curl -fsLS get.chezmoi.io | sudo sh -s -- -b /usr/local/bin")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install chezmoi: %w", err)
	}

	fmt.Println("✓ Chezmoi installed successfully")
	return nil
}

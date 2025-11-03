package installer

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var (
	tailscaleAuthKeyReference = "op://Developer/tailscale auth key/credential"
)

// SetTailscaleAuthKeyReference sets the 1Password reference for Tailscale auth key
func SetTailscaleAuthKeyReference(ref string) {
	tailscaleAuthKeyReference = ref
}

// IsTailscaleSetup checks if Tailscale is configured and running
func IsTailscaleSetup() bool {
	if !IsCommandAvailable("tailscale") {
		return false
	}

	// Run 'tailscale status' to check if it's set up
	cmd := exec.Command("tailscale", "status")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}

	// If the output contains information about the tailscale network, it's set up
	// An unconfigured tailscale will return an error or empty output
	outputStr := strings.TrimSpace(string(output))
	return outputStr != "" && !strings.Contains(outputStr, "Logged out")
}

// InstallTailscale installs the Tailscale client
func InstallTailscale(pkgMgr PackageManager) error {
	if pkgMgr == nil {
		return fmt.Errorf("package manager is nil")
	}

	// Install tailscale package
	if err := pkgMgr.Install("tailscale"); err != nil {
		return fmt.Errorf("failed to install tailscale: %w", err)
	}

	return nil
}

// GetTailscaleAuthKeyFrom1Password retrieves the Tailscale auth key from 1Password
func GetTailscaleAuthKeyFrom1Password() (string, error) {
	fmt.Printf("Fetching Tailscale auth key from 1Password: %s\n", tailscaleAuthKeyReference)

	// Read the secret from 1Password
	cmd := exec.Command("op", "read", tailscaleAuthKeyReference)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to read from 1Password: %w", err)
	}

	authKey := strings.TrimSpace(out.String())
	if authKey == "" {
		return "", fmt.Errorf("empty auth key received from 1Password")
	}

	return authKey, nil
}

// SetupTailscale runs the Tailscale setup process using auth key from 1Password
func SetupTailscale() error {
	if !IsCommandAvailable("tailscale") {
		return fmt.Errorf("tailscale command not found")
	}

	// Get auth key from 1Password
	authKey, err := GetTailscaleAuthKeyFrom1Password()
	if err != nil {
		return fmt.Errorf("failed to get Tailscale auth key: %w", err)
	}

	fmt.Println("Connecting to Tailscale network...")

	// Run 'tailscale up' with the auth key
	cmd := exec.Command("tailscale", "up", "--authkey", authKey)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to setup tailscale: %w", err)
	}

	return nil
}
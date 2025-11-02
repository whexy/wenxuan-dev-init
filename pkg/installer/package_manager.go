package installer

import (
	"fmt"
	"os/exec"
	"runtime"
)

// PackageManager defines the interface for different package managers
type PackageManager interface {
	Name() string
	Install(packages ...string) error
	IsAvailable() bool
}

// DetectPackageManager detects the available package manager on the system
func DetectPackageManager() (PackageManager, error) {
	// Check for devbox first
	if IsCommandAvailable("devbox") {
		return NewDevboxManager(), nil
	}

	// Detect based on OS
	switch runtime.GOOS {
	case "darwin":
		if IsCommandAvailable("brew") {
			return NewBrewManager(), nil
		}
	case "linux":
		if IsCommandAvailable("apt-get") {
			return NewAptManager(), nil
		}
		if IsCommandAvailable("pacman") {
			return NewPacmanManager(), nil
		}
		if IsCommandAvailable("dnf") {
			return NewDnfManager(), nil
		}
		if IsCommandAvailable("yum") {
			return NewYumManager(), nil
		}
	}

	return nil, fmt.Errorf("no supported package manager found")
}

// IsCommandAvailable checks if a command is available in PATH
func IsCommandAvailable(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

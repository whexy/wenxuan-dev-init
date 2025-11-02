package installer

import (
	"os"
	"strings"
)

// IsRunningInContainer detects if the program is running inside a container
func IsRunningInContainer() bool {
	// Check for .dockerenv file
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}

	// Check cgroup to see if we're in a container
	if data, err := os.ReadFile("/proc/1/cgroup"); err == nil {
		content := string(data)
		if strings.Contains(content, "docker") ||
			strings.Contains(content, "lxc") ||
			strings.Contains(content, "containerd") {
			return true
		}
	}

	// Check for container environment variables
	containerEnvVars := []string{
		"DOCKER_CONTAINER",
		"KUBERNETES_SERVICE_HOST",
		"container",
	}

	for _, envVar := range containerEnvVars {
		if os.Getenv(envVar) != "" {
			return true
		}
	}

	return false
}

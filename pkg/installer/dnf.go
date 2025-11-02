package installer

import (
	"fmt"
	"os/exec"
)

type DnfManager struct{}

func NewDnfManager() *DnfManager {
	return &DnfManager{}
}

func (d *DnfManager) Name() string {
	return "dnf"
}

func (d *DnfManager) IsAvailable() bool {
	return IsCommandAvailable("dnf")
}

func (d *DnfManager) Install(packages ...string) error {
	args := append([]string{"install", "-y"}, packages...)
	cmd := exec.Command("sudo", append([]string{"dnf"}, args...)...)
	cmd.Stdout = nil
	cmd.Stderr = nil

	fmt.Printf("Running: sudo dnf install -y %v\n", packages)
	return cmd.Run()
}

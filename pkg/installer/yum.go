package installer

import (
	"fmt"
	"os/exec"
)

type YumManager struct{}

func NewYumManager() *YumManager {
	return &YumManager{}
}

func (y *YumManager) Name() string {
	return "yum"
}

func (y *YumManager) IsAvailable() bool {
	return IsCommandAvailable("yum")
}

func (y *YumManager) Install(packages ...string) error {
	args := append([]string{"install", "-y"}, packages...)
	cmd := exec.Command("sudo", append([]string{"yum"}, args...)...)
	cmd.Stdout = nil
	cmd.Stderr = nil

	fmt.Printf("Running: sudo yum install -y %v\n", packages)
	return cmd.Run()
}

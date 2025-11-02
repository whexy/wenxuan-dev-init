package installer

import (
	"fmt"
	"os/exec"
)

type PacmanManager struct{}

func NewPacmanManager() *PacmanManager {
	return &PacmanManager{}
}

func (p *PacmanManager) Name() string {
	return "pacman"
}

func (p *PacmanManager) IsAvailable() bool {
	return IsCommandAvailable("pacman")
}

func (p *PacmanManager) Install(packages ...string) error {
	args := append([]string{"-S", "--noconfirm"}, packages...)
	cmd := exec.Command("sudo", append([]string{"pacman"}, args...)...)
	cmd.Stdout = nil
	cmd.Stderr = nil

	fmt.Printf("Running: sudo pacman -S --noconfirm %v\n", packages)
	return cmd.Run()
}

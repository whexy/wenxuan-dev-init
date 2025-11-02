package installer

import (
	"fmt"
	"os/exec"
)

type BrewManager struct{}

func NewBrewManager() *BrewManager {
	return &BrewManager{}
}

func (b *BrewManager) Name() string {
	return "brew"
}

func (b *BrewManager) IsAvailable() bool {
	return IsCommandAvailable("brew")
}

func (b *BrewManager) Install(packages ...string) error {
	args := append([]string{"install"}, packages...)
	cmd := exec.Command("brew", args...)
	cmd.Stdout = nil
	cmd.Stderr = nil

	fmt.Printf("Running: brew %v\n", args)
	return cmd.Run()
}

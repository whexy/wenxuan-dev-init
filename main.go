package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/whexy/wenxuan-dev-init/pkg/executor"
	"github.com/whexy/wenxuan-dev-init/pkg/installer"
	"github.com/whexy/wenxuan-dev-init/pkg/logger"
	"github.com/whexy/wenxuan-dev-init/pkg/tui"
)

var (
	githubTokenRef       = flag.String("github-token", "op://Developer/GitHub Personal Access Token/token", "1Password reference for GitHub token")
	tailscaleAuthKeyRef  = flag.String("tailscale-authkey", "op://Developer/tailscale auth key/credential", "1Password reference for Tailscale auth key")
	useServiceAccount    = flag.Bool("use-service-account", false, "Use 1Password service account token (requires OP_SERVICE_ACCOUNT_TOKEN)")
)

func main() {
	flag.Parse()

	// Set the global references
	installer.SetGitHubTokenReference(*githubTokenRef)
	installer.SetTailscaleAuthKeyReference(*tailscaleAuthKeyRef)
	installer.SetUseServiceAccount(*useServiceAccount)

	if err := run(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func run() error {
	// Run the interactive TUI
	options, err := runTUI()
	if err != nil {
		return fmt.Errorf("TUI error: %w", err)
	}

	// User cancelled
	if options == nil {
		fmt.Println("\nSetup cancelled.")
		return nil
	}

	// Execute the workflow
	exec := executor.New(options)
	if err := exec.Execute(); err != nil {
		return fmt.Errorf("execution error: %w", err)
	}

	return nil
}

func runTUI() (map[string]bool, error) {
	model := tui.NewModel()
	p := tea.NewProgram(model, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		return nil, err
	}

	m, ok := finalModel.(tui.Model)
	if !ok {
		return nil, fmt.Errorf("unexpected model type")
	}

	if !m.IsConfirmed() {
		return nil, nil
	}

	return m.GetSelectedOptions(), nil
}

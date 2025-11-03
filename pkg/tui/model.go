package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/whexy/wenxuan-dev-init/pkg/installer"
)

type Dependency struct {
	Name      string
	Command   string
	Available bool
	Icon      string
}

type ConfigOption struct {
	Label       string
	Description string
	Enabled     bool
	Key         string
}

type Model struct {
	dependencies []Dependency
	options      []ConfigOption
	cursor       int
	confirmed    bool
	width        int
	height       int
}

type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Space  key.Binding
	Enter  key.Binding
	Quit   key.Binding
	Help   key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Space: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("space", "toggle"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "confirm"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
}

func NewModel() Model {
	// Check dependencies
	deps := []Dependency{
		{Name: "Git", Command: "git", Available: installer.IsCommandAvailable("git"), Icon: "🔧"},
		{Name: "GitHub CLI", Command: "gh", Available: installer.IsCommandAvailable("gh"), Icon: "🐙"},
		{Name: "1Password CLI", Command: "op", Available: installer.IsCommandAvailable("op"), Icon: "🔐"},
		{Name: "Chezmoi", Command: "chezmoi", Available: installer.IsCommandAvailable("chezmoi"), Icon: "🏠"},
		{Name: "Devbox", Command: "devbox", Available: installer.IsCommandAvailable("devbox"), Icon: "📦"},
		{Name: "Tailscale", Command: "tailscale", Available: installer.IsCommandAvailable("tailscale"), Icon: "🔗"},
	}

	// Detect package manager
	pkgMgrAvailable := false
	pkgMgrName := "none"
	if pkgMgr, err := installer.DetectPackageManager(); err == nil {
		pkgMgrAvailable = true
		pkgMgrName = pkgMgr.Name()
	}

	deps = append(deps, Dependency{
		Name:      fmt.Sprintf("Package Manager (%s)", pkgMgrName),
		Command:   pkgMgrName,
		Available: pkgMgrAvailable,
		Icon:      "📦",
	})

	// Configuration options
	inContainer := installer.IsRunningInContainer()
	devboxDescription := "Install devbox package manager (recommended for Linux)"
	devboxEnabled := !installer.IsCommandAvailable("devbox")

	if inContainer {
		devboxDescription = "⚠️  NOT recommended in containers (requires Nix daemon)"
		devboxEnabled = false // Disable by default in containers
	}

	options := []ConfigOption{
		{
			Label:       "Install Devbox",
			Description: devboxDescription,
			Enabled:     devboxEnabled,
			Key:         "install_devbox",
		},
		{
			Label:       "Install Git",
			Description: "Install git version control system",
			Enabled:     !installer.IsCommandAvailable("git"),
			Key:         "install_git",
		},
		{
			Label:       "Install GitHub CLI",
			Description: "Install gh command-line tool",
			Enabled:     !installer.IsCommandAvailable("gh"),
			Key:         "install_gh",
		},
		{
			Label:       "Install 1Password CLI",
			Description: "Install 1Password command-line tool",
			Enabled:     !installer.IsCommandAvailable("op"),
			Key:         "install_1password",
		},
		{
			Label:       "Install Chezmoi",
			Description: "Install chezmoi dotfile manager",
			Enabled:     !installer.IsCommandAvailable("chezmoi"),
			Key:         "install_chezmoi",
		},
		{
			Label:       "Install Tailscale",
			Description: "Install Tailscale VPN client",
			Enabled:     !installer.IsCommandAvailable("tailscale"),
			Key:         "install_tailscale",
		},
		{
			Label:       "Login to 1Password",
			Description: "Authenticate with 1Password",
			Enabled:     true,
			Key:         "login_1password",
		},
		{
			Label:       "Setup GitHub Authentication",
			Description: "Configure GitHub CLI with 1Password token",
			Enabled:     true,
			Key:         "setup_github",
		},
		{
			Label:       "Initialize Chezmoi",
			Description: "Run chezmoi init --apply whexy",
			Enabled:     true,
			Key:         "init_chezmoi",
		},
		{
			Label:       "Setup Tailscale",
			Description: "Configure and connect to Tailscale network",
			Enabled:     installer.IsCommandAvailable("tailscale") && !installer.IsTailscaleSetup(),
			Key:         "setup_tailscale",
		},
	}

	return Model{
		dependencies: deps,
		options:      options,
		cursor:       0,
		confirmed:    false,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, keys.Up):
			if m.cursor > 0 {
				m.cursor--
			}

		case key.Matches(msg, keys.Down):
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}

		case key.Matches(msg, keys.Space):
			if m.cursor < len(m.options) {
				m.options[m.cursor].Enabled = !m.options[m.cursor].Enabled
			}

		case key.Matches(msg, keys.Enter):
			m.confirmed = true
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Model) View() string {
	if m.confirmed {
		return ""
	}

	var s strings.Builder

	// Styles
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		Background(lipgloss.Color("#1a1a1a")).
		Padding(0, 1).
		MarginBottom(1)

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1)

	availableStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#04B575"))

	unavailableStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF0000"))

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Padding(1, 2).
		MarginBottom(1)

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4")).
		Bold(true)

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Italic(true)

	// Title
	s.WriteString(titleStyle.Render("🚀 Wenxuan Dev Init - Interactive Setup"))
	s.WriteString("\n\n")

	// Top half - Dependencies status
	s.WriteString(headerStyle.Render(" 📊 System Dependencies Status "))
	s.WriteString("\n\n")

	var depLines []string
	for _, dep := range m.dependencies {
		status := "✓"
		style := availableStyle
		statusText := "Available"
		if !dep.Available {
			status = "✗"
			style = unavailableStyle
			statusText = "Missing"
		}
		line := fmt.Sprintf("  %s %s %s %s",
			dep.Icon,
			dep.Name,
			style.Render(status),
			style.Render(statusText),
		)
		depLines = append(depLines, line)
	}

	s.WriteString(boxStyle.Render(strings.Join(depLines, "\n")))
	s.WriteString("\n")

	// Bottom half - Configuration options
	s.WriteString(headerStyle.Render(" ⚙️  Configuration Options "))
	s.WriteString("\n\n")

	for i, option := range m.options {
		cursor := "  "
		checkbox := "[ ]"
		if option.Enabled {
			checkbox = "[✓]"
		}

		line := fmt.Sprintf("%s %s %s", cursor, checkbox, option.Label)

		if i == m.cursor {
			line = selectedStyle.Render(fmt.Sprintf("▶ %s %s", checkbox, option.Label))
			s.WriteString(line)
			s.WriteString("\n")
			s.WriteString("  " + descStyle.Render("  "+option.Description))
			s.WriteString("\n")
		} else {
			s.WriteString(line)
			s.WriteString("\n")
		}
	}

	// Help text
	s.WriteString("\n")
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Padding(1, 0)

	help := "↑/↓: navigate • space: toggle • enter: confirm • q: quit"
	s.WriteString(helpStyle.Render(help))

	return s.String()
}

func (m Model) GetSelectedOptions() map[string]bool {
	result := make(map[string]bool)
	for _, option := range m.options {
		result[option.Key] = option.Enabled
	}
	return result
}

func (m Model) IsConfirmed() bool {
	return m.confirmed
}
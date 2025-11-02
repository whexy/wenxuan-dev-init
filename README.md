# wenxuan-dev-init

A TUI tool for bootstrapping development environments.

Automates the installation of Git, GitHub CLI, 1Password CLI, chezmoi, and related authentication setup.

## What it does

- Detects your system's package manager (apt/brew/pacman/dnf/yum)
- Shows dependency status in an interactive TUI
- Installs missing tools: git, gh, 1password-cli, chezmoi
- Handles GitHub and 1Password authentication
- Sets up dotfiles via chezmoi

## Features

Interactive TUI built with [Bubble Tea](https://github.com/charmbracelet/bubbletea):

- Real-time dependency status
- Checkbox-based configuration
- Vim-style navigation

Supported package managers:

- macOS: Homebrew
- Debian/Ubuntu: apt
- Arch Linux: pacman
- Fedora: dnf
- CentOS/RHEL: yum
- Devbox (optional)

## Usage

```
🚀 Wenxuan Dev Init - Interactive Setup
═══════════════════════════════════════

 📊 System Dependencies Status
╭─────────────────────────────────────╮
│  🔧 Git              ✓ Available   │
│  🐙 GitHub CLI       ✗ Missing     │
│  🔐 1Password CLI    ✗ Missing     │
│  🏠 Chezmoi          ✗ Missing     │
│  📦 Devbox           ✗ Missing     │
│  📦 Package Manager  ✓ Available   │
╰─────────────────────────────────────╯

 ⚙️  Configuration Options

  [ ] Install Devbox
▶ [✓] Install GitHub CLI
      Install gh command-line tool
  [✓] Install 1Password CLI
  [✓] Install Chezmoi
  [✓] Login to 1Password
  [✓] Setup GitHub Authentication
  [✓] Initialize Chezmoi

↑/↓: navigate • space: toggle • enter: confirm • q: quit
```

## Installation

Download the binary:

```bash
# Linux (x86_64)
curl -LO https://github.com/whexy/wenxuan-dev-init/releases/latest/download/wenxuan-dev-init-linux-amd64
chmod +x wenxuan-dev-init-linux-amd64
sudo mv wenxuan-dev-init-linux-amd64 /usr/local/bin/wenxuan-dev-init

# macOS (Intel)
curl -LO https://github.com/whexy/wenxuan-dev-init/releases/latest/download/wenxuan-dev-init-darwin-amd64
chmod +x wenxuan-dev-init-darwin-amd64
sudo mv wenxuan-dev-init-darwin-amd64 /usr/local/bin/wenxuan-dev-init

# macOS (Apple Silicon)
curl -LO https://github.com/whexy/wenxuan-dev-init/releases/latest/download/wenxuan-dev-init-darwin-arm64
chmod +x wenxuan-dev-init-darwin-arm64
sudo mv wenxuan-dev-init-darwin-arm64 /usr/local/bin/wenxuan-dev-init
```

Run it:

```bash
wenxuan-dev-init
```

Or build from source:

```bash
git clone https://github.com/whexy/wenxuan-dev-init.git
cd wenxuan-dev-init
make build
```

## Testing

Test in Docker containers:

```bash
cd deployment
docker-compose run ubuntu-test  # or debian-test, fedora-test
```

See [TESTING.md](TESTING.md) for details.

## Structure

```
pkg/
├── executor/    # workflow orchestration
├── installer/   # package manager implementations
├── tui/         # bubble tea interface
└── logger/      # output formatting
```

Written in Go. Single static binary. MIT licensed.

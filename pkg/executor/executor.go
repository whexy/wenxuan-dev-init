package executor

import (
	"fmt"

	"github.com/whexy/wenxuan-dev-init/pkg/installer"
	"github.com/whexy/wenxuan-dev-init/pkg/logger"
	"github.com/whexy/wenxuan-dev-init/pkg/ui"
)

// Config holds the execution configuration
type Config struct {
	InstallDevbox    bool
	InstallGit       bool
	InstallGH        bool
	Install1Password bool
	InstallChezmoi   bool
	Login1Password   bool
	SetupGitHub      bool
	InitChezmoi      bool
}

// Executor handles the execution workflow
type Executor struct {
	config Config
	pkgMgr installer.PackageManager
}

// New creates a new Executor with the given configuration
func New(options map[string]bool) *Executor {
	return &Executor{
		config: Config{
			InstallDevbox:    options["install_devbox"],
			InstallGit:       options["install_git"],
			InstallGH:        options["install_gh"],
			Install1Password: options["install_1password"],
			InstallChezmoi:   options["install_chezmoi"],
			Login1Password:   options["login_1password"],
			SetupGitHub:      options["setup_github"],
			InitChezmoi:      options["init_chezmoi"],
		},
	}
}

// Execute runs the complete setup workflow
func (e *Executor) Execute() error {
	logger.Step("🚀", "Starting setup process...")
	logger.Println("")

	// Step 1: Setup package manager
	if err := e.setupPackageManager(); err != nil {
		return fmt.Errorf("failed to setup package manager: %w", err)
	}

	// Step 2: Install packages
	if err := e.installPackages(); err != nil {
		return fmt.Errorf("failed to install packages: %w", err)
	}

	// Step 3: Authenticate 1Password
	if e.config.Login1Password {
		// Check if 1Password CLI is available
		if !installer.IsCommandAvailable("op") {
			logger.Warning("1Password CLI not found, skipping authentication.")
			logger.Println("Install 1Password CLI manually and run 'op signin' if needed.")
		} else {
			if err := e.authenticate1Password(); err != nil {
				logger.Error(fmt.Sprintf("1Password authentication failed: %v", err))
				logger.Warning("Continuing without 1Password authentication.")
			}
		}
	}

	// Step 4: Setup GitHub
	if e.config.SetupGitHub {
		// Check if required tools are available
		if !installer.IsCommandAvailable("gh") {
			logger.Warning("GitHub CLI not found, skipping GitHub setup.")
			logger.Println("Install GitHub CLI manually if needed.")
		} else if !installer.IsCommandAvailable("op") {
			logger.Warning("1Password CLI not found, skipping GitHub setup.")
			logger.Println("You'll need to authenticate GitHub manually with 'gh auth login'.")
		} else {
			if err := e.setupGitHub(); err != nil {
				logger.Error(fmt.Sprintf("GitHub setup failed: %v", err))
				logger.Warning("You can authenticate manually with 'gh auth login'.")
			}
		}
	}

	// Step 5: Initialize Chezmoi
	if e.config.InitChezmoi {
		// Check if chezmoi is available
		if !installer.IsCommandAvailable("chezmoi") {
			logger.Warning("Chezmoi not found, skipping initialization.")
			logger.Println("Install chezmoi manually if needed.")
		} else {
			if err := e.initializeChezmoi(); err != nil {
				logger.Error(fmt.Sprintf("Chezmoi initialization failed: %v", err))
				logger.Warning("You can initialize manually with 'chezmoi init --apply whexy'.")
			}
		}
	}

	logger.Println("")
	logger.Success("Setup complete! Your development environment is ready.")
	return nil
}

func (e *Executor) setupPackageManager() error {
	if e.config.InstallDevbox {
		// Check if devbox is already available
		if !installer.IsCommandAvailable("devbox") {
			logger.Step("📦", "Installing devbox...")
			if err := installer.InstallDevbox(); err != nil {
				logger.Error(fmt.Sprintf("Devbox installation failed: %v", err))
				logger.Warning("Devbox installation encountered errors.")

				// Offer fallback
				if ui.AskYesNo("Would you like to use the system package manager instead?") {
					logger.Info("Falling back to system package manager...")
					var detectErr error
					e.pkgMgr, detectErr = installer.DetectPackageManager()
					if detectErr != nil {
						return fmt.Errorf("failed to detect system package manager: %w", detectErr)
					}
				} else {
					return fmt.Errorf("setup cancelled by user")
				}
			} else {
				logger.Success("Devbox installed successfully!")
				logger.Println("")
				logger.Info("⚠️  Devbox requires environment initialization.")
				logger.Println("")
				logger.Info("Please run the following commands to activate devbox:")
				logger.Println("   eval \"$(devbox global shellenv --init-hook)\"")
				logger.Println("")
				logger.Info("Then rerun this program to continue the setup.")
				logger.Println("")
				return fmt.Errorf("devbox installed - please reload your shell and rerun")
			}
		} else {
			logger.Info("Devbox is already installed")
			e.pkgMgr = installer.NewDevboxManager()
		}
	} else {
		var err error
		e.pkgMgr, err = installer.DetectPackageManager()
		if err != nil {
			return err
		}
	}

	logger.Info(fmt.Sprintf("Using package manager: %s", e.pkgMgr.Name()))
	logger.Println("")

	return nil
}

func (e *Executor) installPackages() error {
	packagesToInstall := e.getPackagesToInstall()

	if len(packagesToInstall) == 0 {
		return nil
	}

	logger.Step("📦", fmt.Sprintf("Installing packages: %v", packagesToInstall))
	if err := e.pkgMgr.Install(packagesToInstall...); err != nil {
		logger.Println("") // Add spacing after error output
		logger.Error(fmt.Sprintf("Package installation failed: %v", err))

		// If using devbox and it failed, offer to fallback
		if e.pkgMgr.Name() == "devbox" {
			logger.Warning("Devbox package installation failed (this is common in containers).")

			if ui.AskYesNo("Would you like to try with the system package manager instead?") {
				logger.Info("Switching to system package manager...")

				// Detect and switch to system package manager
				systemPkgMgr, detectErr := installer.DetectPackageManager()
				if detectErr != nil {
					logger.Error(fmt.Sprintf("Failed to detect system package manager: %v", detectErr))
					logger.Warning("You can install packages manually later.")
					return nil
				}

				e.pkgMgr = systemPkgMgr
				logger.Info(fmt.Sprintf("Using package manager: %s", e.pkgMgr.Name()))

				// Retry installation with system package manager
				logger.Step("📦", fmt.Sprintf("Retrying installation: %v", packagesToInstall))
				if retryErr := e.pkgMgr.Install(packagesToInstall...); retryErr != nil {
					logger.Println("") // Add spacing
					logger.Error(fmt.Sprintf("Installation failed again: %v", retryErr))
					logger.Warning("You can install packages manually later.")
					return nil
				}

				logger.Success("Packages installed successfully with system package manager!")
				return nil
			}
		}

		logger.Warning("Some packages may not have been installed.")
		logger.Println("Please check the error messages above for details.")
		return nil
	}
	logger.Success("Packages installed successfully")

	// Show devbox shell initialization hint
	if e.pkgMgr.Name() == "devbox" {
		logger.Println("")
		logger.Warning("Don't forget to initialize devbox shell:")
		logger.Println("   eval \"$(devbox global shellenv --init-hook)\"")
	}

	return nil
}

func (e *Executor) getPackagesToInstall() []string {
	var packages []string

	if e.config.InstallGit {
		packages = append(packages, "git")
	}

	if e.config.InstallGH {
		packages = append(packages, "gh")
	}

	if e.config.Install1Password {
		if e.pkgMgr.Name() == "devbox" {
			packages = append(packages, "_1password-cli")
		} else {
			packages = append(packages, "1password-cli")
		}
	}

	if e.config.InstallChezmoi {
		packages = append(packages, "chezmoi")
	}

	return packages
}

func (e *Executor) authenticate1Password() error {
	logger.Println("")
	logger.Step("🔐", "Logging in to 1Password...")
	if err := installer.Login1Password(); err != nil {
		return err
	}
	logger.Success("Logged in to 1Password")
	return nil
}

func (e *Executor) setupGitHub() error {
	logger.Println("")
	logger.Step("🐙", "Setting up GitHub authentication...")

	token, err := installer.GetGitHubTokenFrom1Password()
	if err != nil {
		return err
	}

	if err := installer.AuthenticateGitHub(token); err != nil {
		return err
	}

	logger.Success("GitHub authentication successful")
	return nil
}

func (e *Executor) initializeChezmoi() error {
	logger.Println("")
	logger.Step("🏠", "Initializing chezmoi...")

	if err := installer.InitChezmoi("whexy"); err != nil {
		return err
	}

	logger.Success("Chezmoi initialized successfully")
	return nil
}

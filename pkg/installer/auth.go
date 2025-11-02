package installer

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var (
	githubTokenReference = "op://Developer/GitHub Personal Access Token/token"
	useServiceAccount    = false
)

// SetGitHubTokenReference sets the 1Password reference for GitHub token
func SetGitHubTokenReference(ref string) {
	githubTokenReference = ref
}

// SetUseServiceAccount sets whether to use service account authentication
func SetUseServiceAccount(use bool) {
	useServiceAccount = use
}

// Login1Password prompts the user to log in to 1Password
func Login1Password() error {
	// If using service account, check for token
	if useServiceAccount {
		return ensureServiceAccountToken()
	}

	// Regular interactive signin
	fmt.Println("Logging in to 1Password...")
	fmt.Println("Please follow the prompts to authenticate.")

	// Use --force flag to bypass the eval warning
	cmd := exec.Command("op", "signin", "--force")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// ensureServiceAccountToken checks for OP_SERVICE_ACCOUNT_TOKEN and prompts if not set
func ensureServiceAccountToken() error {
	token := os.Getenv("OP_SERVICE_ACCOUNT_TOKEN")

	if token != "" {
		fmt.Println("✓ Using 1Password service account token from environment")
		return nil
	}

	// Token not set, prompt user
	fmt.Println("1Password service account mode enabled, but OP_SERVICE_ACCOUNT_TOKEN is not set.")
	fmt.Println("")
	fmt.Println("Please enter your 1Password service account token:")
	fmt.Print("> ")

	var inputToken string
	fmt.Scanln(&inputToken)

	if inputToken == "" {
		return fmt.Errorf("no service account token provided")
	}

	// Set the environment variable for current process and children
	os.Setenv("OP_SERVICE_ACCOUNT_TOKEN", inputToken)
	fmt.Println("✓ Service account token set successfully")

	return nil
}

// GetGitHubTokenFrom1Password retrieves the GitHub token from 1Password
func GetGitHubTokenFrom1Password() (string, error) {
	fmt.Printf("Fetching GitHub token from 1Password: %s\n", githubTokenReference)

	// Read the secret from 1Password
	cmd := exec.Command("op", "read", githubTokenReference)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to read from 1Password: %w", err)
	}

	token := strings.TrimSpace(out.String())
	if token == "" {
		return "", fmt.Errorf("empty token received from 1Password")
	}

	return token, nil
}

// AuthenticateGitHub authenticates GitHub CLI with a token
func AuthenticateGitHub(token string) error {
	fmt.Println("Authenticating GitHub CLI...")

	// Use gh auth login with token via stdin
	cmd := exec.Command("gh", "auth", "login", "--with-token")
	cmd.Stdin = strings.NewReader(token)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to authenticate GitHub: %w", err)
	}

	// Configure git to use gh as credential helper
	gitCmd := exec.Command("gh", "auth", "setup-git")
	gitCmd.Stdout = os.Stdout
	gitCmd.Stderr = os.Stderr

	if err := gitCmd.Run(); err != nil {
		return fmt.Errorf("failed to setup git authentication: %w", err)
	}

	fmt.Println("GitHub authentication successful!")
	return nil
}

// InitChezmoi initializes chezmoi with the specified GitHub username
func InitChezmoi(username string) error {
	// If using service account, setup chezmoi config first
	if useServiceAccount {
		if err := setupChezmoiConfig(); err != nil {
			return fmt.Errorf("failed to setup chezmoi config: %w", err)
		}
	}

	fmt.Printf("Initializing chezmoi with GitHub user: %s\n", username)

	cmd := exec.Command("chezmoi", "init", "--apply", username)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to initialize chezmoi: %w", err)
	}

	fmt.Println("Chezmoi initialized successfully!")
	return nil
}

// setupChezmoiConfig creates chezmoi config with 1Password service mode
func setupChezmoiConfig() error {
	home := os.Getenv("HOME")
	configDir := home + "/.config/chezmoi"
	configFile := configDir + "/chezmoi.toml"

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Create/overwrite config file with 1Password service mode
	config := `[onepassword]
mode = "service"
`

	if err := os.WriteFile(configFile, []byte(config), 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	fmt.Printf("✓ Created chezmoi config with 1Password service mode\n")
	return nil
}

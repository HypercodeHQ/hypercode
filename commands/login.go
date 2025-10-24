package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/urfave/cli/v3"
)

// LoginCommand creates the login command for CLI authentication
func LoginCommand() *cli.Command {
	return NewCommandWithFlags(
		"login",
		"Authenticate with Hypercommit",
		`Authenticate the CLI with your Hypercommit account using device flow.

This command will open a browser where you can confirm the authentication code.
Once confirmed, your credentials will be stored locally for Git operations.`,
		[]cli.Flag{
			&cli.StringFlag{
				Name:    "host",
				Aliases: []string{"h"},
				Usage:   "Hypercommit server host",
				Value:   "localhost:3000",
			},
		},
		runLogin,
	)
}

type deviceCodeResponse struct {
	SessionID       string `json:"session_id"`
	UserCode        string `json:"user_code"`
	VerificationURL string `json:"verification_url"`
	ExpiresAt       int64  `json:"expires_at"`
	Interval        int    `json:"interval"`
}

type pollResponse struct {
	Status      string `json:"status"`
	AccessToken string `json:"access_token,omitempty"`
	Username    string `json:"username,omitempty"`
	Error       string `json:"error,omitempty"`
}

func runLogin(ctx context.Context, cmd *cli.Command) error {
	host := cmd.String("host")
	protocol := "http"
	if host != "localhost:3000" && host != "127.0.0.1:3000" {
		protocol = "https"
	}
	baseURL := fmt.Sprintf("%s://%s", protocol, host)

	fmt.Println("🔐 Authenticating with Hypercommit...")
	fmt.Println()

	// Step 1: Initiate device auth flow
	fmt.Print("⏳ Requesting authentication code... ")
	deviceResp, err := initiateDeviceAuth(baseURL)
	if err != nil {
		fmt.Println("✗")
		return fmt.Errorf("failed to initiate authentication: %w", err)
	}
	fmt.Println("✓")

	// Step 2: Display code and URL
	fmt.Println("✓ Authentication code generated")
	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("  Your one-time code: \033[1;36m%s\033[0m\n", deviceResp.UserCode)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	// Step 3: Ask if user wants to open browser
	var openBrowser bool
	openForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Open browser to authenticate?").
				Description(fmt.Sprintf("URL: %s", deviceResp.VerificationURL)).
				Affirmative("Yes").
				Negative("No").
				Value(&openBrowser),
		),
	)

	if err := openForm.Run(); err != nil {
		return fmt.Errorf("cancelled: %w", err)
	}

	if openBrowser {
		if err := openURL(deviceResp.VerificationURL); err != nil {
			fmt.Printf("⚠️  Could not open browser automatically: %v\n", err)
			fmt.Printf("Please visit: %s\n\n", deviceResp.VerificationURL)
		} else {
			fmt.Println("✓ Opened browser")
			fmt.Println()
		}
	} else {
		fmt.Printf("Please visit: %s\n\n", deviceResp.VerificationURL)
	}

	// Step 4: Poll for confirmation
	fmt.Println("⏳ Waiting for confirmation...")
	fmt.Println()

	pollResult, err := pollForConfirmation(baseURL, deviceResp.SessionID, deviceResp.Interval, deviceResp.ExpiresAt)
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Step 5: Store credentials
	if err := storeCredentials(host, pollResult.Username, pollResult.AccessToken); err != nil {
		return fmt.Errorf("failed to store credentials: %w", err)
	}

	fmt.Println()
	fmt.Println("✨ Successfully authenticated!")
	fmt.Printf("   Logged in as: \033[1;32m@%s\033[0m\n", pollResult.Username)
	fmt.Println()
	fmt.Println("You can now use Git commands with your Hypercommit repositories.")

	return nil
}

func initiateDeviceAuth(baseURL string) (deviceCodeResponse, error) {
	var resp deviceCodeResponse

	httpResp, err := http.Post(baseURL+"/api/auth/device/code", "application/json", nil)
	if err != nil {
		return resp, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResp.Body)
		return resp, fmt.Errorf("server returned %d: %s", httpResp.StatusCode, string(body))
	}

	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return resp, err
	}

	return resp, nil
}

func pollForConfirmation(baseURL, sessionID string, interval int, expiresAt int64) (*pollResponse, error) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	timeout := time.Until(time.Unix(expiresAt, 0))
	timeoutTimer := time.NewTimer(timeout)
	defer timeoutTimer.Stop()

	pollURL := fmt.Sprintf("%s/api/auth/device/poll?session_id=%s", baseURL, sessionID)

	for {
		select {
		case <-timeoutTimer.C:
			return nil, fmt.Errorf("authentication timed out")

		case <-ticker.C:
			httpResp, err := http.Get(pollURL)
			if err != nil {
				continue // Retry on network errors
			}

			body, err := io.ReadAll(httpResp.Body)
			httpResp.Body.Close()
			if err != nil {
				continue
			}

			var pollResp pollResponse
			if err := json.Unmarshal(body, &pollResp); err != nil {
				continue
			}

			switch pollResp.Status {
			case "confirmed":
				fmt.Println("✓ Authentication confirmed!")
				return &pollResp, nil

			case "expired":
				return nil, fmt.Errorf("authentication code expired: %s", pollResp.Error)

			case "pending":
				// Continue polling
				continue

			default:
				// Unknown status, continue polling
				continue
			}
		}
	}
}

func storeCredentials(host, username, accessToken string) error {
	// Get config directory
	configDir, err := getConfigDir()
	if err != nil {
		return err
	}

	// Ensure directory exists
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return err
	}

	// Store in Git credential helper
	protocol := "http"
	if host != "localhost:3000" && host != "127.0.0.1:3000" {
		protocol = "https"
	}

	// Try to approve credentials with git credential helper
	cmd := exec.Command("git", "credential", "approve")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	// Write credentials in git credential format
	credInput := fmt.Sprintf("protocol=%s\nhost=%s\nusername=%s\npassword=%s\n\n", protocol, host, username, accessToken)
	if _, err := stdin.Write([]byte(credInput)); err != nil {
		stdin.Close()
		cmd.Wait()
		return err
	}
	stdin.Close()

	if err := cmd.Wait(); err != nil {
		// If git credential fails, store in our own config file as fallback
		return storeFallbackCredentials(configDir, host, username, accessToken)
	}

	return nil
}

func storeFallbackCredentials(configDir, host, username, accessToken string) error {
	configPath := filepath.Join(configDir, "credentials")

	config := map[string]interface{}{
		"host":         host,
		"username":     username,
		"access_token": accessToken,
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0600)
}

func getConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".hypercommit"), nil
}

func openURL(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		// Try xdg-open first, fallback to common browsers
		if _, err := exec.LookPath("xdg-open"); err == nil {
			cmd = exec.Command("xdg-open", url)
		} else if _, err := exec.LookPath("firefox"); err == nil {
			cmd = exec.Command("firefox", url)
		} else if _, err := exec.LookPath("google-chrome"); err == nil {
			cmd = exec.Command("google-chrome", url)
		} else {
			return fmt.Errorf("no suitable browser found")
		}
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	return cmd.Start()
}

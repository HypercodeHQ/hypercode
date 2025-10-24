package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hypercommithq/hypercommit/database/repositories"
	"github.com/hypercommithq/hypercommit/httperror"
	"github.com/hypercommithq/hypercommit/services"
)

type GitHubAuthController interface {
	Login(w http.ResponseWriter, r *http.Request) error
	Callback(w http.ResponseWriter, r *http.Request) error
}

type githubAuthController struct {
	users       repositories.UsersRepository
	authService services.AuthService
	githubOAuth services.GitHubOAuthService
}

func NewGitHubAuthController(
	users repositories.UsersRepository,
	authService services.AuthService,
	githubOAuth services.GitHubOAuthService,
) GitHubAuthController {
	return &githubAuthController{
		users:       users,
		authService: authService,
		githubOAuth: githubOAuth,
	}
}

func (c *githubAuthController) Login(w http.ResponseWriter, r *http.Request) error {
	// Generate random state for CSRF protection
	state := generateRandomState()

	// Store state in cookie for verification in callback
	http.SetCookie(w, &http.Cookie{
		Name:     "github_oauth_state",
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   600, // 10 minutes
	})

	authURL := c.githubOAuth.GetAuthURL(state)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
	return nil
}

func (c *githubAuthController) Callback(w http.ResponseWriter, r *http.Request) error {
	// Verify state parameter
	stateCookie, err := r.Cookie("github_oauth_state")
	if err != nil {
		return httperror.New(http.StatusBadRequest, "Missing state cookie")
	}

	stateParam := r.URL.Query().Get("state")
	if stateParam != stateCookie.Value {
		return httperror.New(http.StatusBadRequest, "Invalid state parameter")
	}

	// Clear state cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "github_oauth_state",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	// Get code from query parameter
	code := r.URL.Query().Get("code")
	if code == "" {
		return httperror.New(http.StatusBadRequest, "Missing authorization code")
	}

	// Exchange code for token
	token, err := c.githubOAuth.ExchangeCode(code)
	if err != nil {
		return fmt.Errorf("failed to exchange code: %w", err)
	}

	// Get user info from GitHub
	githubUser, err := c.githubOAuth.GetUserInfo(token)
	if err != nil {
		return fmt.Errorf("failed to get user info: %w", err)
	}

	// Check if user already exists by GitHub ID
	githubUserID := strconv.FormatInt(githubUser.ID, 10)
	user, err := c.users.FindByGitHubUserID(githubUserID)
	if err != nil {
		return err
	}

	// If user doesn't exist by GitHub ID, check by email
	if user == nil && githubUser.Email != "" {
		user, err = c.users.FindByEmail(githubUser.Email)
		if err != nil {
			return err
		}
	}

	// Create new user if doesn't exist
	if user == nil {
		displayName := githubUser.Name
		if displayName == "" {
			displayName = githubUser.Login
		}

		user, err = c.users.CreateFromGitHub(
			githubUser.Login,
			githubUser.Email,
			displayName,
			githubUserID,
		)
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
	} else if user.GitHubUserID == nil {
		// Update existing user with GitHub ID if they signed in with password before
		user.GitHubUserID = &githubUserID
		if err := c.users.Update(user); err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}
	}

	// Set user cookie
	c.authService.SetUserCookie(w, user.ID)

	// Redirect to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

func generateRandomState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

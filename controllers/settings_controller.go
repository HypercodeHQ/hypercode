package controllers

import (
	"net/http"
	"strings"

	"github.com/hypercommithq/hypercommit/database/models"
	"github.com/hypercommithq/hypercommit/database/repositories"
	"github.com/hypercommithq/hypercommit/httperror"
	"github.com/hypercommithq/hypercommit/middleware"
	"github.com/hypercommithq/hypercommit/services"
	"github.com/hypercommithq/hypercommit/views/pages"
)

type SettingsController interface {
	Show(w http.ResponseWriter, r *http.Request) error
	UpdateGeneral(w http.ResponseWriter, r *http.Request) error
	UpdatePassword(w http.ResponseWriter, r *http.Request) error
}

type settingsController struct {
	users        repositories.UsersRepository
	accessTokens repositories.AccessTokensRepository
	authService  services.AuthService
}

func NewSettingsController(users repositories.UsersRepository, accessTokens repositories.AccessTokensRepository, authService services.AuthService) SettingsController {
	return &settingsController{
		users:        users,
		accessTokens: accessTokens,
		authService:  authService,
	}
}

func (c *settingsController) Show(w http.ResponseWriter, r *http.Request) error {
	user := middleware.GetUserFromContext(r)
	if user == nil {
		http.Redirect(w, r, "/auth/sign-in", http.StatusSeeOther)
		return nil
	}

	// Load access tokens
	tokens, err := c.accessTokens.FindByUserID(user.ID)
	if err != nil {
		return err
	}
	if tokens == nil {
		tokens = []*models.AccessToken{}
	}

	// Get flash messages from cookies
	newToken := ""
	tokenSuccess := ""
	tokenError := ""

	if cookie, err := r.Cookie("new_access_token"); err == nil {
		newToken = cookie.Value
		// Clear cookie
		http.SetCookie(w, &http.Cookie{
			Name:   "new_access_token",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
	}

	if cookie, err := r.Cookie("access_token_success"); err == nil {
		tokenSuccess = cookie.Value
		// Clear cookie
		http.SetCookie(w, &http.Cookie{
			Name:   "access_token_success",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
	}

	if cookie, err := r.Cookie("access_token_error"); err == nil {
		tokenError = cookie.Value
		// Clear cookie
		http.SetCookie(w, &http.Cookie{
			Name:   "access_token_error",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
	}

	return pages.Settings(r, &pages.SettingsData{
		User:               user,
		AccessTokens:       tokens,
		NewAccessToken:     newToken,
		AccessTokenSuccess: tokenSuccess,
		AccessTokenError:   tokenError,
	}).Render(w, r)
}

func (c *settingsController) UpdateGeneral(w http.ResponseWriter, r *http.Request) error {
	user := middleware.GetUserFromContext(r)
	if user == nil {
		http.Redirect(w, r, "/auth/sign-in", http.StatusSeeOther)
		return nil
	}

	if err := r.ParseForm(); err != nil {
		return httperror.New(http.StatusBadRequest, "Invalid form data")
	}

	displayName := strings.TrimSpace(r.FormValue("display_name"))
	username := strings.TrimSpace(r.FormValue("username"))

	settingsData := &pages.SettingsData{
		User:        user,
		DisplayName: displayName,
		Username:    username,
	}

	hasErrors := false

	if displayName == "" {
		settingsData.DisplayNameError = "Display name is required"
		hasErrors = true
	}

	if username == "" {
		settingsData.UsernameError = "Username is required"
		hasErrors = true
	}

	// Check if username is already taken by another user
	if username != user.Username {
		existingUser, err := c.users.FindByUsername(username)
		if err != nil {
			return err
		}
		if existingUser != nil {
			settingsData.UsernameError = "Username already taken"
			hasErrors = true
		}
	}

	if hasErrors {
		return pages.Settings(r, settingsData).Render(w, r)
	}

	// Update user
	user.DisplayName = displayName
	user.Username = username

	if err := c.users.Update(user); err != nil {
		return err
	}

	settingsData.GeneralSuccess = "Settings updated successfully!"
	return pages.Settings(r, settingsData).Render(w, r)
}

func (c *settingsController) UpdatePassword(w http.ResponseWriter, r *http.Request) error {
	user := middleware.GetUserFromContext(r)
	if user == nil {
		http.Redirect(w, r, "/auth/sign-in", http.StatusSeeOther)
		return nil
	}

	if err := r.ParseForm(); err != nil {
		return httperror.New(http.StatusBadRequest, "Invalid form data")
	}

	currentPassword := r.FormValue("current_password")
	newPassword := r.FormValue("new_password")
	confirmPassword := r.FormValue("confirm_password")

	settingsData := &pages.SettingsData{
		User: user,
	}

	hasErrors := false

	// Check if user has a password (not a GitHub OAuth user)
	if user.Password == nil {
		settingsData.CurrentPasswordError = "Cannot change password for OAuth accounts"
		hasErrors = true
	} else if currentPassword == "" {
		settingsData.CurrentPasswordError = "Current password is required"
		hasErrors = true
	} else if !c.authService.CheckPassword(currentPassword, *user.Password) {
		settingsData.CurrentPasswordError = "Current password is incorrect"
		hasErrors = true
	}

	if newPassword == "" {
		settingsData.NewPasswordError = "New password is required"
		hasErrors = true
	} else if len(newPassword) < 8 {
		settingsData.NewPasswordError = "Password must be at least 8 characters"
		hasErrors = true
	}

	if confirmPassword == "" {
		settingsData.ConfirmPasswordError = "Please confirm your new password"
		hasErrors = true
	} else if newPassword != confirmPassword {
		settingsData.ConfirmPasswordError = "Passwords do not match"
		hasErrors = true
	}

	if hasErrors {
		return pages.Settings(r, settingsData).Render(w, r)
	}

	// Hash new password
	hashedPassword, err := c.authService.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Update password
	user.Password = &hashedPassword
	if err := c.users.Update(user); err != nil {
		return err
	}

	settingsData.PasswordSuccess = "Password updated successfully"
	return pages.Settings(r, settingsData).Render(w, r)
}

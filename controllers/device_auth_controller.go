package controllers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hypercommithq/hypercommit/database/repositories"
	"github.com/hypercommithq/hypercommit/httperror"
	"github.com/hypercommithq/hypercommit/middleware"
	"github.com/hypercommithq/hypercommit/views/pages"
)

type DeviceAuthController interface {
	InitiateDeviceAuth(w http.ResponseWriter, r *http.Request) error
	PollDeviceAuth(w http.ResponseWriter, r *http.Request) error
	ShowDeviceAuthPage(w http.ResponseWriter, r *http.Request) error
	ConfirmDeviceAuth(w http.ResponseWriter, r *http.Request) error
}

type deviceAuthController struct {
	sessions     repositories.DeviceAuthSessionsRepository
	accessTokens repositories.AccessTokensRepository
	users        repositories.UsersRepository
}

func NewDeviceAuthController(
	sessions repositories.DeviceAuthSessionsRepository,
	accessTokens repositories.AccessTokensRepository,
	users repositories.UsersRepository,
) DeviceAuthController {
	return &deviceAuthController{
		sessions:     sessions,
		accessTokens: accessTokens,
		users:        users,
	}
}

// generateDeviceCode generates a user-friendly device code
func (c *deviceAuthController) generateDeviceCode() (string, error) {
	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789" // Exclude ambiguous characters
	const codeLength = 8

	code := make([]byte, codeLength)
	for i := range code {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		code[i] = charset[num.Int64()]
	}

	// Format as XXXX-XXXX
	return fmt.Sprintf("%s-%s", string(code[:4]), string(code[4:])), nil
}

// generateAccessToken generates a cryptographically secure access token
func (c *deviceAuthController) generateAccessToken() (string, string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", "", err
	}

	rawToken := base64.URLEncoding.EncodeToString(b)
	hash := sha256.Sum256([]byte(rawToken))
	tokenHash := fmt.Sprintf("%x", hash)

	return rawToken, tokenHash, nil
}

// InitiateDeviceAuth creates a new device auth session
func (c *deviceAuthController) InitiateDeviceAuth(w http.ResponseWriter, r *http.Request) error {
	sessionID := uuid.New().String()
	code, err := c.generateDeviceCode()
	if err != nil {
		return err
	}

	// Session expires in 10 minutes
	expiresAt := time.Now().Add(10 * time.Minute).Unix()

	session, err := c.sessions.Create(sessionID, code, expiresAt)
	if err != nil {
		return err
	}

	// Return JSON response with session details
	verificationURL := fmt.Sprintf("%s://%s/auth/device?code=%s",
		func() string {
			if r.TLS != nil {
				return "https"
			}
			return "http"
		}(),
		r.Host,
		code,
	)

	response := map[string]interface{}{
		"session_id":       session.ID,
		"user_code":        session.Code,
		"verification_url": verificationURL,
		"expires_at":       session.ExpiresAt,
		"interval":         1, // Poll every 1 second
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

// PollDeviceAuth checks the status of a device auth session
func (c *deviceAuthController) PollDeviceAuth(w http.ResponseWriter, r *http.Request) error {
	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		return httperror.New(http.StatusBadRequest, "session_id is required")
	}

	session, err := c.sessions.FindByID(sessionID)
	if err != nil {
		return err
	}

	if session == nil {
		return httperror.New(http.StatusNotFound, "Session not found")
	}

	// Check if expired
	if time.Now().Unix() > session.ExpiresAt {
		_ = c.sessions.UpdateStatus(sessionID, "expired")
		response := map[string]interface{}{
			"status": "expired",
			"error":  "The device code has expired. Please start a new authentication flow.",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusGone)
		return json.NewEncoder(w).Encode(response)
	}

	// Return current status
	response := map[string]interface{}{
		"status": session.Status,
	}

	if session.Status == "confirmed" && session.AccessToken != nil && session.UserID != nil {
		// Fetch username
		user, err := c.users.FindByID(*session.UserID)
		if err != nil {
			return err
		}

		response["access_token"] = *session.AccessToken
		response["username"] = user.Username
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

// ShowDeviceAuthPage shows the device confirmation page
func (c *deviceAuthController) ShowDeviceAuthPage(w http.ResponseWriter, r *http.Request) error {
	code := r.URL.Query().Get("code")
	if code != "" {
		// Normalize code (remove hyphens and convert to uppercase)
		code = strings.ToUpper(strings.ReplaceAll(code, "-", ""))
		// Add hyphen back in the middle for display
		if len(code) == 8 {
			code = fmt.Sprintf("%s-%s", code[:4], code[4:])
		}
	}

	user := middleware.GetUserFromContext(r)

	data := &pages.DeviceAuthData{
		User: user,
		Code: code,
	}

	return pages.DeviceAuth(r, data).Render(w, r)
}

// ConfirmDeviceAuth confirms a device auth session
func (c *deviceAuthController) ConfirmDeviceAuth(w http.ResponseWriter, r *http.Request) error {
	user := middleware.GetUserFromContext(r)
	if user == nil {
		http.Redirect(w, r, "/auth/sign-in", http.StatusSeeOther)
		return nil
	}

	if err := r.ParseForm(); err != nil {
		return httperror.New(http.StatusBadRequest, "Invalid form data")
	}

	code := r.FormValue("code")
	if code == "" {
		return httperror.New(http.StatusBadRequest, "Code is required")
	}

	// Normalize code
	code = strings.ToUpper(strings.ReplaceAll(code, "-", ""))
	if len(code) == 8 {
		code = fmt.Sprintf("%s-%s", code[:4], code[4:])
	}

	// Find session by code
	session, err := c.sessions.FindByCode(code)
	if err != nil {
		return err
	}

	if session == nil {
		data := &pages.DeviceAuthData{
			User:  user,
			Code:  code,
			Error: "Invalid or expired code. Please try again.",
		}
		return pages.DeviceAuth(r, data).Render(w, r)
	}

	// Check if expired
	if time.Now().Unix() > session.ExpiresAt {
		_ = c.sessions.UpdateStatus(session.ID, "expired")
		data := &pages.DeviceAuthData{
			User:  user,
			Code:  code,
			Error: "This code has expired. Please start a new authentication flow.",
		}
		return pages.DeviceAuth(r, data).Render(w, r)
	}

	// Check if already confirmed
	if session.Status == "confirmed" {
		data := &pages.DeviceAuthData{
			User:    user,
			Code:    code,
			Success: true,
		}
		return pages.DeviceAuth(r, data).Render(w, r)
	}

	// Generate access token
	rawToken, tokenHash, err := c.generateAccessToken()
	if err != nil {
		return err
	}

	// Create access token in database
	tokenName := fmt.Sprintf("CLI Device Auth - %s", time.Now().Format("2006-01-02 15:04:05"))
	_, err = c.accessTokens.Create(user.ID, tokenName, tokenHash)
	if err != nil {
		return err
	}

	// Confirm the session
	err = c.sessions.Confirm(session.ID, user.ID, rawToken)
	if err != nil {
		return err
	}

	// Show success page
	data := &pages.DeviceAuthData{
		User:    user,
		Code:    code,
		Success: true,
	}
	return pages.DeviceAuth(r, data).Render(w, r)
}

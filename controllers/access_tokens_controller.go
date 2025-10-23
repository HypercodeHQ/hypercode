package controllers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/hyperstitieux/hypercode/database/repositories"
	"github.com/hyperstitieux/hypercode/httperror"
	"github.com/hyperstitieux/hypercode/middleware"
)

type AccessTokensController interface {
	Create(w http.ResponseWriter, r *http.Request) error
	Delete(w http.ResponseWriter, r *http.Request) error
}

type accessTokensController struct {
	tokens repositories.AccessTokensRepository
}

func NewAccessTokensController(
	tokens repositories.AccessTokensRepository,
) AccessTokensController {
	return &accessTokensController{
		tokens: tokens,
	}
}

// generateToken generates a cryptographically secure random token
func (c *accessTokensController) generateToken() (string, string, error) {
	// Generate 32 random bytes
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", "", err
	}

	// Encode as base64 for the raw token (shown to user once)
	rawToken := base64.URLEncoding.EncodeToString(b)

	// Hash the token for storage
	hash := sha256.Sum256([]byte(rawToken))
	tokenHash := fmt.Sprintf("%x", hash)

	return rawToken, tokenHash, nil
}

func (c *accessTokensController) Create(w http.ResponseWriter, r *http.Request) error {
	user := middleware.GetUserFromContext(r)
	if user == nil {
		http.Redirect(w, r, "/auth/sign-in", http.StatusSeeOther)
		return nil
	}

	if err := r.ParseForm(); err != nil {
		return httperror.New(http.StatusBadRequest, "Invalid form data")
	}

	name := r.FormValue("name")
	if name == "" {
		// Store error in cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "access_token_error",
			Value:    "Token name is required",
			Path:     "/",
			HttpOnly: true,
			MaxAge:   10,
		})
		http.Redirect(w, r, "/settings", http.StatusSeeOther)
		return nil
	}

	// Generate token
	rawToken, tokenHash, err := c.generateToken()
	if err != nil {
		return err
	}

	// Save to database
	_, err = c.tokens.Create(user.ID, name, tokenHash)
	if err != nil {
		return err
	}

	// Flash the token to show it once
	http.SetCookie(w, &http.Cookie{
		Name:     "new_access_token",
		Value:    rawToken,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   10,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token_success",
		Value:    "Access token created successfully! Make sure to copy it now - you won't be able to see it again.",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   10,
	})
	http.Redirect(w, r, "/settings#access-tokens", http.StatusSeeOther)
	return nil
}

func (c *accessTokensController) Delete(w http.ResponseWriter, r *http.Request) error {
	user := middleware.GetUserFromContext(r)
	if user == nil {
		http.Redirect(w, r, "/auth/sign-in", http.StatusSeeOther)
		return nil
	}

	tokenIDStr := chi.URLParam(r, "id")
	tokenID, err := strconv.ParseInt(tokenIDStr, 10, 64)
	if err != nil {
		return httperror.New(http.StatusBadRequest, "Invalid token ID")
	}

	// Verify the token belongs to the user
	token, err := c.tokens.FindByID(tokenID)
	if err != nil {
		return err
	}

	if token == nil {
		return httperror.New(http.StatusNotFound, "Token not found")
	}

	if token.UserID != user.ID {
		return httperror.New(http.StatusForbidden, "You don't have permission to delete this token")
	}

	// Delete the token
	if err := c.tokens.Delete(tokenID); err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token_success",
		Value:    "Access token deleted successfully",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   10,
	})
	http.Redirect(w, r, "/settings#access-tokens", http.StatusSeeOther)
	return nil
}

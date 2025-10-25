package controllers

import (
	"net/http"
	"strings"

	"github.com/hypercommithq/hypercommit/database/repositories"
	"github.com/hypercommithq/hypercommit/httperror"
	"github.com/hypercommithq/hypercommit/services"
	"github.com/hypercommithq/hypercommit/views/pages"
)

type SignUpController interface {
	Show(w http.ResponseWriter, r *http.Request) error
	Handle(w http.ResponseWriter, r *http.Request) error
}

type signUpController struct {
	users        repositories.UsersRepository
	authService  services.AuthService
	flashService services.FlashService
}

func NewSignUpController(users repositories.UsersRepository, authService services.AuthService, flashService services.FlashService) SignUpController {
	return &signUpController{
		users:        users,
		authService:  authService,
		flashService: flashService,
	}
}

func (c *signUpController) Show(w http.ResponseWriter, r *http.Request) error {
	return pages.SignUp(r, nil).Render(w, r)
}

func (c *signUpController) Handle(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return httperror.New(http.StatusBadRequest, "Invalid form data")
	}

	username := strings.TrimSpace(r.FormValue("username"))
	email := strings.TrimSpace(r.FormValue("email"))
	displayName := strings.TrimSpace(r.FormValue("display_name"))
	password := r.FormValue("password")

	signUpData := &pages.SignUpData{
		DisplayName: displayName,
		Username:    username,
		Email:       email,
	}

	hasErrors := false

	if displayName == "" {
		signUpData.DisplayNameError = "Display name is required"
		hasErrors = true
	}

	if username == "" {
		signUpData.UsernameError = "Username is required"
		hasErrors = true
	}

	if email == "" {
		signUpData.EmailError = "Email is required"
		hasErrors = true
	}

	if password == "" {
		signUpData.PasswordError = "Password is required"
		hasErrors = true
	} else if len(password) < 8 {
		signUpData.PasswordError = "Password must be at least 8 characters"
		hasErrors = true
	}

	if hasErrors {
		return pages.SignUp(r, signUpData).Render(w, r)
	}

	existingUser, err := c.users.FindByEmail(email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		signUpData.EmailError = "Email already in use"
		return pages.SignUp(r, signUpData).Render(w, r)
	}

	existingUser, err = c.users.FindByUsername(username)
	if err != nil {
		return err
	}
	if existingUser != nil {
		signUpData.UsernameError = "Username already taken"
		return pages.SignUp(r, signUpData).Render(w, r)
	}

	hashedPassword, err := c.authService.HashPassword(password)
	if err != nil {
		return err
	}

	user, err := c.users.Create(username, email, displayName, hashedPassword)
	if err != nil {
		return err
	}

	c.authService.SetUserCookie(w, r, user.ID)
	c.flashService.Set(w, r, services.FlashCelebration)

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

package controllers

import (
	"net/http"
	"strings"

	"github.com/hyperstitieux/hypercode/database/repositories"
	"github.com/hyperstitieux/hypercode/httperror"
	"github.com/hyperstitieux/hypercode/services"
	"github.com/hyperstitieux/hypercode/views/pages"
)

type SignInController interface {
	Show(w http.ResponseWriter, r *http.Request) error
	Handle(w http.ResponseWriter, r *http.Request) error
}

type signInController struct {
	users       repositories.UsersRepository
	authService services.AuthService
}

func NewSignInController(users repositories.UsersRepository, authService services.AuthService) SignInController {
	return &signInController{
		users:       users,
		authService: authService,
	}
}

func (c *signInController) Show(w http.ResponseWriter, r *http.Request) error {
	return pages.SignIn(r, nil).Render(w, r)
}

func (c *signInController) Handle(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return httperror.New(http.StatusBadRequest, "Invalid form data")
	}

	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")

	if email == "" || password == "" {
		return pages.SignIn(r, &pages.SignInData{
			Error: "Email and password are required",
		}).Render(w, r)
	}

	user, err := c.users.FindByEmail(email)
	if err != nil {
		return err
	}

	if user == nil || !c.authService.CheckPassword(password, user.Password) {
		return pages.SignIn(r, &pages.SignInData{
			Error: "Invalid email or password",
		}).Render(w, r)
	}

	c.authService.SetUserCookie(w, user.ID)

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

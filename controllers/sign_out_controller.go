package controllers

import (
	"net/http"

	"github.com/hypercodehq/hypercode/services"
)

type SignOutController interface {
	Handle(w http.ResponseWriter, r *http.Request) error
}

type signOutController struct {
	authService services.AuthService
}

func NewSignOutController(authService services.AuthService) SignOutController {
	return &signOutController{
		authService: authService,
	}
}

func (c *signOutController) Handle(w http.ResponseWriter, r *http.Request) error {
	c.authService.ClearUserCookie(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

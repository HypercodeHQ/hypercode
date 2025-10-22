package controllers

import "net/http"

type ResetPasswordController interface {
	Show(w http.ResponseWriter, r *http.Request) error
	Handle(w http.ResponseWriter, r *http.Request) error
}

type resetPasswordController struct{}

func NewResetPasswordController() ResetPasswordController {
	return &resetPasswordController{}
}

func (c *resetPasswordController) Show(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (c *resetPasswordController) Handle(w http.ResponseWriter, r *http.Request) error {
	return nil
}

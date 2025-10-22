package controllers

import "net/http"

type ForgotPasswordController interface {
	Show(w http.ResponseWriter, r *http.Request) error
	Handle(w http.ResponseWriter, r *http.Request) error
}

type forgotPasswordController struct{}

func NewForgotPasswordController() ForgotPasswordController {
	return &forgotPasswordController{}
}

func (c *forgotPasswordController) Show(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (c *forgotPasswordController) Handle(w http.ResponseWriter, r *http.Request) error {
	return nil
}

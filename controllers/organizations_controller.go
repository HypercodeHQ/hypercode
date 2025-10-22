package controllers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hyperstitieux/hypercode/httperror"
	custommiddleware "github.com/hyperstitieux/hypercode/middleware"
)

type OrganizationsController interface {
	New(w http.ResponseWriter, r *http.Request) error
	Create(w http.ResponseWriter, r *http.Request) error
	Store(w http.ResponseWriter, r *http.Request) error
	Show(w http.ResponseWriter, r *http.Request) error
	Settings(w http.ResponseWriter, r *http.Request) error
	Update(w http.ResponseWriter, r *http.Request) error
	Delete(w http.ResponseWriter, r *http.Request) error
}

type organizationsController struct{}

func NewOrganizationsController() OrganizationsController {
	return &organizationsController{}
}

func (c *organizationsController) New(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (c *organizationsController) Create(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (c *organizationsController) Store(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (c *organizationsController) Show(w http.ResponseWriter, r *http.Request) error {
	owner := chi.URLParam(r, "owner")
	ownerType, ok := custommiddleware.GetOwnerType(r.Context())
	if !ok {
		return httperror.NotFound("owner not found")
	}

	ownerID, _ := custommiddleware.GetOwnerID(r.Context())

	w.Write([]byte("<h1>Owner: " + owner + "</h1>"))
	w.Write([]byte("<p>Type: " + string(ownerType) + "</p>"))
	w.Write([]byte("<p>ID: " + string(rune(ownerID+'0')) + "</p>"))

	return nil
}

func (c *organizationsController) Settings(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (c *organizationsController) Update(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (c *organizationsController) Delete(w http.ResponseWriter, r *http.Request) error {
	return nil
}

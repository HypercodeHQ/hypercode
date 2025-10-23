package controllers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/hyperstitieux/hypercode/database/models"
	"github.com/hyperstitieux/hypercode/database/repositories"
	"github.com/hyperstitieux/hypercode/httperror"
	custommiddleware "github.com/hyperstitieux/hypercode/middleware"
	"github.com/hyperstitieux/hypercode/services"
	"github.com/hyperstitieux/hypercode/views/pages"
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

type organizationsController struct {
	orgs        repositories.OrganizationsRepository
	users       repositories.UsersRepository
	repos       repositories.RepositoriesRepository
	authService services.AuthService
}

func NewOrganizationsController(orgs repositories.OrganizationsRepository, users repositories.UsersRepository, repos repositories.RepositoriesRepository, authService services.AuthService) OrganizationsController {
	return &organizationsController{
		orgs:        orgs,
		users:       users,
		repos:       repos,
		authService: authService,
	}
}

func (c *organizationsController) New(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (c *organizationsController) Create(w http.ResponseWriter, r *http.Request) error {
	user, err := c.authService.GetUserFromCookie(r)
	if err != nil {
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
		return nil
	}

	if user == nil {
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
		return nil
	}

	return pages.NewOrganization(r, nil).Render(w, r)
}

func (c *organizationsController) Store(w http.ResponseWriter, r *http.Request) error {
	user, err := c.authService.GetUserFromCookie(r)
	if err != nil {
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
		return nil
	}

	if user == nil {
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
		return nil
	}

	if err := r.ParseForm(); err != nil {
		return httperror.BadRequest("invalid form data")
	}

	username := r.FormValue("username")
	displayName := r.FormValue("display_name")

	orgData := &pages.NewOrganizationData{
		Username:    username,
		DisplayName: displayName,
	}

	hasErrors := false

	if username == "" {
		orgData.UsernameError = "Organization username is required"
		hasErrors = true
	}

	if displayName == "" {
		orgData.DisplayNameError = "Display name is required"
		hasErrors = true
	}

	if hasErrors {
		return pages.NewOrganization(r, orgData).Render(w, r)
	}

	existingOrg, err := c.orgs.FindByUsername(username)
	if err != nil {
		slog.Error("failed to check for existing organization", "error", err)
	}
	if existingOrg != nil {
		orgData.UsernameError = "Organization username already exists"
		return pages.NewOrganization(r, orgData).Render(w, r)
	}

	org, err := c.orgs.Create(username, displayName)
	if err != nil {
		slog.Error("failed to create organization", "error", err)
		orgData.UsernameError = "Failed to create organization"
		return pages.NewOrganization(r, orgData).Render(w, r)
	}

	slog.Info("organization created", "username", username, "displayName", displayName, "creator", user.Username)

	http.Redirect(w, r, fmt.Sprintf("/%s", org.Username), http.StatusSeeOther)
	return nil
}

func (c *organizationsController) Show(w http.ResponseWriter, r *http.Request) error {
	ownerType, ok := custommiddleware.GetOwnerType(r.Context())
	if !ok {
		return httperror.NotFound("owner not found")
	}

	ownerID, _ := custommiddleware.GetOwnerID(r.Context())

	currentUser := custommiddleware.GetUserFromContext(r)

	// Handle user profile
	if ownerType == custommiddleware.OwnerTypeUser {
		profileUser, err := c.users.FindByID(ownerID)
		if err != nil {
			return httperror.NotFound("user not found")
		}

		userRepos, err := c.repos.FindAllByUser(ownerID)
		if err != nil {
			slog.Error("failed to fetch user repositories", "error", err)
			userRepos = []*models.Repository{}
		}

		return pages.UserProfile(r, &pages.UserProfileData{
			User:         currentUser,
			ProfileUser:  profileUser,
			Repositories: userRepos,
		}).Render(w, r)
	}

	// Handle organization profile (TODO: implement organization profile page)
	if ownerType == custommiddleware.OwnerTypeOrg {
		http.Error(w, "Organization profiles not yet implemented", http.StatusNotImplemented)
		return nil
	}

	return httperror.NotFound("owner not found")
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

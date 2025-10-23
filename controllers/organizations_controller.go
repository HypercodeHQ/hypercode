package controllers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/hypercodehq/hypercode/database/models"
	"github.com/hypercodehq/hypercode/database/repositories"
	"github.com/hypercodehq/hypercode/httperror"
	custommiddleware "github.com/hypercodehq/hypercode/middleware"
	"github.com/hypercodehq/hypercode/services"
	"github.com/hypercodehq/hypercode/views/pages"
)

type OrganizationsController interface {
	New(w http.ResponseWriter, r *http.Request) error
	Create(w http.ResponseWriter, r *http.Request) error
	Store(w http.ResponseWriter, r *http.Request) error
	Show(w http.ResponseWriter, r *http.Request) error
	Repositories(w http.ResponseWriter, r *http.Request) error
	Stars(w http.ResponseWriter, r *http.Request) error
	Settings(w http.ResponseWriter, r *http.Request) error
	Update(w http.ResponseWriter, r *http.Request) error
	Delete(w http.ResponseWriter, r *http.Request) error
}

type organizationsController struct {
	orgs        repositories.OrganizationsRepository
	users       repositories.UsersRepository
	repos       repositories.RepositoriesRepository
	stars       repositories.StarsRepository
	authService services.AuthService
}

func NewOrganizationsController(orgs repositories.OrganizationsRepository, users repositories.UsersRepository, repos repositories.RepositoriesRepository, stars repositories.StarsRepository, authService services.AuthService) OrganizationsController {
	return &organizationsController{
		orgs:        orgs,
		users:       users,
		repos:       repos,
		stars:       stars,
		authService: authService,
	}
}

func (c *organizationsController) New(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (c *organizationsController) Create(w http.ResponseWriter, r *http.Request) error {
	user, err := c.authService.GetUserFromCookie(r)
	if err != nil {
		http.Redirect(w, r, "/auth/sign-in", http.StatusSeeOther)
		return nil
	}

	if user == nil {
		http.Redirect(w, r, "/auth/sign-in", http.StatusSeeOther)
		return nil
	}

	return pages.NewOrganization(r, nil).Render(w, r)
}

func (c *organizationsController) Store(w http.ResponseWriter, r *http.Request) error {
	user, err := c.authService.GetUserFromCookie(r)
	if err != nil {
		http.Redirect(w, r, "/auth/sign-in", http.StatusSeeOther)
		return nil
	}

	if user == nil {
		http.Redirect(w, r, "/auth/sign-in", http.StatusSeeOther)
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

		// Fetch star counts for each repository
		starCounts := make(map[int64]int64)
		for _, repo := range userRepos {
			starCount, err := c.stars.CountByRepository(repo.ID)
			if err != nil {
				starCount = 0
			}
			starCounts[repo.ID] = starCount
		}

		// User profiles don't have settings tab
		return pages.UserProfile(r, &pages.UserProfileData{
			User:         currentUser,
			ProfileUser:  profileUser,
			Repositories: userRepos,
			StarCounts:   starCounts,
			CanManage:    false,
			CurrentTab:   "overview",
		}).Render(w, r)
	}

	// Handle organization profile
	if ownerType == custommiddleware.OwnerTypeOrg {
		org, err := c.orgs.FindByID(ownerID)
		if err != nil {
			return httperror.NotFound("organization not found")
		}

		orgRepos, err := c.repos.FindAllByOrg(ownerID)
		if err != nil {
			slog.Error("failed to fetch organization repositories", "error", err)
			orgRepos = []*models.Repository{}
		}

		// Fetch star counts for each repository
		starCounts := make(map[int64]int64)
		for _, repo := range orgRepos {
			starCount, err := c.stars.CountByRepository(repo.ID)
			if err != nil {
				starCount = 0
			}
			starCounts[repo.ID] = starCount
		}

		// Check if current user can manage this organization
		// TODO: Implement proper organization membership/ownership check
		canManage := false
		if currentUser != nil {
			// For now, we'll need to implement organization membership
			// This is a placeholder - proper implementation needed
			canManage = false
		}

		return pages.OrganizationProfile(r, &pages.OrganizationProfileData{
			User:         currentUser,
			Organization: org,
			Repositories: orgRepos,
			StarCounts:   starCounts,
			CanManage:    canManage,
			CurrentTab:   "overview",
		}).Render(w, r)
	}

	return httperror.NotFound("owner not found")
}

func (c *organizationsController) Repositories(w http.ResponseWriter, r *http.Request) error {
	ownerType, ok := custommiddleware.GetOwnerType(r.Context())
	if !ok {
		return httperror.NotFound("owner not found")
	}

	ownerID, _ := custommiddleware.GetOwnerID(r.Context())

	currentUser := custommiddleware.GetUserFromContext(r)

	// Handle user repositories tab
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

		// Fetch star counts for each repository
		starCounts := make(map[int64]int64)
		for _, repo := range userRepos {
			starCount, err := c.stars.CountByRepository(repo.ID)
			if err != nil {
				starCount = 0
			}
			starCounts[repo.ID] = starCount
		}

		// User profiles don't have settings tab
		return pages.UserProfile(r, &pages.UserProfileData{
			User:         currentUser,
			ProfileUser:  profileUser,
			Repositories: userRepos,
			StarCounts:   starCounts,
			CanManage:    false,
			CurrentTab:   "repositories",
		}).Render(w, r)
	}

	// Handle organization repositories tab
	if ownerType == custommiddleware.OwnerTypeOrg {
		org, err := c.orgs.FindByID(ownerID)
		if err != nil {
			return httperror.NotFound("organization not found")
		}

		orgRepos, err := c.repos.FindAllByOrg(ownerID)
		if err != nil {
			slog.Error("failed to fetch organization repositories", "error", err)
			orgRepos = []*models.Repository{}
		}

		// Fetch star counts for each repository
		starCounts := make(map[int64]int64)
		for _, repo := range orgRepos {
			starCount, err := c.stars.CountByRepository(repo.ID)
			if err != nil {
				starCount = 0
			}
			starCounts[repo.ID] = starCount
		}

		// Check if current user can manage this organization
		canManage := false
		if currentUser != nil {
			// TODO: Implement proper organization membership/ownership check
			canManage = false
		}

		return pages.OrganizationProfile(r, &pages.OrganizationProfileData{
			User:         currentUser,
			Organization: org,
			Repositories: orgRepos,
			StarCounts:   starCounts,
			CanManage:    canManage,
			CurrentTab:   "repositories",
		}).Render(w, r)
	}

	return httperror.NotFound("owner not found")
}

func (c *organizationsController) Stars(w http.ResponseWriter, r *http.Request) error {
	ownerType, ok := custommiddleware.GetOwnerType(r.Context())
	if !ok {
		return httperror.NotFound("owner not found")
	}

	// Organizations cannot star repositories
	if ownerType == custommiddleware.OwnerTypeOrg {
		return httperror.NotFound("organizations cannot star repositories")
	}

	ownerID, _ := custommiddleware.GetOwnerID(r.Context())
	currentUser := custommiddleware.GetUserFromContext(r)

	// Handle user stars tab
	if ownerType == custommiddleware.OwnerTypeUser {
		profileUser, err := c.users.FindByID(ownerID)
		if err != nil {
			return httperror.NotFound("user not found")
		}

		// Fetch repositories starred by the user
		starredRepos, err := c.stars.FindStarredRepositoriesByUser(ownerID)
		if err != nil {
			slog.Error("failed to fetch starred repositories", "error", err)
			starredRepos = []*models.Repository{}
		}

		// Build repository with owner information
		reposWithOwner := []pages.RepositoryWithOwner{}
		for _, repo := range starredRepos {
			// Get owner username
			var ownerUsername string
			if repo.OwnerUserID != nil {
				ownerUser, err := c.users.FindByID(*repo.OwnerUserID)
				if err == nil && ownerUser != nil {
					ownerUsername = ownerUser.Username
				}
			} else if repo.OwnerOrgID != nil {
				ownerOrg, err := c.orgs.FindByID(*repo.OwnerOrgID)
				if err == nil && ownerOrg != nil {
					ownerUsername = ownerOrg.Username
				}
			}

			// Get star count
			starCount, err := c.stars.CountByRepository(repo.ID)
			if err != nil {
				starCount = 0
			}

			reposWithOwner = append(reposWithOwner, pages.RepositoryWithOwner{
				Repository:    repo,
				OwnerUsername: ownerUsername,
				StarCount:     starCount,
			})
		}

		// User profiles don't have settings tab
		return pages.UserProfile(r, &pages.UserProfileData{
			User:                  currentUser,
			ProfileUser:           profileUser,
			RepositoriesWithOwner: reposWithOwner,
			CanManage:             false,
			CurrentTab:            "stars",
		}).Render(w, r)
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

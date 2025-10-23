package controllers

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/hyperstitieux/hypercode/database/models"
	"github.com/hyperstitieux/hypercode/database/repositories"
	"github.com/hyperstitieux/hypercode/httperror"
	"github.com/hyperstitieux/hypercode/services"
	"github.com/hyperstitieux/hypercode/views/pages"
)

type RepositoriesController interface {
	Create(w http.ResponseWriter, r *http.Request) error
	Store(w http.ResponseWriter, r *http.Request) error
	Show(w http.ResponseWriter, r *http.Request) error
	Settings(w http.ResponseWriter, r *http.Request) error
	UpdateSettings(w http.ResponseWriter, r *http.Request) error
	Delete(w http.ResponseWriter, r *http.Request) error
	Star(w http.ResponseWriter, r *http.Request) error
	Unstar(w http.ResponseWriter, r *http.Request) error
}

type repositoriesController struct {
	repos         repositories.RepositoriesRepository
	users         repositories.UsersRepository
	contributors  repositories.ContributorsRepository
	stars         repositories.StarsRepository
	orgs          repositories.OrganizationsRepository
	authService   services.AuthService
	reposBasePath string
}

func NewRepositoriesController(
	repos repositories.RepositoriesRepository,
	users repositories.UsersRepository,
	contributors repositories.ContributorsRepository,
	stars repositories.StarsRepository,
	orgs repositories.OrganizationsRepository,
	authService services.AuthService,
	reposBasePath string,
) RepositoriesController {
	return &repositoriesController{
		repos:         repos,
		users:         users,
		contributors:  contributors,
		stars:         stars,
		orgs:          orgs,
		authService:   authService,
		reposBasePath: reposBasePath,
	}
}

func (c *repositoriesController) Create(w http.ResponseWriter, r *http.Request) error {
	user, err := c.authService.GetUserFromCookie(r)
	if err != nil {
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
		return nil
	}

	if user == nil {
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
		return nil
	}

	// Get all organizations (for now, show all orgs - later add membership check)
	// TODO: Filter by organizations where user is a member
	orgs, err := c.orgs.FindAll()
	if err != nil {
		slog.Error("failed to fetch organizations", "error", err)
		orgs = []*models.Organization{}
	}

	return pages.NewRepository(r, &pages.NewRepositoryData{
		User:          user,
		Organizations: orgs,
		DefaultBranch: "main",
	}).Render(w, r)
}

func (c *repositoriesController) Store(w http.ResponseWriter, r *http.Request) error {
	user, err := c.authService.GetUserFromCookie(r)
	if err != nil {
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
		return nil
	}

	if err := r.ParseForm(); err != nil {
		return httperror.BadRequest("invalid form data")
	}

	name := r.FormValue("name")
	visibility := r.FormValue("visibility")
	defaultBranch := r.FormValue("default_branch")
	ownerUsername := r.FormValue("owner")

	if visibility == "" {
		visibility = "public"
	}

	if defaultBranch == "" {
		defaultBranch = "main"
	}

	if ownerUsername == "" {
		ownerUsername = user.Username
	}

	// Get all organizations for error handling
	orgs, _ := c.orgs.FindAll()

	repoData := &pages.NewRepositoryData{
		Name:          name,
		DefaultBranch: defaultBranch,
		Visibility:    visibility,
		Owner:         ownerUsername,
		User:          user,
		Organizations: orgs,
	}

	hasErrors := false

	if name == "" {
		repoData.NameError = "Repository name is required"
		hasErrors = true
	}

	if hasErrors {
		return pages.NewRepository(r, repoData).Render(w, r)
	}

	// Determine if owner is user or organization
	var repo *models.Repository
	var ownerIDForPath string

	if ownerUsername == user.Username {
		// Check for existing repo under user
		existingRepo, err := c.repos.FindByUserAndName(user.ID, name)
		if err != nil {
			slog.Error("failed to check for existing repository", "error", err)
		}
		if existingRepo != nil {
			repoData.NameError = "Repository name already exists"
			return pages.NewRepository(r, repoData).Render(w, r)
		}

		// Create for user
		repo, err = c.repos.CreateForUser(user.ID, name, visibility, defaultBranch, nil)
		if err != nil {
			slog.Error("failed to create repository", "error", err)
			repoData.NameError = "Failed to create repository"
			return pages.NewRepository(r, repoData).Render(w, r)
		}
		ownerIDForPath = fmt.Sprintf("%d", user.ID)
	} else {
		// Find organization
		org, err := c.orgs.FindByUsername(ownerUsername)
		if err != nil {
			slog.Error("failed to find organization", "error", err)
			repoData.NameError = "Organization not found"
			return pages.NewRepository(r, repoData).Render(w, r)
		}
		if org == nil {
			repoData.NameError = "Organization not found"
			return pages.NewRepository(r, repoData).Render(w, r)
		}

		// Check for existing repo under org
		existingRepo, err := c.repos.FindByOrgAndName(org.ID, name)
		if err != nil {
			slog.Error("failed to check for existing repository", "error", err)
		}
		if existingRepo != nil {
			repoData.NameError = "Repository name already exists"
			return pages.NewRepository(r, repoData).Render(w, r)
		}

		// Create for organization
		repo, err = c.repos.CreateForOrg(org.ID, name, visibility, defaultBranch, nil)
		if err != nil {
			slog.Error("failed to create repository", "error", err)
			repoData.NameError = "Failed to create repository"
			return pages.NewRepository(r, repoData).Render(w, r)
		}
		ownerIDForPath = fmt.Sprintf("org_%d", org.ID)
	}

	// Create admin contributor
	_, err = c.contributors.Create(repo.ID, user.ID, "admin")
	if err != nil {
		slog.Error("failed to create admin contributor", "error", err)
	}

	repoPath := filepath.Join(c.reposBasePath, ownerIDForPath, name)
	if err := os.MkdirAll(repoPath, 0755); err != nil {
		slog.Error("failed to create repository directory", "error", err)
		return httperror.New(http.StatusInternalServerError, "failed to create repository")
	}

	cmd := exec.Command("git", "init", "--bare")
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		slog.Error("failed to initialize git repository", "error", err)
		return httperror.New(http.StatusInternalServerError, "failed to initialize repository")
	}

	configCmd := exec.Command("git", "config", "http.receivepack", "true")
	configCmd.Dir = repoPath
	if err := configCmd.Run(); err != nil {
		slog.Error("failed to configure git repository", "error", err)
		return httperror.New(http.StatusInternalServerError, "failed to configure repository")
	}

	slog.Info("repository created", "owner", ownerUsername, "name", name, "visibility", visibility, "creator", user.Username)

	http.Redirect(w, r, fmt.Sprintf("/%s/%s", ownerUsername, name), http.StatusSeeOther)
	return nil
}

func (c *repositoriesController) Show(w http.ResponseWriter, r *http.Request) error {
	owner := chi.URLParam(r, "owner")
	repoName := chi.URLParam(r, "repo")

	repo, err := c.repos.FindByOwnerAndName(owner, repoName)
	if err != nil {
		return httperror.NotFound("repository not found")
	}

	if repo == nil {
		return httperror.NotFound("repository not found")
	}

	if repo.Visibility == "private" {
		user, err := c.authService.GetUserFromCookie(r)
		if err != nil || user == nil {
			return httperror.Unauthorized("authentication required")
		}
		if repo.OwnerUserID != nil && *repo.OwnerUserID != user.ID {
			return httperror.Forbidden("access denied")
		}
	}

	user, _ := c.authService.GetUserFromCookie(r)

	// Determine if user can manage the repository
	canManage := false
	if user != nil && repo.OwnerUserID != nil && *repo.OwnerUserID == user.ID {
		canManage = true
	} else if user != nil {
		// Check if user is an admin contributor
		contributor, err := c.contributors.FindByRepositoryAndUser(repo.ID, user.ID)
		if err == nil && contributor != nil && contributor.Role == "admin" {
			canManage = true
		}
	}

	// Get star count
	starCount, err := c.stars.CountByRepository(repo.ID)
	if err != nil {
		slog.Error("failed to count stars", "error", err)
		starCount = 0
	}

	// Check if user has starred
	hasStarred := false
	if user != nil {
		star, err := c.stars.FindByUserAndRepository(repo.ID, user.ID)
		if err != nil {
			slog.Error("failed to check if user starred", "error", err)
		}
		if star != nil {
			hasStarred = true
		}
	}

	host := r.Host
	cloneURL := fmt.Sprintf("https://%s/%s/%s", host, owner, repoName)

	data := &pages.ShowRepositoryData{
		User:          user,
		Repository:    repo,
		OwnerUsername: owner,
		CloneURL:      cloneURL,
		IsPublic:      repo.Visibility == "public",
		CanManage:     canManage,
		StarCount:     starCount,
		HasStarred:    hasStarred,
	}

	return pages.ShowRepository(r, data).Render(w, r)
}

func (c *repositoriesController) Settings(w http.ResponseWriter, r *http.Request) error {
	owner := chi.URLParam(r, "owner")
	repoName := chi.URLParam(r, "repo")

	repo, err := c.repos.FindByOwnerAndName(owner, repoName)
	if err != nil {
		return httperror.NotFound("repository not found")
	}

	if repo == nil {
		return httperror.NotFound("repository not found")
	}

	user, _ := c.authService.GetUserFromCookie(r)
	if user == nil {
		return httperror.Unauthorized("authentication required")
	}

	// Check if user can manage the repository
	canManage := false
	if repo.OwnerUserID != nil && *repo.OwnerUserID == user.ID {
		canManage = true
	} else {
		contributor, err := c.contributors.FindByRepositoryAndUser(repo.ID, user.ID)
		if err == nil && contributor != nil && contributor.Role == "admin" {
			canManage = true
		}
	}

	if !canManage {
		return httperror.Forbidden("access denied")
	}

	return pages.RepositorySettings(r, &pages.RepositorySettingsData{
		User:          user,
		Repository:    repo,
		OwnerUsername: owner,
	}).Render(w, r)
}

func (c *repositoriesController) UpdateSettings(w http.ResponseWriter, r *http.Request) error {
	owner := chi.URLParam(r, "owner")
	repoName := chi.URLParam(r, "repo")

	repo, err := c.repos.FindByOwnerAndName(owner, repoName)
	if err != nil {
		return httperror.NotFound("repository not found")
	}

	if repo == nil {
		return httperror.NotFound("repository not found")
	}

	user, _ := c.authService.GetUserFromCookie(r)
	if user == nil {
		return httperror.Unauthorized("authentication required")
	}

	// Check if user can manage the repository
	canManage := false
	if repo.OwnerUserID != nil && *repo.OwnerUserID == user.ID {
		canManage = true
	} else {
		contributor, err := c.contributors.FindByRepositoryAndUser(repo.ID, user.ID)
		if err == nil && contributor != nil && contributor.Role == "admin" {
			canManage = true
		}
	}

	if !canManage {
		return httperror.Forbidden("access denied")
	}

	if err := r.ParseForm(); err != nil {
		return httperror.BadRequest("invalid form data")
	}

	name := r.FormValue("name")
	defaultBranch := r.FormValue("default_branch")
	visibility := r.FormValue("visibility")

	settingsData := &pages.RepositorySettingsData{
		User:          user,
		Repository:    repo,
		OwnerUsername: owner,
		Name:          name,
		DefaultBranch: defaultBranch,
		Visibility:    visibility,
	}

	hasErrors := false

	if name == "" {
		settingsData.NameError = "Repository name is required"
		hasErrors = true
	}

	if visibility == "" {
		visibility = "public"
	}

	// Check if name changed and if new name is already taken
	if name != repo.Name {
		existingRepo, err := c.repos.FindByOwnerAndName(owner, name)
		if err != nil {
			slog.Error("failed to check for existing repository", "error", err)
		}
		if existingRepo != nil {
			settingsData.NameError = "Repository name already exists"
			hasErrors = true
		}
	}

	if hasErrors {
		return pages.RepositorySettings(r, settingsData).Render(w, r)
	}

	// Update repository
	repo.Name = name
	repo.DefaultBranch = defaultBranch
	repo.Visibility = visibility

	if err := c.repos.Update(repo); err != nil {
		slog.Error("failed to update repository", "error", err)
		settingsData.NameError = "Failed to update repository"
		return pages.RepositorySettings(r, settingsData).Render(w, r)
	}

	settingsData.GeneralSuccess = "Settings updated successfully!"

	// If name changed, redirect to new URL
	if name != repoName {
		http.Redirect(w, r, fmt.Sprintf("/%s/%s/settings", owner, name), http.StatusSeeOther)
		return nil
	}

	return pages.RepositorySettings(r, settingsData).Render(w, r)
}

func (c *repositoriesController) Delete(w http.ResponseWriter, r *http.Request) error {
	owner := chi.URLParam(r, "owner")
	repoName := chi.URLParam(r, "repo")

	repo, err := c.repos.FindByOwnerAndName(owner, repoName)
	if err != nil {
		return httperror.NotFound("repository not found")
	}

	if repo == nil {
		return httperror.NotFound("repository not found")
	}

	user, _ := c.authService.GetUserFromCookie(r)
	if user == nil {
		return httperror.Unauthorized("authentication required")
	}

	// Check if user can manage the repository
	canManage := false
	if repo.OwnerUserID != nil && *repo.OwnerUserID == user.ID {
		canManage = true
	} else {
		contributor, err := c.contributors.FindByRepositoryAndUser(repo.ID, user.ID)
		if err == nil && contributor != nil && contributor.Role == "admin" {
			canManage = true
		}
	}

	if !canManage {
		return httperror.Forbidden("access denied")
	}

	// Delete from database
	if err := c.repos.Delete(repo.ID); err != nil {
		slog.Error("failed to delete repository", "error", err)
		return httperror.New(http.StatusInternalServerError, "failed to delete repository")
	}

	// Delete repository directory
	var ownerIDForPath string
	if repo.OwnerUserID != nil {
		ownerIDForPath = fmt.Sprintf("%d", *repo.OwnerUserID)
	} else if repo.OwnerOrgID != nil {
		ownerIDForPath = fmt.Sprintf("org_%d", *repo.OwnerOrgID)
	}

	repoPath := filepath.Join(c.reposBasePath, ownerIDForPath, fmt.Sprintf("%d", repo.ID))
	if err := os.RemoveAll(repoPath); err != nil {
		slog.Error("failed to delete repository directory", "error", err)
	}

	slog.Info("repository deleted", "owner", owner, "name", repoName)

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

func (c *repositoriesController) Star(w http.ResponseWriter, r *http.Request) error {
	owner := chi.URLParam(r, "owner")
	repoName := chi.URLParam(r, "repo")

	repo, err := c.repos.FindByOwnerAndName(owner, repoName)
	if err != nil {
		return httperror.NotFound("repository not found")
	}

	if repo == nil {
		return httperror.NotFound("repository not found")
	}

	user, _ := c.authService.GetUserFromCookie(r)
	if user == nil {
		return httperror.Unauthorized("authentication required")
	}

	// Check if already starred
	existingStar, err := c.stars.FindByUserAndRepository(repo.ID, user.ID)
	if err != nil {
		slog.Error("failed to check existing star", "error", err)
	}

	if existingStar != nil {
		// Already starred, just redirect back
		referer := r.Header.Get("Referer")
		if referer == "" {
			referer = fmt.Sprintf("/%s/%s", owner, repoName)
		}
		http.Redirect(w, r, referer, http.StatusSeeOther)
		return nil
	}

	// Create star
	_, err = c.stars.Create(repo.ID, user.ID)
	if err != nil {
		slog.Error("failed to star repository", "error", err)
		return httperror.New(http.StatusInternalServerError, "failed to star repository")
	}

	slog.Info("repository starred", "user", user.Username, "repo", repoName)

	referer := r.Header.Get("Referer")
	if referer == "" {
		referer = fmt.Sprintf("/%s/%s", owner, repoName)
	}
	http.Redirect(w, r, referer, http.StatusSeeOther)
	return nil
}

func (c *repositoriesController) Unstar(w http.ResponseWriter, r *http.Request) error {
	owner := chi.URLParam(r, "owner")
	repoName := chi.URLParam(r, "repo")

	repo, err := c.repos.FindByOwnerAndName(owner, repoName)
	if err != nil {
		return httperror.NotFound("repository not found")
	}

	if repo == nil {
		return httperror.NotFound("repository not found")
	}

	user, _ := c.authService.GetUserFromCookie(r)
	if user == nil {
		return httperror.Unauthorized("authentication required")
	}

	// Delete star
	err = c.stars.Delete(repo.ID, user.ID)
	if err != nil {
		slog.Error("failed to unstar repository", "error", err)
		return httperror.New(http.StatusInternalServerError, "failed to unstar repository")
	}

	slog.Info("repository unstarred", "user", user.Username, "repo", repoName)

	referer := r.Header.Get("Referer")
	if referer == "" {
		referer = fmt.Sprintf("/%s/%s", owner, repoName)
	}
	http.Redirect(w, r, referer, http.StatusSeeOther)
	return nil
}

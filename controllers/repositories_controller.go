package controllers

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/hyperstitieux/hypercode/database/repositories"
	"github.com/hyperstitieux/hypercode/httperror"
	"github.com/hyperstitieux/hypercode/services"
	"github.com/hyperstitieux/hypercode/views/pages"
)

type RepositoriesController interface {
	Create(w http.ResponseWriter, r *http.Request) error
	Store(w http.ResponseWriter, r *http.Request) error
	Show(w http.ResponseWriter, r *http.Request) error
	Delete(w http.ResponseWriter, r *http.Request) error
}

type repositoriesController struct {
	repos         repositories.RepositoriesRepository
	users         repositories.UsersRepository
	contributors  repositories.ContributorsRepository
	authService   services.AuthService
	reposBasePath string
}

func NewRepositoriesController(
	repos repositories.RepositoriesRepository,
	users repositories.UsersRepository,
	contributors repositories.ContributorsRepository,
	authService services.AuthService,
	reposBasePath string,
) RepositoriesController {
	return &repositoriesController{
		repos:         repos,
		users:         users,
		contributors:  contributors,
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

	return pages.NewRepository(r, nil).Render(w, r)
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

	if visibility == "" {
		visibility = "public"
	}

	if defaultBranch == "" {
		defaultBranch = "main"
	}

	repoData := &pages.NewRepositoryData{
		Name:          name,
		DefaultBranch: defaultBranch,
		Visibility:    visibility,
	}

	hasErrors := false

	if name == "" {
		repoData.NameError = "Repository name is required"
		hasErrors = true
	}

	if hasErrors {
		return pages.NewRepository(r, repoData).Render(w, r)
	}

	existingRepo, err := c.repos.FindByUserAndName(user.ID, name)
	if err != nil {
		slog.Error("failed to check for existing repository", "error", err)
	}
	if existingRepo != nil {
		repoData.NameError = "Repository name already exists"
		return pages.NewRepository(r, repoData).Render(w, r)
	}

	repo, err := c.repos.CreateForUser(user.ID, name, visibility, nil)
	if err != nil {
		slog.Error("failed to create repository", "error", err)
		repoData.NameError = "Failed to create repository"
		return pages.NewRepository(r, repoData).Render(w, r)
	}

	_, err = c.contributors.Create(repo.ID, user.ID, "admin")
	if err != nil {
		slog.Error("failed to create admin contributor", "error", err)
	}

	repoPath := filepath.Join(c.reposBasePath, fmt.Sprintf("%d", user.ID), name)
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

	slog.Info("repository created", "owner", user.Username, "name", name, "visibility", visibility)

	http.Redirect(w, r, fmt.Sprintf("/%s/%s", user.Username, name), http.StatusSeeOther)
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

	host := r.Host
	cloneURL := fmt.Sprintf("https://%s/%s/%s", host, owner, repoName)

	data := &pages.ShowRepositoryData{
		User:          user,
		Repository:    repo,
		OwnerUsername: owner,
		CloneURL:      cloneURL,
		IsPublic:      repo.Visibility == "public",
	}

	return pages.ShowRepository(r, data).Render(w, r)
}

func (c *repositoriesController) Delete(w http.ResponseWriter, r *http.Request) error {
	return nil
}

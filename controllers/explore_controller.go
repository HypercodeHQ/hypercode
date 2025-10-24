package controllers

import (
	"log/slog"
	"net/http"

	"github.com/hypercommithq/hypercommit/database/models"
	"github.com/hypercommithq/hypercommit/database/repositories"
	custommiddleware "github.com/hypercommithq/hypercommit/middleware"
	"github.com/hypercommithq/hypercommit/services"
	"github.com/hypercommithq/hypercommit/views/pages"
)

type ExploreController interface {
	Repositories(w http.ResponseWriter, r *http.Request) error
	Users(w http.ResponseWriter, r *http.Request) error
	Organizations(w http.ResponseWriter, r *http.Request) error
}

type exploreController struct {
	repos       repositories.RepositoriesRepository
	users       repositories.UsersRepository
	orgs        repositories.OrganizationsRepository
	stars       repositories.StarsRepository
	authService services.AuthService
}

func NewExploreController(
	repos repositories.RepositoriesRepository,
	users repositories.UsersRepository,
	orgs repositories.OrganizationsRepository,
	stars repositories.StarsRepository,
	authService services.AuthService,
) ExploreController {
	return &exploreController{
		repos:       repos,
		users:       users,
		orgs:        orgs,
		stars:       stars,
		authService: authService,
	}
}

func (c *exploreController) Repositories(w http.ResponseWriter, r *http.Request) error {
	currentUser := custommiddleware.GetUserFromContext(r)

	allRepos, err := c.repos.FindPublic()
	if err != nil {
		slog.Error("failed to fetch all repositories", "error", err)
		allRepos = []*models.Repository{}
	}

	// Build RepositoryWithOwner slice
	reposWithOwner := make([]pages.RepositoryWithOwner, 0, len(allRepos))
	for _, repo := range allRepos {
		var ownerUsername string
		if repo.OwnerUserID != nil {
			user, err := c.users.FindByID(*repo.OwnerUserID)
			if err == nil && user != nil {
				ownerUsername = user.Username
			}
		} else if repo.OwnerOrgID != nil {
			org, err := c.orgs.FindByID(*repo.OwnerOrgID)
			if err == nil && org != nil {
				ownerUsername = org.Username
			}
		}

		if ownerUsername != "" {
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
	}

	return pages.ExploreRepositories(r, &pages.ExploreRepositoriesData{
		User:         currentUser,
		Repositories: reposWithOwner,
	}).Render(w, r)
}

func (c *exploreController) Users(w http.ResponseWriter, r *http.Request) error {
	currentUser := custommiddleware.GetUserFromContext(r)

	allUsers, err := c.users.FindAll()
	if err != nil {
		slog.Error("failed to fetch all users", "error", err)
		allUsers = []*models.User{}
	}

	return pages.ExploreUsers(r, &pages.ExploreUsersData{
		User:  currentUser,
		Users: allUsers,
	}).Render(w, r)
}

func (c *exploreController) Organizations(w http.ResponseWriter, r *http.Request) error {
	currentUser := custommiddleware.GetUserFromContext(r)

	allOrgs, err := c.orgs.FindAll()
	if err != nil {
		slog.Error("failed to fetch all organizations", "error", err)
		allOrgs = []*models.Organization{}
	}

	return pages.ExploreOrganizations(r, &pages.ExploreOrganizationsData{
		User:          currentUser,
		Organizations: allOrgs,
	}).Render(w, r)
}

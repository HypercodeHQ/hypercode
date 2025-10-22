package controllers

import (
	"net/http"

	"github.com/hyperstitieux/hypercode/database/repositories"
	"github.com/hyperstitieux/hypercode/middleware"
	"github.com/hyperstitieux/hypercode/views/pages"
)

type HomeController interface {
	Show(w http.ResponseWriter, r *http.Request) error
}

type homeController struct {
	repos repositories.RepositoriesRepository
	users repositories.UsersRepository
	orgs  repositories.OrganizationsRepository
}

func NewHomeController(
	repos repositories.RepositoriesRepository,
	users repositories.UsersRepository,
	orgs repositories.OrganizationsRepository,
) HomeController {
	return &homeController{
		repos: repos,
		users: users,
		orgs:  orgs,
	}
}

func (c *homeController) Show(w http.ResponseWriter, r *http.Request) error {
	user := middleware.GetUserFromContext(r)

	data := &pages.HomeData{
		User:         user,
		Repositories: []pages.RepositoryWithOwner{},
	}

	if user != nil {
		userRepos, err := c.repos.FindAllByUser(user.ID)
		if err == nil && userRepos != nil {
			for _, repo := range userRepos {
				data.Repositories = append(data.Repositories, pages.RepositoryWithOwner{
					Repository:    repo,
					OwnerUsername: user.Username,
				})
			}
		}
	}

	return pages.Home(r, data).Render(w, r)
}

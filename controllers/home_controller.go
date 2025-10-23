package controllers

import (
	"net/http"

	"github.com/hypercodehq/hypercode/database/repositories"
	"github.com/hypercodehq/hypercode/middleware"
	"github.com/hypercodehq/hypercode/views/pages"
)

type HomeController interface {
	Show(w http.ResponseWriter, r *http.Request) error
}

type homeController struct {
	repos repositories.RepositoriesRepository
	users repositories.UsersRepository
	orgs  repositories.OrganizationsRepository
	stars repositories.StarsRepository
}

func NewHomeController(
	repos repositories.RepositoriesRepository,
	users repositories.UsersRepository,
	orgs repositories.OrganizationsRepository,
	stars repositories.StarsRepository,
) HomeController {
	return &homeController{
		repos: repos,
		users: users,
		orgs:  orgs,
		stars: stars,
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
				starCount, err := c.stars.CountByRepository(repo.ID)
				if err != nil {
					starCount = 0
				}
				data.Repositories = append(data.Repositories, pages.RepositoryWithOwner{
					Repository:    repo,
					OwnerUsername: user.Username,
					StarCount:     starCount,
				})
			}
		}
	}

	return pages.Home(r, data).Render(w, r)
}

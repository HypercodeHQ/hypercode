package pages

import (
	"net/http"

	html "github.com/hypercodehq/libhtml"
	"github.com/hypercommithq/hypercommit/database/models"
	"github.com/hypercommithq/hypercommit/views/components/layouts"
	"github.com/hypercommithq/hypercommit/views/components/ui"
	"github.com/hypercommithq/libhtml/attr"
)

type ExploreRepositoriesData struct {
	User         *models.User
	Repositories []RepositoryWithOwner
}

type ExploreUsersData struct {
	User  *models.User
	Users []*models.User
}

type ExploreOrganizationsData struct {
	User          *models.User
	Organizations []*models.Organization
}

func ExploreRepositories(r *http.Request, data *ExploreRepositoriesData) html.Node {
	if data == nil {
		data = &ExploreRepositoriesData{}
	}

	var content html.Node
	if len(data.Repositories) > 0 {
		repoCards := make([]html.Node, len(data.Repositories))
		for i, repo := range data.Repositories {
			repoCards[i] = ui.RepositoryCard(ui.RepositoryCardProps{
				OwnerUsername: repo.OwnerUsername,
				Name:          repo.Repository.Name,
				IsPublic:      repo.Repository.Visibility == "public",
				StarCount:     repo.StarCount,
			})
		}

		content = html.Div(
			attr.Class("grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4"),
			html.For(repoCards, func(card html.Node) html.Node {
				return card
			}),
		)
	} else {
		content = ui.EmptyState(ui.EmptyStateProps{
			Icon:        ui.SVGIcon(ui.IconRepository, "size-6"),
			Title:       "No repositories yet",
			Description: "There are no repositories to explore yet.",
			ShowAction:  false,
		})
	}

	return layouts.Explore(r,
		"Explore Repositories - Hypercommit",
		layouts.ExploreLayoutOptions{
			CurrentTab: "repositories",
		},
		html.Main(
			attr.Class("w-full mx-auto max-w-7xl px-4 py-8"),
			html.Div(
				attr.Class("space-y-6"),
				html.H1(
					attr.Class("text-2xl font-semibold"),
					html.Text("Explore repositories"),
				),
				content,
			),
		),
	)
}

func ExploreUsers(r *http.Request, data *ExploreUsersData) html.Node {
	if data == nil {
		data = &ExploreUsersData{}
	}

	var content html.Node
	if len(data.Users) > 0 {
		userCards := make([]html.Node, len(data.Users))
		for i, user := range data.Users {
			userCards[i] = ui.UserCard(ui.UserCardProps{
				Username:    user.Username,
				DisplayName: user.DisplayName,
			})
		}

		content = html.Div(
			attr.Class("grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4"),
			html.For(userCards, func(card html.Node) html.Node {
				return card
			}),
		)
	} else {
		content = ui.EmptyState(ui.EmptyStateProps{
			Icon:        ui.SVGIcon(ui.IconUser, "size-6"),
			Title:       "No users yet",
			Description: "There are no users to explore yet.",
			ShowAction:  false,
		})
	}

	return layouts.Explore(r,
		"Explore Users - Hypercommit",
		layouts.ExploreLayoutOptions{
			CurrentTab: "users",
		},
		html.Main(
			attr.Class("w-full mx-auto max-w-7xl px-4 py-8"),
			html.Div(
				attr.Class("space-y-6"),
				html.H1(
					attr.Class("text-2xl font-semibold"),
					html.Text("Explore users"),
				),
				content,
			),
		),
	)
}

func ExploreOrganizations(r *http.Request, data *ExploreOrganizationsData) html.Node {
	if data == nil {
		data = &ExploreOrganizationsData{}
	}

	var content html.Node
	if len(data.Organizations) > 0 {
		orgCards := make([]html.Node, len(data.Organizations))
		for i, org := range data.Organizations {
			orgCards[i] = ui.OrganizationCard(ui.OrganizationCardProps{
				Username:    org.Username,
				DisplayName: org.DisplayName,
			})
		}

		content = html.Div(
			attr.Class("grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4"),
			html.For(orgCards, func(card html.Node) html.Node {
				return card
			}),
		)
	} else {
		content = ui.EmptyState(ui.EmptyStateProps{
			Icon:        ui.SVGIcon(ui.IconUsers, "size-6"),
			Title:       "No organizations yet",
			Description: "There are no organizations to explore yet.",
			ShowAction:  false,
		})
	}

	return layouts.Explore(r,
		"Explore Organizations - Hypercommit",
		layouts.ExploreLayoutOptions{
			CurrentTab: "organizations",
		},
		html.Main(
			attr.Class("w-full mx-auto max-w-7xl px-4 py-8"),
			html.Div(
				attr.Class("space-y-6"),
				html.H1(
					attr.Class("text-2xl font-semibold"),
					html.Text("Explore organizations"),
				),
				content,
			),
		),
	)
}

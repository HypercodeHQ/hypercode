package pages

import (
	"net/http"

	html "github.com/hypercodehq/libhtml"
	"github.com/hypercommithq/hypercommit/database/models"
	"github.com/hypercommithq/hypercommit/views/components/layouts"
	"github.com/hypercommithq/hypercommit/views/components/ui"
	"github.com/hypercommithq/libhtml/attr"
)

type HomeData struct {
	User         *models.User
	Repositories []RepositoryWithOwner
}

func Home(r *http.Request, data *HomeData) html.Node {
	if data == nil {
		data = &HomeData{}
	}

	var mainContent html.Node

	if data.User != nil {
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

			mainContent = html.Div(
				attr.Class("space-y-6"),
				html.H1(
					attr.Class("text-xl font-medium"),
					html.Text("My repositories"),
				),
				html.Div(
					attr.Class("grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4"),
					html.For(repoCards, func(card html.Node) html.Node {
						return card
					}),
				),
			)
		} else {
			mainContent = ui.EmptyState(ui.EmptyStateProps{
				Icon:        ui.SVGIcon(ui.IconGitBranch, "size-6"),
				Title:       "No repositories yet",
				Description: "Get started by creating your first repository.",
				ActionText:  "Create repository",
				ActionHref:  "/repositories/new",
				ShowAction:  true,
			})
		}
	} else {
		mainContent = ui.EmptyState(ui.EmptyStateProps{
			Icon:        ui.SVGIcon(ui.IconGitBranch, "size-6"),
			Title:       "No repositories yet",
			Description: "Sign up to create and manage your repositories.",
			ActionText:  "Create repository",
			ActionHref:  "/auth/sign-up",
			ShowAction:  true,
		})
	}

	return layouts.Main(r,
		"Hypercommit: An open-source alternative to GitHub",
		html.Main(
			attr.Class("w-full mx-auto max-w-7xl px-4 py-8"),
			mainContent,
		),
	)
}

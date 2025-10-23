package pages

import (
	"net/http"

	"github.com/hyperstitieux/hypercode/database/models"
	"github.com/hyperstitieux/hypercode/html"
	"github.com/hyperstitieux/hypercode/html/attr"
	"github.com/hyperstitieux/hypercode/views/components/layouts"
	"github.com/hyperstitieux/hypercode/views/components/ui"
)

type UserProfileData struct {
	User         *models.User
	ProfileUser  *models.User
	Repositories []*models.Repository
}

func UserProfile(r *http.Request, data *UserProfileData) html.Node {
	if data == nil || data.ProfileUser == nil {
		data = &UserProfileData{}
	}

	var mainContent html.Node

	if len(data.Repositories) > 0 {
		repoCards := make([]html.Node, len(data.Repositories))
		for i, repo := range data.Repositories {
			repoCards[i] = ui.RepositoryCard(ui.RepositoryCardProps{
				OwnerUsername: data.ProfileUser.Username,
				Name:          repo.Name,
				IsPublic:      repo.Visibility == "public",
			})
		}

		mainContent = html.Div(
			attr.Class("space-y-6"),
			// User header
			html.Div(
				attr.Class("flex items-center gap-4 pb-6 border-b"),
				html.Div(
					attr.Class("flex flex-col"),
					html.H1(
						attr.Class("text-2xl font-semibold"),
						html.Text(data.ProfileUser.DisplayName),
					),
					html.P(
						attr.Class("text-muted-foreground"),
						html.Text("@"+data.ProfileUser.Username),
					),
				),
			),
			// Repositories section
			html.Div(
				attr.Class("space-y-4"),
				html.H2(
					attr.Class("text-xl font-medium"),
					html.Text("Repositories"),
				),
				html.Div(
					attr.Class("grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4"),
					html.For(repoCards, func(card html.Node) html.Node {
						return card
					}),
				),
			),
		)
	} else {
		mainContent = html.Div(
			attr.Class("space-y-6"),
			// User header
			html.Div(
				attr.Class("flex items-center gap-4 pb-6 border-b"),
				html.Div(
					attr.Class("flex flex-col"),
					html.H1(
						attr.Class("text-2xl font-semibold"),
						html.Text(data.ProfileUser.DisplayName),
					),
					html.P(
						attr.Class("text-muted-foreground"),
						html.Text("@"+data.ProfileUser.Username),
					),
				),
			),
			// Empty state
			ui.EmptyState(ui.EmptyStateProps{
				Icon:        ui.SVGIcon(ui.IconGitBranch, "size-6"),
				Title:       "No repositories yet",
				Description: "This user hasn't created any repositories yet.",
				ShowAction:  false,
			}),
		)
	}

	return layouts.Main(r,
		data.ProfileUser.DisplayName+" - Hypercode",
		html.Main(
			attr.Class("w-full mx-auto max-w-6xl px-4 py-8"),
			mainContent,
		),
	)
}

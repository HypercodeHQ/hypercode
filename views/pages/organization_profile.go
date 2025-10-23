package pages

import (
	"net/http"

	"github.com/hyperstitieux/hypercode/database/models"
	"github.com/hyperstitieux/hypercode/html"
	"github.com/hyperstitieux/hypercode/html/attr"
	"github.com/hyperstitieux/hypercode/views/components/layouts"
	"github.com/hyperstitieux/hypercode/views/components/ui"
)

type OrganizationProfileData struct {
	User         *models.User
	Organization *models.Organization
	Repositories []*models.Repository
	StarCounts   map[int64]int64
	CanManage    bool
	CurrentTab   string
}

func OrganizationProfile(r *http.Request, data *OrganizationProfileData) html.Node {
	if data == nil || data.Organization == nil {
		data = &OrganizationProfileData{}
	}

	// Default to overview tab if not specified
	if data.CurrentTab == "" {
		data.CurrentTab = "overview"
	}

	var tabContent html.Node

	switch data.CurrentTab {
	case "repositories":
		tabContent = renderRepositoriesTab(data)
	case "stars":
		tabContent = renderStarsTab(data)
	default: // overview
		tabContent = renderOverviewTab(data)
	}

	return layouts.Profile(r,
		data.Organization.DisplayName+" - Hypercode",
		layouts.ProfileLayoutOptions{
			Username:     data.Organization.Username,
			DisplayName:  data.Organization.DisplayName,
			IsOrg:        true,
			CurrentTab:   data.CurrentTab,
			ShowSettings: data.CanManage,
		},
		html.Main(
			attr.Class("w-full mx-auto max-w-7xl px-4 py-8"),
			html.Div(
				attr.Class("space-y-6"),
				// Organization header with display name
				html.Div(
					attr.Class("flex items-center gap-4 pb-6 border-b"),
					html.Div(
						attr.Class("flex items-center gap-3"),
						html.Div(
							attr.Class("p-3 rounded-full bg-muted"),
							ui.SVGIcon(ui.IconUsers, "size-8"),
						),
						html.Div(
							attr.Class("flex flex-col"),
							html.H1(
								attr.Class("text-2xl font-semibold"),
								html.Text(data.Organization.DisplayName),
							),
							html.P(
								attr.Class("text-muted-foreground"),
								html.Text("@"+data.Organization.Username),
							),
						),
					),
				),
				// Tab content
				tabContent,
			),
		),
	)
}

func renderOverviewTab(data *OrganizationProfileData) html.Node {
	if len(data.Repositories) == 0 {
		return ui.EmptyState(ui.EmptyStateProps{
			Icon:        ui.SVGIcon(ui.IconRepository, "size-6"),
			Title:       "No repositories yet",
			Description: "This organization hasn't created any repositories yet.",
			ShowAction:  false,
		})
	}

	// Show top 6 repositories for overview
	repoCount := len(data.Repositories)
	if repoCount > 6 {
		repoCount = 6
	}

	repoCards := make([]html.Node, repoCount)
	for i := 0; i < repoCount; i++ {
		repo := data.Repositories[i]
		starCount := int64(0)
		if data.StarCounts != nil {
			starCount = data.StarCounts[repo.ID]
		}
		repoCards[i] = ui.RepositoryCard(ui.RepositoryCardProps{
			OwnerUsername: data.Organization.Username,
			Name:          repo.Name,
			IsPublic:      repo.Visibility == "public",
			StarCount:     starCount,
		})
	}

	return html.Div(
		attr.Class("space-y-4"),
		html.Div(
			attr.Class("flex justify-between items-center"),
			html.H2(
				attr.Class("text-xl font-medium"),
				html.Text("Popular repositories"),
			),
			html.If(
				len(data.Repositories) > 6,
				html.A(
					attr.Href("/"+data.Organization.Username+"/repositories"),
					attr.Class("text-sm text-primary hover:underline"),
					html.Text("View all repositories"),
				),
			),
		),
		html.Div(
			attr.Class("grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4"),
			html.For(repoCards, func(card html.Node) html.Node {
				return card
			}),
		),
	)
}

func renderRepositoriesTab(data *OrganizationProfileData) html.Node {
	if len(data.Repositories) == 0 {
		return ui.EmptyState(ui.EmptyStateProps{
			Icon:        ui.SVGIcon(ui.IconRepository, "size-6"),
			Title:       "No repositories yet",
			Description: "This organization hasn't created any repositories yet.",
			ShowAction:  false,
		})
	}

	repoCards := make([]html.Node, len(data.Repositories))
	for i, repo := range data.Repositories {
		starCount := int64(0)
		if data.StarCounts != nil {
			starCount = data.StarCounts[repo.ID]
		}
		repoCards[i] = ui.RepositoryCard(ui.RepositoryCardProps{
			OwnerUsername: data.Organization.Username,
			Name:          repo.Name,
			IsPublic:      repo.Visibility == "public",
			StarCount:     starCount,
		})
	}

	return html.Div(
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
	)
}

func renderStarsTab(data *OrganizationProfileData) html.Node {
	// TODO: Implement starred repositories for organizations
	return ui.EmptyState(ui.EmptyStateProps{
		Icon:        ui.SVGIcon(ui.IconStar, "size-6"),
		Title:       "No starred repositories yet",
		Description: "This organization hasn't starred any repositories yet.",
		ShowAction:  false,
	})
}

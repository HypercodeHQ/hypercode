package components

import (
	"github.com/hyperstitieux/hypercode/database/models"
	"github.com/hyperstitieux/hypercode/html"
	"github.com/hyperstitieux/hypercode/html/attr"
	"github.com/hyperstitieux/hypercode/views/components/ui"
)

type RepositoryHeaderData struct {
	User          *models.User
	OwnerUsername string
	RepoName      string
	IsPublic      bool
}

func RepositoryHeader(data *RepositoryHeaderData) html.Node {
	if data == nil {
		data = &RepositoryHeaderData{}
	}

	return html.Header(
		attr.Class("bg-background border-b px-4 py-3 flex flex-wrap justify-between items-center gap-2"),
		html.Div(
			attr.Class("flex flex-wrap items-center gap-3"),
			html.A(
				attr.Href("/"),
				attr.DataTooltip("Go back home"),
				attr.DataSide("bottom"),
				html.Img(
					attr.Src("/logo.png"),
					attr.Alt("Hypercode"),
					attr.Class("h-7"),
				),
			),
			html.Element("span",
				attr.Class("text-muted-foreground text-lg"),
				html.Text("/"),
			),
			html.A(
				attr.Href("/"+data.OwnerUsername),
				attr.Class("text-lg font-medium hover:underline"),
				html.Text(data.OwnerUsername),
			),
			html.Element("span",
				attr.Class("text-muted-foreground text-lg"),
				html.Text("/"),
			),
			html.A(
				attr.Href("/"+data.OwnerUsername+"/"+data.RepoName),
				attr.Class("text-lg font-medium hover:underline"),
				html.Text(data.RepoName),
			),
			ui.Badge(
				ui.BadgeProps{
					Variant: ui.BadgeOutline,
					Class:   "bg-card",
				},
				html.Text(visibilityText(data.IsPublic)),
			),
		),
		html.Div(
			attr.Class("flex flex-wrap items-center gap-4"),
			html.IfElsef(
				data.User != nil,
				func() html.Node { return repositoryLoggedInActions(data.User) },
				func() html.Node { return loggedOutActions() },
			),
		),
	)
}

func repositoryLoggedInActions(user *models.User) html.Node {
	return html.Div(
		attr.Class("flex flex-wrap items-center gap-4"),
		createNewDropdown(),
		userAccountDropdown(user),
	)
}

func visibilityText(isPublic bool) string {
	if isPublic {
		return "Public"
	}
	return "Private"
}

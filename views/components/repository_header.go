package components

import (
	"fmt"

	"github.com/hypercodehq/hypercode/database/models"
	"github.com/hypercodehq/libhtml"
	"github.com/hypercodehq/libhtml/attr"
	"github.com/hypercodehq/hypercode/views/components/ui"
)

type RepositoryHeaderData struct {
	User          *models.User
	OwnerUsername string
	RepoName      string
	IsPublic      bool
	CurrentTab    string
	ShowSettings  bool
	StarCount     int64
	HasStarred    bool
	DefaultBranch string
	CloneURL      string
	RepositoryURL string
}

func RepositoryHeader(data *RepositoryHeaderData) html.Node {
	if data == nil {
		data = &RepositoryHeaderData{}
	}

	return html.Group(
		html.Header(
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
					attr.Class("text-muted-foreground text-md"),
					html.Text("/"),
				),
				html.A(
					attr.Href("/"+data.OwnerUsername),
					attr.Class("text-md font-medium hover:underline"),
					html.Text(data.OwnerUsername),
				),
				html.Element("span",
					attr.Class("text-muted-foreground text-md"),
					html.Text("/"),
				),
				html.A(
					attr.Href("/"+data.OwnerUsername+"/"+data.RepoName),
					attr.Class("text-md font-medium hover:underline"),
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
		),
		html.Div(
			attr.Class("bg-background border-b px-4 flex flex-wrap justify-between items-center gap-4"),
			ui.RepositoryTabs(ui.RepositoryTabsProps{
				OwnerUsername: data.OwnerUsername,
				RepoName:      data.RepoName,
				CurrentTab:    data.CurrentTab,
				ShowSettings:  data.ShowSettings,
				DefaultBranch: data.DefaultBranch,
			}),
			html.If(
				data.User != nil,
				html.Div(
					attr.Class("flex items-center gap-2"),
					ShareDropdown(&RepositoryActionsDropdownData{
						OwnerUsername: data.OwnerUsername,
						RepoName:      data.RepoName,
						CloneURL:      data.CloneURL,
						RepositoryURL: data.RepositoryURL,
					}),
					CloneDropdown(&RepositoryActionsDropdownData{
						OwnerUsername: data.OwnerUsername,
						RepoName:      data.RepoName,
						CloneURL:      data.CloneURL,
						RepositoryURL: data.RepositoryURL,
					}),
					starButton(data),
				),
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

func starButton(data *RepositoryHeaderData) html.Node {
	if data.HasStarred {
		return html.Element("form",
			attr.Method("post"),
			attr.Action("/"+data.OwnerUsername+"/"+data.RepoName+"/unstar"),
			html.Element("button",
				attr.Type("submit"),
				attr.Class("btn-outline inline-flex items-center gap-2"),
				starIconFilled(),
				html.Text("Unstar "),
				html.Element("span",
					attr.Class("bg-muted px-1.5 py-0.5 rounded text-sm"),
					html.Text(formatStarCount(data.StarCount)),
				),
			),
		)
	}
	return html.Element("form",
		attr.Method("post"),
		attr.Action("/"+data.OwnerUsername+"/"+data.RepoName+"/star"),
		html.Element("button",
			attr.Type("submit"),
			attr.Class("btn-outline inline-flex items-center gap-2"),
			ui.SVGIcon(ui.IconStar, "size-4"),
			html.Text("Star "),
			html.Element("span",
				attr.Class("bg-muted px-1.5 py-0.5 rounded text-sm"),
				html.Text(formatStarCount(data.StarCount)),
			),
		),
	)
}

func formatStarCount(count int64) string {
	return fmt.Sprintf("%d", count)
}

func starIconFilled() html.Node {
	return html.Element("svg",
		attr.Xmlns("http://www.w3.org/2000/svg"),
		attr.Width("16"),
		attr.Height("16"),
		attr.ViewBox("0 0 24 24"),
		attr.Fill("#facc15"),
		attr.Stroke("#facc15"),
		attr.StrokeWidth("2"),
		attr.StrokeLinecap("round"),
		attr.StrokeLinejoin("round"),
		attr.Class("size-4"),
		html.Element("polygon", attr.Points("12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2")),
	)
}

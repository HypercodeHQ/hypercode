package components

import (
	"github.com/hypercommithq/hypercommit/database/models"
	"github.com/hypercommithq/libhtml"
	"github.com/hypercommithq/libhtml/attr"
	"github.com/hypercommithq/hypercommit/views/components/ui"
)

type ProfileHeaderData struct {
	User         *models.User
	Username     string
	DisplayName  string
	IsOrg        bool
	CurrentTab   string
	ShowSettings bool
}

func ProfileHeader(data *ProfileHeaderData) html.Node {
	if data == nil {
		data = &ProfileHeaderData{}
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
				html.Span(
					attr.Class("text-muted-foreground text-md"),
					html.Text("/"),
				),
				html.A(
					attr.Href("/"+data.Username),
					attr.Class("text-md font-medium hover:underline"),
					html.Text(data.Username),
				),
			),
			html.Div(
				attr.Class("flex flex-wrap items-center gap-4"),
				html.IfElsef(
					data.User != nil,
					func() html.Node { return profileLoggedInActions(data.User) },
					func() html.Node { return loggedOutActions() },
				),
			),
		),
		html.Div(
			attr.Class("bg-background border-b px-4"),
			html.Div(
				attr.Class("flex flex-wrap justify-between items-center gap-4"),
				ui.ProfileTabs(ui.ProfileTabsProps{
					Username:     data.Username,
					CurrentTab:   data.CurrentTab,
					ShowSettings: data.ShowSettings,
					IsOrg:        data.IsOrg,
				}),
			),
		),
	)
}

func profileLoggedInActions(user *models.User) html.Node {
	return html.Div(
		attr.Class("flex flex-wrap items-center gap-4"),
		createNewDropdown(),
		userAccountDropdown(user),
	)
}

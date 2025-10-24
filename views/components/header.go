package components

import (
	"strings"

	"github.com/hypercommithq/hypercommit/database/models"
	"github.com/hypercommithq/hypercommit/views/components/ui"
	html "github.com/hypercommithq/libhtml"
	"github.com/hypercommithq/libhtml/attr"
)

type HeaderData struct {
	User *models.User
}

func Header(data *HeaderData, children ...html.Node) html.Node {
	if data == nil {
		data = &HeaderData{}
	}

	return html.Header(
		attr.Class("bg-background border-b px-4 py-3 flex flex-wrap justify-between items-center gap-2"),
		html.Div(
			attr.Class("flex flex-wrap items-center gap-4"),
			html.A(
				attr.Href("/"),
				attr.DataTooltip("Go back home"),
				attr.DataSide("bottom"),
				html.Img(
					attr.Src("/logo.png"),
					attr.Alt("Hypercommit"),
					attr.Class("h-7"),
				),
			),
			html.Group(children...),
		),
		html.Div(
			attr.Class("flex flex-wrap items-center gap-4"),
			html.IfElsef(
				data.User != nil,
				func() html.Node { return loggedInActions(data.User) },
				func() html.Node { return loggedOutActions() },
			),
		),
	)
}

func loggedInActions(user *models.User) html.Node {
	return html.Div(
		attr.Class("flex flex-wrap items-center gap-4"),
		createNewDropdown(),
		userAccountDropdown(user),
	)
}

func socialLinks() html.Node {
	return html.Div(
		attr.Class("flex flex-wrap items-center gap-2"),
		html.A(
			attr.Href("https://x.com/hypercommit2099"),
			attr.Target("_blank"),
			attr.Rel("noopener noreferrer"),
			attr.AriaLabel("Follow us on X"),
			attr.DataTooltip("Follow us on X"),
			attr.DataSide("bottom"),
			attr.Class("btn-icon-ghost"),
			ui.SVGIcon(ui.IconTwitter, ""),
		),
		html.A(
			attr.Href("https://discord.gg/9ntvVASRKf"),
			attr.Target("_blank"),
			attr.Rel("noopener noreferrer"),
			attr.AriaLabel("Join our Discord"),
			attr.DataTooltip("Join our Discord"),
			attr.DataSide("bottom"),
			attr.Class("btn-icon-ghost"),
			ui.SVGIcon(ui.IconDiscord, ""),
		),
		html.A(
			attr.Href("https://bsky.app/profile/hypercommit.com"),
			attr.Target("_blank"),
			attr.Rel("noopener noreferrer"),
			attr.AriaLabel("Follow us on Bluesky"),
			attr.DataTooltip("Follow us on Bluesky"),
			attr.DataSide("bottom"),
			attr.Class("btn-icon-ghost"),
			ui.SVGIcon(ui.IconBluesky, ""),
		),
		html.A(
			attr.Href("https://github.com/hypercommithq/hypercommit"),
			attr.Target("_blank"),
			attr.Rel("noopener noreferrer"),
			attr.AriaLabel("Star us on GitHub"),
			attr.DataTooltip("Star us on GitHub"),
			attr.DataSide("bottom"),
			attr.Class("btn-icon-ghost"),
			ui.SVGIcon(ui.IconGitHub, ""),
		),
	)
}

func loggedOutActions() html.Node {
	return html.Div(
		attr.Class("flex flex-wrap items-center gap-4"),
		socialLinks(),
		html.A(
			attr.Href("/auth/sign-in"),
			attr.Class("btn-outline"),
			html.Text("Sign in"),
		),
		html.A(
			attr.Href("/auth/sign-up"),
			attr.Class("btn"),
			html.Text("Create an account"),
		),
	)
}

func createNewDropdown() html.Node {
	return html.Div(
		attr.Class("dropdown-menu"),
		html.Element("button",
			attr.Type("button"),
			attr.Id("create-new-trigger"),
			attr.Class("btn-icon-ghost"),
			attr.AriaHaspopup("menu"),
			attr.AriaControls("create-new-menu"),
			attr.AriaExpanded("false"),
			ui.SVGIcon(ui.IconPlus, ""),
		),
		html.Div(
			attr.Id("create-new-popover"),
			attr.DataPopover(""),
			attr.AriaHidden("true"),
			attr.Class("min-w-56 right-0 left-auto"),
			html.Div(
				attr.Role("menu"),
				attr.Id("create-new-menu"),
				attr.AriaLabelledby("create-new-trigger"),
				html.Div(
					attr.Role("group"),
					attr.AriaLabelledby("create-new-label"),
					html.Div(
						attr.Role("heading"),
						attr.Id("create-new-label"),
						html.Text("Create"),
					),
					html.Hr(attr.Role("separator")),
					html.A(
						attr.Class("cursor-pointer"),
						attr.Href("/repositories/new"),
						attr.Role("menuitem"),
						ui.SVGIcon(ui.IconRepository, ""),
						html.Text("New repository"),
					),
					html.Hr(attr.Role("separator")),
					html.A(
						attr.Class("cursor-pointer"),
						attr.Href("/organizations/new"),
						attr.Role("menuitem"),
						ui.SVGIcon(ui.IconBuilding, ""),
						html.Text("New organization"),
					),
				),
			),
		),
	)
}

func userAccountDropdown(user *models.User) html.Node {
	return html.Div(
		attr.Class("dropdown-menu"),
		html.Element("button",
			attr.Type("button"),
			attr.Id("user-account-trigger"),
			attr.Class("btn-icon-ghost size-9"),
			attr.AriaHaspopup("menu"),
			attr.AriaControls("user-account-menu"),
			attr.AriaExpanded("false"),
			html.Div(
				attr.Class("size-8 rounded-full bg-primary text-primary-foreground flex items-center justify-center text-xs font-medium"),
				html.Text(getUserInitials(user)),
			),
		),
		html.Div(
			attr.Id("user-account-popover"),
			attr.DataPopover(""),
			attr.AriaHidden("true"),
			attr.Class("min-w-56 right-0 left-auto"),
			html.Div(
				attr.Role("menu"),
				attr.Id("user-account-menu"),
				attr.AriaLabelledby("user-account-trigger"),
				html.Div(
					attr.Role("group"),
					attr.AriaLabelledby("user-label"),
					html.Div(
						attr.Role("heading"),
						attr.Id("user-label"),
						html.Text("@"+user.Username),
					),
					html.Hr(attr.Role("separator")),
					html.A(
						attr.Href("/"+user.Username),
						attr.Role("menuitem"),
						attr.Class("cursor-pointer"),
						ui.SVGIcon(ui.IconUser, ""),
						html.Text("Profile"),
					),
					html.A(
						attr.Href("/settings"),
						attr.Role("menuitem"),
						attr.Class("cursor-pointer"),
						ui.SVGIcon(ui.IconSettings, ""),
						html.Text("Settings"),
					),
					html.Hr(attr.Role("separator")),
					html.A(
						attr.Href("https://x.com/hypercommit2099"),
						attr.Target("_blank"),
						attr.Rel("noopener noreferrer"),
						attr.Role("menuitem"),
						attr.Class("cursor-pointer"),
						ui.SVGIcon(ui.IconTwitter, ""),
						html.Text("Follow us on X"),
					),
					html.A(
						attr.Href("https://discord.gg/9ntvVASRKf"),
						attr.Target("_blank"),
						attr.Rel("noopener noreferrer"),
						attr.Role("menuitem"),
						attr.Class("cursor-pointer"),
						ui.SVGIcon(ui.IconDiscord, ""),
						html.Text("Join our Discord"),
					),
					html.A(
						attr.Href("https://bsky.app/profile/hypercommit.com"),
						attr.Target("_blank"),
						attr.Rel("noopener noreferrer"),
						attr.Role("menuitem"),
						attr.Class("cursor-pointer"),
						ui.SVGIcon(ui.IconBluesky, ""),
						html.Text("Follow us on Bluesky"),
					),
					html.A(
						attr.Href("https://github.com/hypercommithq/hypercommit"),
						attr.Target("_blank"),
						attr.Rel("noopener noreferrer"),
						attr.Role("menuitem"),
						attr.Class("cursor-pointer"),
						ui.SVGIcon(ui.IconGitHub, ""),
						html.Text("Star us on GitHub"),
					),
					html.Hr(attr.Role("separator")),
					html.A(
						attr.Href("/auth/sign-out"),
						attr.Role("menuitem"),
						attr.Class("cursor-pointer text-destructive"),
						ui.SVGIcon(ui.IconLogOut, "text-destructive"),
						html.Text("Sign out"),
					),
				),
			),
		),
	)
}

func getUserInitials(user *models.User) string {
	if user == nil || user.Username == "" {
		return "?"
	}
	if len(user.Username) >= 2 {
		return strings.ToUpper(user.Username[0:2])
	}
	return strings.ToUpper(string(user.Username[0]))
}

type MainHeaderData struct {
	User  *models.User
	Class string
}

func MainHeader(data *MainHeaderData) html.Node {
	return Header(&HeaderData{User: data.User}, html.A(
		attr.Class("btn-ghost "+data.Class),
		attr.Href("/explore/repositories"),
		html.Text("Explore"),
	))
}

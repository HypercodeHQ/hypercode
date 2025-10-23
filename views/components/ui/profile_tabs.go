package ui

import (
	"github.com/hyperstitieux/hypercode/html"
	"github.com/hyperstitieux/hypercode/html/attr"
)

type ProfileTabsProps struct {
	Username     string
	CurrentTab   string
	ShowSettings bool
	IsOrg        bool
}

func ProfileTabs(props ProfileTabsProps) html.Node {
	tabs := []html.Node{
		profileTab(
			props.Username,
			"overview",
			props.CurrentTab,
			IconLayoutGrid,
			"Overview",
		),
		profileTab(
			props.Username,
			"repositories",
			props.CurrentTab,
			IconRepository,
			"Repositories",
		),
	}

	// Only show stars tab for users, not organizations
	if !props.IsOrg {
		tabs = append(tabs, profileTab(
			props.Username,
			"stars",
			props.CurrentTab,
			IconStar,
			"Stars",
		))
	}

	if props.ShowSettings {
		tabs = append(tabs, profileTab(
			props.Username,
			"settings",
			props.CurrentTab,
			IconSettings,
			"Settings",
		))
	}

	return html.Nav(
		attr.Class("flex flex-wrap items-center gap-4"),
		html.Group(tabs...),
	)
}

func profileTab(username, tab, currentTab string, icon Icon, label string) html.Node {
	href := "/" + username
	if tab != "overview" {
		href += "/" + tab
	}

	isActive := tab == currentTab
	spanClasses := "btn-ghost inline-flex items-center gap-2"
	if isActive {
		spanClasses += " font-medium"
	} else {
		spanClasses += " text-muted-foreground"
	}

	borderClass := "border-transparent"
	if isActive {
		borderClass = "border-zinc-900"
	}

	return html.A(
		attr.Href(href),
		attr.Class("inline-flex mt-2 pb-2 border-b-2 transition-colors "+borderClass),
		html.Element("span",
			attr.Class(spanClasses),
			smallSVGIcon(icon, ""),
			html.Text(label),
		),
	)
}

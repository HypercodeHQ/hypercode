package ui

import (
	"github.com/hypercommithq/libhtml"
	"github.com/hypercommithq/libhtml/attr"
)

type ExploreTabsProps struct {
	CurrentTab string
}

func ExploreTabs(props ExploreTabsProps) html.Node {
	tabs := []html.Node{
		exploreTab(
			"repositories",
			props.CurrentTab,
			IconRepository,
			"Repositories",
		),
		exploreTab(
			"users",
			props.CurrentTab,
			IconUser,
			"Users",
		),
		exploreTab(
			"organizations",
			props.CurrentTab,
			IconUsers,
			"Organizations",
		),
	}

	return html.Nav(
		attr.Class("flex flex-wrap items-center gap-4"),
		html.Group(tabs...),
	)
}

func exploreTab(tab, currentTab string, icon Icon, label string) html.Node {
	href := "/explore/" + tab
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
		attr.Class("inline-flex pb-2 border-b-2 transition-colors "+borderClass),
		html.Element("span",
			attr.Class(spanClasses),
			smallSVGIcon(icon, ""),
			html.Text(label),
		),
	)
}

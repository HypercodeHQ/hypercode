package ui

import (
	"github.com/hypercodehq/hypercode/html"
	"github.com/hypercodehq/hypercode/html/attr"
)

type RepositoryTabsProps struct {
	OwnerUsername string
	RepoName      string
	CurrentTab    string
	ShowSettings  bool
	DefaultBranch string
}

func RepositoryTabs(props RepositoryTabsProps) html.Node {
	tabs := []html.Node{
		repositoryTab(
			props.OwnerUsername,
			props.RepoName,
			props.DefaultBranch,
			"overview",
			props.CurrentTab,
			IconLayoutGrid,
			"Overview",
		),
		repositoryTab(
			props.OwnerUsername,
			props.RepoName,
			props.DefaultBranch,
			"tree",
			props.CurrentTab,
			IconCode,
			"Code",
		),
		repositoryTab(
			props.OwnerUsername,
			props.RepoName,
			props.DefaultBranch,
			"tickets",
			props.CurrentTab,
			IconCircle,
			"Tickets",
		),
	}

	if props.ShowSettings {
		tabs = append(tabs, repositoryTab(
			props.OwnerUsername,
			props.RepoName,
			props.DefaultBranch,
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

func repositoryTab(ownerUsername, repoName, defaultBranch, tab, currentTab string, icon Icon, label string) html.Node {
	href := "/" + ownerUsername + "/" + repoName
	if tab == "tree" {
		// For tree tab, link to /tree/<defaultBranch>
		href += "/tree/" + defaultBranch
	} else if tab != "overview" {
		href += "/" + tab
	}

	isActive := tab == currentTab
	spanClasses := "btn-ghost inline-flex items-center gap-2"
	if isActive {
		spanClasses += " font-medium"
	} else {
		spanClasses += " text-muted-foreground"
	}

	return html.A(
		attr.Href(href),
		attr.Class("inline-flex mt-2 pb-2 border-b-2 transition-colors "+borderClass(isActive)),
		html.Element("span",
			attr.Class(spanClasses),
			smallSVGIcon(icon, ""),
			html.Text(label),
		),
	)
}

func borderClass(isActive bool) string {
	if isActive {
		return "border-zinc-900"
	}
	return "border-transparent"
}

func smallSVGIcon(icon Icon, class string) html.Node {
	svgAttrs := []html.Node{
		attr.Xmlns("http://www.w3.org/2000/svg"),
		attr.Width("16"),
		attr.Height("16"),
		attr.ViewBox("0 0 24 24"),
		attr.Fill("none"),
		attr.Stroke("currentColor"),
		attr.StrokeWidth("2"),
		attr.StrokeLinecap("round"),
		attr.StrokeLinejoin("round"),
	}

	if class != "" {
		svgAttrs = append(svgAttrs, attr.Class(class))
	}

	var paths []html.Node

	switch icon {
	case IconLayoutGrid:
		paths = []html.Node{
			html.Element("rect", attr.Width("7"), attr.Height("7"), attr.X("3"), attr.Y("3"), attr.Rx("1")),
			html.Element("rect", attr.Width("7"), attr.Height("7"), attr.X("14"), attr.Y("3"), attr.Rx("1")),
			html.Element("rect", attr.Width("7"), attr.Height("7"), attr.X("14"), attr.Y("14"), attr.Rx("1")),
			html.Element("rect", attr.Width("7"), attr.Height("7"), attr.X("3"), attr.Y("14"), attr.Rx("1")),
		}
	case IconCode:
		paths = []html.Node{
			html.Element("polyline", attr.Points("16 18 22 12 16 6")),
			html.Element("polyline", attr.Points("8 6 2 12 8 18")),
		}
	case IconRepository:
		paths = []html.Node{
			html.Element("path", attr.D("M4 19.5v-15A2.5 2.5 0 0 1 6.5 2H19a1 1 0 0 1 1 1v18a1 1 0 0 1-1 1H6.5a1 1 0 0 1 0-5H20")),
			html.Element("path", attr.D("M8 7h6")),
			html.Element("path", attr.D("M8 11h8")),
		}
	case IconUser:
		paths = []html.Node{
			html.Element("path", attr.D("M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2")),
			html.Element("circle", attr.Cx("12"), attr.Cy("7"), attr.R("4")),
		}
	case IconUsers:
		paths = []html.Node{
			html.Element("path", attr.D("M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2")),
			html.Element("circle", attr.Cx("9"), attr.Cy("7"), attr.R("4")),
			html.Element("path", attr.D("M22 21v-2a4 4 0 0 0-3-3.87")),
			html.Element("path", attr.D("M16 3.13a4 4 0 0 1 0 7.75")),
		}
	case IconSettings:
		paths = []html.Node{
			html.Element("path", attr.D("M12.22 2h-.44a2 2 0 0 0-2 2v.18a2 2 0 0 1-1 1.73l-.43.25a2 2 0 0 1-2 0l-.15-.08a2 2 0 0 0-2.73.73l-.22.38a2 2 0 0 0 .73 2.73l.15.1a2 2 0 0 1 1 1.72v.51a2 2 0 0 1-1 1.74l-.15.09a2 2 0 0 0-.73 2.73l.22.38a2 2 0 0 0 2.73.73l.15-.08a2 2 0 0 1 2 0l.43.25a2 2 0 0 1 1 1.73V20a2 2 0 0 0 2 2h.44a2 2 0 0 0 2-2v-.18a2 2 0 0 1 1-1.73l.43-.25a2 2 0 0 1 2 0l.15.08a2 2 0 0 0 2.73-.73l.22-.39a2 2 0 0 0-.73-2.73l-.15-.08a2 2 0 0 1-1-1.74v-.5a2 2 0 0 1 1-1.74l.15-.09a2 2 0 0 0 .73-2.73l-.22-.38a2 2 0 0 0-2.73-.73l-.15.08a2 2 0 0 1-2 0l-.43-.25a2 2 0 0 1-1-1.73V4a2 2 0 0 0-2-2z")),
			html.Element("circle", attr.Cx("12"), attr.Cy("12"), attr.R("3")),
		}
	case IconStar:
		paths = []html.Node{
			html.Element("polygon", attr.Points("12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2")),
		}
	case IconCircle:
		paths = []html.Node{
			html.Element("circle", attr.Cx("12"), attr.Cy("12"), attr.R("10")),
			html.Element("circle", attr.Cx("12"), attr.Cy("12"), attr.R("1"), attr.Fill("currentColor")),
		}
	}

	return html.Element("svg", append(svgAttrs, paths...)...)
}

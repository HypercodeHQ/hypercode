package ui

import (
	"github.com/hyperstitieux/hypercode/html"
	"github.com/hyperstitieux/hypercode/html/attr"
)

type Icon string

const (
	IconCheck        Icon = "check"
	IconChevronDown  Icon = "chevron-down"
	IconX            Icon = "x"
	IconAlertCircle  Icon = "alert-circle"
	IconInfo         Icon = "info"
	IconSend         Icon = "send"
	IconArrowRight   Icon = "arrow-right"
	IconLoader       Icon = "loader"
	IconDownload     Icon = "download"
	IconUpload       Icon = "upload"
	IconMoreVertical Icon = "more-vertical"
	IconTrash        Icon = "trash"
	IconPlus         Icon = "plus"
	IconRepository   Icon = "repository"
	IconBuilding     Icon = "building"
	IconUser         Icon = "user"
	IconSettings     Icon = "settings"
	IconLogOut       Icon = "log-out"
	IconMail         Icon = "mail"
	IconAtSign       Icon = "at-sign"
	IconLock         Icon = "lock"
	IconGitBranch    Icon = "git-branch"
	IconLayoutGrid   Icon = "layout-grid"
	IconCode         Icon = "code"
	IconCopy         Icon = "copy"
	IconGlobe        Icon = "globe"
	IconTwitter      Icon = "twitter"
	IconDiscord      Icon = "discord"
	IconGitHub       Icon = "github"
	IconBluesky      Icon = "bluesky"
	IconStar         Icon = "star"
)

func SVGIcon(icon Icon, class string) html.Node {
	svgAttrs := []html.Node{
		attr.Xmlns("http://www.w3.org/2000/svg"),
		attr.Width("24"),
		attr.Height("24"),
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
	case IconCheck:
		paths = []html.Node{
			html.Element("path", attr.D("M20 6 9 17l-5-5")),
		}
	case IconChevronDown:
		paths = []html.Node{
			html.Element("path", attr.D("m6 9 6 6 6-6")),
		}
	case IconX:
		paths = []html.Node{
			html.Element("path", attr.D("M18 6 6 18")),
			html.Element("path", attr.D("m6 6 12 12")),
		}
	case IconAlertCircle:
		paths = []html.Node{
			html.Element("circle", attr.Cx("12"), attr.Cy("12"), attr.R("10")),
			html.Element("line", attr.X1("12"), attr.X2("12"), attr.Y1("8"), attr.Y2("12")),
			html.Element("line", attr.X1("12"), attr.X2("12.01"), attr.Y1("16"), attr.Y2("16")),
		}
	case IconInfo:
		paths = []html.Node{
			html.Element("circle", attr.Cx("12"), attr.Cy("12"), attr.R("10")),
			html.Element("path", attr.D("M12 16v-4")),
			html.Element("path", attr.D("M12 8h.01")),
		}
	case IconSend:
		paths = []html.Node{
			html.Element("path", attr.D("M14.536 21.686a.5.5 0 0 0 .937-.024l6.5-19a.496.496 0 0 0-.635-.635l-19 6.5a.5.5 0 0 0-.024.937l7.93 3.18a2 2 0 0 1 1.112 1.11z")),
			html.Element("path", attr.D("m21.854 2.147-10.94 10.939")),
		}
	case IconArrowRight:
		paths = []html.Node{
			html.Element("path", attr.D("M5 12h14")),
			html.Element("path", attr.D("m12 5 7 7-7 7")),
		}
	case IconLoader:
		paths = []html.Node{
			html.Element("path", attr.D("M12 2v4")),
			html.Element("path", attr.D("m16.2 7.8 2.9-2.9")),
			html.Element("path", attr.D("M18 12h4")),
			html.Element("path", attr.D("m16.2 16.2 2.9 2.9")),
			html.Element("path", attr.D("M12 18v4")),
			html.Element("path", attr.D("m4.9 19.1 2.9-2.9")),
			html.Element("path", attr.D("M2 12h4")),
			html.Element("path", attr.D("m4.9 4.9 2.9 2.9")),
		}
	case IconDownload:
		paths = []html.Node{
			html.Element("path", attr.D("M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4")),
			html.Element("polyline", attr.Points("7 10 12 15 17 10")),
			html.Element("line", attr.X1("12"), attr.X2("12"), attr.Y1("15"), attr.Y2("3")),
		}
	case IconUpload:
		paths = []html.Node{
			html.Element("path", attr.D("M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4")),
			html.Element("polyline", attr.Points("17 8 12 3 7 8")),
			html.Element("line", attr.X1("12"), attr.X2("12"), attr.Y1("3"), attr.Y2("15")),
		}
	case IconMoreVertical:
		paths = []html.Node{
			html.Element("circle", attr.Cx("12"), attr.Cy("12"), attr.R("1")),
			html.Element("circle", attr.Cx("12"), attr.Cy("5"), attr.R("1")),
			html.Element("circle", attr.Cx("12"), attr.Cy("19"), attr.R("1")),
		}
	case IconTrash:
		paths = []html.Node{
			html.Element("path", attr.D("M3 6h18")),
			html.Element("path", attr.D("M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6")),
			html.Element("path", attr.D("M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2")),
			html.Element("line", attr.X1("10"), attr.X2("10"), attr.Y1("11"), attr.Y2("17")),
			html.Element("line", attr.X1("14"), attr.X2("14"), attr.Y1("11"), attr.Y2("17")),
		}
	case IconPlus:
		paths = []html.Node{
			html.Element("path", attr.D("M5 12h14")),
			html.Element("path", attr.D("M12 5v14")),
		}
	case IconRepository:
		paths = []html.Node{
			html.Element("path", attr.D("M4 19.5v-15A2.5 2.5 0 0 1 6.5 2H19a1 1 0 0 1 1 1v18a1 1 0 0 1-1 1H6.5a1 1 0 0 1 0-5H20")),
			html.Element("path", attr.D("M8 7h6")),
			html.Element("path", attr.D("M8 11h8")),
		}
	case IconBuilding:
		paths = []html.Node{
			html.Element("path", attr.D("M3 21h18")),
			html.Element("path", attr.D("M5 21V7l8-4v18")),
			html.Element("path", attr.D("M19 21V11l-6-4")),
			html.Element("path", attr.D("M9 9v.01")),
			html.Element("path", attr.D("M9 12v.01")),
			html.Element("path", attr.D("M9 15v.01")),
			html.Element("path", attr.D("M9 18v.01")),
		}
	case IconUser:
		paths = []html.Node{
			html.Element("path", attr.D("M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2")),
			html.Element("circle", attr.Cx("12"), attr.Cy("7"), attr.R("4")),
		}
	case IconSettings:
		paths = []html.Node{
			html.Element("path", attr.D("M12.22 2h-.44a2 2 0 0 0-2 2v.18a2 2 0 0 1-1 1.73l-.43.25a2 2 0 0 1-2 0l-.15-.08a2 2 0 0 0-2.73.73l-.22.38a2 2 0 0 0 .73 2.73l.15.1a2 2 0 0 1 1 1.72v.51a2 2 0 0 1-1 1.74l-.15.09a2 2 0 0 0-.73 2.73l.22.38a2 2 0 0 0 2.73.73l.15-.08a2 2 0 0 1 2 0l.43.25a2 2 0 0 1 1 1.73V20a2 2 0 0 0 2 2h.44a2 2 0 0 0 2-2v-.18a2 2 0 0 1 1-1.73l.43-.25a2 2 0 0 1 2 0l.15.08a2 2 0 0 0 2.73-.73l.22-.39a2 2 0 0 0-.73-2.73l-.15-.08a2 2 0 0 1-1-1.74v-.5a2 2 0 0 1 1-1.74l.15-.09a2 2 0 0 0 .73-2.73l-.22-.38a2 2 0 0 0-2.73-.73l-.15.08a2 2 0 0 1-2 0l-.43-.25a2 2 0 0 1-1-1.73V4a2 2 0 0 0-2-2z")),
			html.Element("circle", attr.Cx("12"), attr.Cy("12"), attr.R("3")),
		}
	case IconLogOut:
		paths = []html.Node{
			html.Element("path", attr.D("M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4")),
			html.Element("polyline", attr.Points("16 17 21 12 16 7")),
			html.Element("line", attr.X1("21"), attr.X2("9"), attr.Y1("12"), attr.Y2("12")),
		}
	case IconMail:
		paths = []html.Node{
			html.Element("rect", attr.Width("20"), attr.Height("16"), attr.X("2"), attr.Y("4"), attr.Rx("2")),
			html.Element("path", attr.D("m22 7-8.97 5.7a1.94 1.94 0 0 1-2.06 0L2 7")),
		}
	case IconAtSign:
		paths = []html.Node{
			html.Element("circle", attr.Cx("12"), attr.Cy("12"), attr.R("4")),
			html.Element("path", attr.D("M16 8v5a3 3 0 0 0 6 0v-1a10 10 0 1 0-4 8")),
		}
	case IconLock:
		paths = []html.Node{
			html.Element("rect", attr.Width("18"), attr.Height("11"), attr.X("3"), attr.Y("11"), attr.Rx("2"), attr.Ry("2")),
			html.Element("path", attr.D("M7 11V7a5 5 0 0 1 10 0v4")),
		}
	case IconGitBranch:
		paths = []html.Node{
			html.Element("line", attr.X1("6"), attr.X2("6"), attr.Y1("3"), attr.Y2("15")),
			html.Element("circle", attr.Cx("18"), attr.Cy("6"), attr.R("3")),
			html.Element("circle", attr.Cx("6"), attr.Cy("18"), attr.R("3")),
			html.Element("path", attr.D("M18 9a9 9 0 0 1-9 9")),
		}
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
	case IconCopy:
		paths = []html.Node{
			html.Element("rect", attr.Width("14"), attr.Height("14"), attr.X("8"), attr.Y("8"), attr.Rx("2"), attr.Ry("2")),
			html.Element("path", attr.D("M4 16c-1.1 0-2-.9-2-2V4c0-1.1.9-2 2-2h10c1.1 0 2 .9 2 2")),
		}
	case IconGlobe:
		paths = []html.Node{
			html.Element("circle", attr.Cx("12"), attr.Cy("12"), attr.R("10")),
			html.Element("path", attr.D("M12 2a14.5 14.5 0 0 0 0 20 14.5 14.5 0 0 0 0-20")),
			html.Element("path", attr.D("M2 12h20")),
		}
	case IconTwitter:
		paths = []html.Node{
			html.Element("path", attr.D("M4 4l11.733 16h4.267l-11.733 -16z")),
			html.Element("path", attr.D("M4 20l6.768 -6.768m2.46 -2.46l6.772 -6.772")),
		}
	case IconDiscord:
		discordAttrs := []html.Node{
			attr.Xmlns("http://www.w3.org/2000/svg"),
			attr.Width("24"),
			attr.Height("24"),
			attr.ViewBox("0 0 127 96"),
			attr.Fill("currentColor"),
		}
		if class != "" {
			discordAttrs = append(discordAttrs, attr.Class(class))
		}
		discordAttrs = append(discordAttrs, html.Element("path", attr.D("M81.15,0c-1.2376,2.1973-2.3489,4.4704-3.3591,6.794-9.5975-1.4396-19.3718-1.4396-28.9945,0-.985-2.3236-2.1216-4.5967-3.3591-6.794-9.0166,1.5407-17.8059,4.2431-26.1405,8.0568C2.779,32.5304-1.6914,56.3725.5312,79.8863c9.6732,7.1476,20.5083,12.603,32.0505,16.0884,2.6014-3.4854,4.8998-7.1981,6.8698-11.0623-3.738-1.3891-7.3497-3.1318-10.8098-5.1523.9092-.6567,1.7932-1.3386,2.6519-1.9953,20.281,9.547,43.7696,9.547,64.0758,0,.8587.7072,1.7427,1.3891,2.6519,1.9953-3.4601,2.0457-7.0718,3.7632-10.835,5.1776,1.97,3.8642,4.2683,7.5769,6.8698,11.0623,11.5419-3.4854,22.3769-8.9156,32.0509-16.0631,2.626-27.2771-4.496-50.9172-18.817-71.8548C98.9811,4.2684,90.1918,1.5659,81.1752.0505l-.0252-.0505ZM42.2802,65.4144c-6.2383,0-11.4159-5.6575-11.4159-12.6535s4.9755-12.6788,11.3907-12.6788,11.5169,5.708,11.4159,12.6788c-.101,6.9708-5.026,12.6535-11.3907,12.6535ZM84.3576,65.4144c-6.2637,0-11.3907-5.6575-11.3907-12.6535s4.9755-12.6788,11.3907-12.6788,11.4917,5.708,11.3906,12.6788c-.101,6.9708-5.026,12.6535-11.3906,12.6535Z")))
		return html.Element("svg", discordAttrs...)
	case IconGitHub:
		githubAttrs := []html.Node{
			attr.Xmlns("http://www.w3.org/2000/svg"),
			attr.Width("24"),
			attr.Height("24"),
			attr.ViewBox("0 0 20 20"),
			attr.Fill("currentColor"),
		}
		if class != "" {
			githubAttrs = append(githubAttrs, attr.Class(class))
		}
		githubAttrs = append(githubAttrs, html.Element("path", attr.D("M10 0C4.477 0 0 4.477 0 10c0 4.42 2.865 8.17 6.839 9.49.5.092.682-.217.682-.482 0-.237-.008-.866-.013-1.7-2.782.603-3.369-1.34-3.369-1.34-.454-1.156-1.11-1.463-1.11-1.463-.908-.62.069-.608.069-.608 1.003.07 1.531 1.03 1.531 1.03.892 1.529 2.341 1.087 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.11-4.555-4.943 0-1.091.39-1.984 1.029-2.683-.103-.253-.446-1.27.098-2.647 0 0 .84-.269 2.75 1.025A9.578 9.578 0 0 1 10 4.836c.85.004 1.705.114 2.504.336 1.909-1.294 2.747-1.025 2.747-1.025.546 1.377.203 2.394.1 2.647.64.699 1.028 1.592 1.028 2.683 0 3.842-2.339 4.687-4.566 4.935.359.309.678.919.678 1.852 0 1.336-.012 2.415-.012 2.743 0 .267.18.578.688.48C17.138 18.167 20 14.418 20 10c0-5.523-4.477-10-10-10z")))
		return html.Element("svg", githubAttrs...)
	case IconBluesky:
		blueskyAttrs := []html.Node{
			attr.Xmlns("http://www.w3.org/2000/svg"),
			attr.Width("24"),
			attr.Height("24"),
			attr.ViewBox("0 0 600 530"),
			attr.Fill("currentColor"),
		}
		if class != "" {
			blueskyAttrs = append(blueskyAttrs, attr.Class(class))
		}
		blueskyAttrs = append(blueskyAttrs, html.Element("path", attr.D("M135.72 44.03c66.496 49.921 138.02 151.14 164.28 205.46 26.262-54.316 97.782-155.54 164.28-205.46 47.98-36.021 125.72-63.892 125.72 24.795 0 17.712-10.155 148.79-16.111 170.07-20.703 73.984-96.144 92.854-163.25 81.433 117.3 19.964 147.14 86.092 82.697 152.22-122.39 125.59-175.91-31.511-189.63-71.766-2.514-7.3797-3.6904-10.832-3.7077-7.8964-.0174-2.9357-1.1937.51669-3.7077 7.8964-13.714 40.255-67.233 197.36-189.63 71.766-64.444-66.128-34.605-132.26 82.697-152.22-67.108 11.421-142.55-7.4491-163.25-81.433-5.9562-21.282-16.111-152.36-16.111-170.07 0-88.687 77.742-60.816 125.72-24.795z")))
		return html.Element("svg", blueskyAttrs...)
	case IconStar:
		paths = []html.Node{
			html.Element("polygon", attr.Points("12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2")),
		}
	}

	return html.Element("svg", append(svgAttrs, paths...)...)
}

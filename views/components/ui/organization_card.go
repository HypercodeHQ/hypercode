package ui

import (
	"github.com/hyperstitieux/hypercode/html"
	"github.com/hyperstitieux/hypercode/html/attr"
)

type OrganizationCardProps struct {
	Username    string
	DisplayName string
}

func OrganizationCard(props OrganizationCardProps) html.Node {
	return html.A(
		attr.Href("/"+props.Username),
		attr.Class("card block p-4 hover:border-zinc-300 transition-all"),
		html.Div(
			attr.Class("space-y-3"),
			// Icon and username
			html.Div(
				attr.Class("flex items-center gap-3"),
				html.Div(
					attr.Class("flex items-center justify-center w-10 h-10 rounded-full bg-muted"),
					SVGIcon(IconUsers, "h-5 w-5 text-muted-foreground"),
				),
				html.Div(
					attr.Class("flex flex-col min-w-0"),
					html.H3(
						attr.Class("font-medium text-sm truncate"),
						html.Text(props.DisplayName),
					),
					html.P(
						attr.Class("text-xs text-muted-foreground truncate"),
						html.Text("@"+props.Username),
					),
				),
			),
		),
	)
}

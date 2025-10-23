package ui

import (
	"fmt"

	"github.com/hyperstitieux/hypercode/html"
	"github.com/hyperstitieux/hypercode/html/attr"
)

type RepositoryCardProps struct {
	OwnerUsername string
	Name          string
	IsPublic      bool
	StarCount     int64
}

func RepositoryCard(props RepositoryCardProps) html.Node {
	visibilityText := "Private"
	if props.IsPublic {
		visibilityText = "Public"
	}

	return html.A(
		attr.Href("/"+props.OwnerUsername+"/"+props.Name),
		attr.Class("card hover:opacity-70 transition-opacity"),
		html.Element("header",
			attr.Class("flex flex-col flex-wrap gap-4"),
			html.Div(
				attr.Class("flex items-center justify-between gap-2"),
				html.Element("span",
					attr.Class("badge-outline"),
					html.Text(visibilityText),
				),
				html.Div(
					attr.Class("flex items-center gap-1 text-muted-foreground text-sm"),
					SVGIcon(IconStar, "size-4"),
					html.Text(formatStarCount(props.StarCount)),
				),
			),
			html.H2(
				html.Text(props.OwnerUsername+"/"+props.Name),
			),
		),
	)
}

func formatStarCount(count int64) string {
	return fmt.Sprintf("%d", count)
}

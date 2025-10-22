package ui

import (
	"github.com/hyperstitieux/hypercode/html"
	"github.com/hyperstitieux/hypercode/html/attr"
)

type RepositoryCardProps struct {
	OwnerUsername string
	Name          string
	IsPublic      bool
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
			html.Element("span",
				attr.Class("badge-outline"),
				html.Text(visibilityText),
			),
			html.H2(
				html.Text(props.OwnerUsername+"/"+props.Name),
			),
		),
	)
}

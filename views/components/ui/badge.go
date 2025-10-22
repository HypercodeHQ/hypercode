package ui

import (
	"github.com/hyperstitieux/hypercode/html"
	"github.com/hyperstitieux/hypercode/html/attr"
)

type BadgeVariant string

const (
	BadgePrimary     BadgeVariant = "badge"
	BadgeSecondary   BadgeVariant = "badge-secondary"
	BadgeOutline     BadgeVariant = "badge-outline"
	BadgeDestructive BadgeVariant = "badge-destructive"
)

type BadgeProps struct {
	Variant BadgeVariant
	Class   string
	Href    string
}

func Badge(props BadgeProps, children ...html.Node) html.Node {
	className := string(props.Variant)
	if props.Class != "" {
		className += " " + props.Class
	}

	attrs := []html.Node{attr.Class(className)}

	if props.Href != "" {
		attrs = append(attrs, attr.Href(props.Href))
		return html.A(append(attrs, children...)...)
	}

	return html.Element("span", append(attrs, children...)...)
}

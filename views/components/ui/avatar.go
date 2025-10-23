package ui

import (
	"github.com/hypercodehq/hypercode/html"
	"github.com/hypercodehq/hypercode/html/attr"
)

type AvatarProps struct {
	Src      string
	Alt      string
	Fallback string
	Size     string
	Rounded  string
	Class    string
}

func Avatar(props AvatarProps) html.Node {
	className := "shrink-0 object-cover"

	if props.Size != "" {
		className += " " + props.Size
	} else {
		className += " size-8"
	}

	if props.Rounded != "" {
		className += " " + props.Rounded
	} else {
		className += " rounded-full"
	}

	if props.Class != "" {
		className += " " + props.Class
	}

	if props.Src != "" {
		return html.Element("img",
			attr.Class(className),
			attr.Alt(props.Alt),
			attr.Src(props.Src),
		)
	}

	if props.Fallback != "" {
		return html.Element("span",
			attr.Class(props.Size+" shrink-0 bg-muted flex items-center justify-center "+props.Rounded),
			html.Text(props.Fallback),
		)
	}

	return nil
}

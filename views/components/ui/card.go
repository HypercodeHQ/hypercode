package ui

import (
	"github.com/hypercodehq/hypercode/html"
	"github.com/hypercodehq/hypercode/html/attr"
)

type CardProps struct {
	Header      html.Node
	Title       string
	Description string
	Content     html.Node
	Footer      html.Node
	Class       string
}

func Card(props CardProps) html.Node {
	className := "card"
	if props.Class != "" {
		className += " " + props.Class
	}

	children := []html.Node{}

	if props.Header != nil {
		children = append(children, html.Element("header", props.Header))
	} else if props.Title != "" || props.Description != "" {
		headerChildren := []html.Node{}
		if props.Title != "" {
			headerChildren = append(headerChildren, html.H2(html.Text(props.Title)))
		}
		if props.Description != "" {
			headerChildren = append(headerChildren, html.P(html.Text(props.Description)))
		}
		children = append(children, html.Element("header", headerChildren...))
	}

	if props.Content != nil {
		children = append(children, html.Element("section", props.Content))
	}

	if props.Footer != nil {
		children = append(children, html.Element("footer", props.Footer))
	}

	return html.Div(
		attr.Class(className),
		html.For(children, func(child html.Node) html.Node {
			return child
		}),
	)
}

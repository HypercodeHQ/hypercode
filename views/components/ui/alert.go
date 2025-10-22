package ui

import (
	"github.com/hyperstitieux/hypercode/html"
	"github.com/hyperstitieux/hypercode/html/attr"
)

type AlertVariant string

const (
	AlertDefault     AlertVariant = "alert"
	AlertDestructive AlertVariant = "alert-destructive"
)

type AlertProps struct {
	Variant     AlertVariant
	Title       string
	Description string
	Icon        html.Node
	Class       string
}

func Alert(props AlertProps) html.Node {
	className := string(props.Variant)
	if props.Class != "" {
		className += " " + props.Class
	}

	children := []html.Node{}

	if props.Icon != nil {
		children = append(children, props.Icon)
	}

	if props.Title != "" {
		children = append(children, html.H2(html.Text(props.Title)))
	}

	if props.Description != "" {
		children = append(children, html.Element("section", html.Text(props.Description)))
	}

	return html.Div(
		attr.Class(className),
		html.For(children, func(child html.Node) html.Node {
			return child
		}),
	)
}

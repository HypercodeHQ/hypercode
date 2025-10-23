package ui

import (
	"github.com/hypercodehq/libhtml"
	"github.com/hypercodehq/libhtml/attr"
)

type InputProps struct {
	Type        string
	Name        string
	Id          string
	Placeholder string
	Value       string
	Required    bool
	Disabled    bool
	Class       string
}

func Input(props InputProps) html.Node {
	className := "input"
	if props.Class != "" {
		className += " " + props.Class
	}

	attrs := []html.Node{attr.Class(className)}

	if props.Type != "" {
		attrs = append(attrs, attr.Type(props.Type))
	} else {
		attrs = append(attrs, attr.Type("text"))
	}

	if props.Name != "" {
		attrs = append(attrs, attr.Name(props.Name))
	}

	if props.Id != "" {
		attrs = append(attrs, attr.Id(props.Id))
	}

	if props.Placeholder != "" {
		attrs = append(attrs, attr.Placeholder(props.Placeholder))
	}

	if props.Value != "" {
		attrs = append(attrs, attr.Value(props.Value))
	}

	if props.Required {
		attrs = append(attrs, attr.Required())
	}

	if props.Disabled {
		attrs = append(attrs, attr.Disabled())
	}

	return html.Input(attrs...)
}

func LabelFor(forId string, children ...html.Node) html.Node {
	return html.Label(
		attr.For(forId),
		html.For(children, func(child html.Node) html.Node {
			return child
		}),
	)
}

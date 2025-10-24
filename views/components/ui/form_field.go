package ui

import (
	"github.com/hypercommithq/libhtml"
	"github.com/hypercommithq/libhtml/attr"
)

type FormFieldProps struct {
	Label        string
	Id           string
	Name         string
	Type         string
	Placeholder  string
	Icon         Icon
	Required     bool
	Value        string
	Class        string
	WrapperClass string
	Error        string
}

func FormField(props FormFieldProps) html.Node {
	labelClass := "label"

	iconNode := html.Node(nil)
	inputClass := "input"

	if props.Icon != "" {
		iconNode = html.Div(
			attr.Class("absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground"),
			SVGIcon(props.Icon, "h-4 w-4"),
		)
		inputClass = "input pl-9"
	}

	if props.Error != "" {
		inputClass += " border-destructive focus:ring-destructive"
	}

	if props.Class != "" {
		inputClass += " " + props.Class
	}

	inputAttrs := []html.Node{
		attr.Class(inputClass),
		attr.Type(props.Type),
		attr.Id(props.Id),
		attr.Name(props.Name),
	}

	if props.Placeholder != "" {
		inputAttrs = append(inputAttrs, attr.Placeholder(props.Placeholder))
	}

	if props.Value != "" {
		inputAttrs = append(inputAttrs, attr.Value(props.Value))
	}

	if props.Required {
		inputAttrs = append(inputAttrs, attr.Required())
	}

	wrapperClass := "space-y-2"
	if props.WrapperClass != "" {
		wrapperClass += " " + props.WrapperClass
	}

	return html.Div(
		attr.Class(wrapperClass),
		html.Label(
			attr.For(props.Id),
			attr.Class(labelClass),
			html.Text(props.Label),
		),
		html.Div(
			attr.Class("relative"),
			iconNode,
			html.Input(inputAttrs...),
		),
		html.If(props.Error != "", html.P(
			attr.Class("text-sm text-destructive"),
			html.Text(props.Error),
		)),
	)
}

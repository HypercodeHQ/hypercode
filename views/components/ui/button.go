package ui

import (
	"github.com/hypercommithq/libhtml"
	"github.com/hypercommithq/libhtml/attr"
)

type ButtonVariant string

const (
	ButtonPrimary     ButtonVariant = "btn-primary"
	ButtonSecondary   ButtonVariant = "btn-secondary"
	ButtonOutline     ButtonVariant = "btn-outline"
	ButtonGhost       ButtonVariant = "btn-ghost"
	ButtonDestructive ButtonVariant = "btn-destructive"
	ButtonLink        ButtonVariant = "btn-link"
)

type ButtonSize string

const (
	ButtonDefault ButtonSize = ""
	ButtonSmall   ButtonSize = "sm"
	ButtonLarge   ButtonSize = "lg"
)

type ButtonProps struct {
	Variant  ButtonVariant
	Size     ButtonSize
	Icon     bool
	Disabled bool
	Type     string
	Class    string
	OnClick  string
}

func Button(props ButtonProps, children ...html.Node) html.Node {
	className := string(props.Variant)

	if props.Size != "" {
		className = "btn-" + string(props.Size) + "-" + string(props.Variant)[4:]
	}

	if props.Icon {
		if props.Size != "" {
			className = "btn-" + string(props.Size) + "-icon-" + string(props.Variant)[4:]
		} else {
			className = "btn-icon-" + string(props.Variant)[4:]
		}
	}

	if props.Class != "" {
		className += " " + props.Class
	}

	attrs := []html.Node{attr.Class(className)}

	if props.Type != "" {
		attrs = append(attrs, attr.Type(props.Type))
	} else {
		attrs = append(attrs, attr.Type("button"))
	}

	if props.Disabled {
		attrs = append(attrs, attr.Disabled())
	}

	if props.OnClick != "" {
		attrs = append(attrs, attr.Onclick(props.OnClick))
	}

	return html.Element("button", append(attrs, children...)...)
}

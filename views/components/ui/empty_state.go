package ui

import (
	"github.com/hyperstitieux/hypercode/html"
	"github.com/hyperstitieux/hypercode/html/attr"
)

type EmptyStateProps struct {
	Icon        html.Node
	Title       string
	Description string
	ActionText  string
	ActionHref  string
	ShowAction  bool
}

func EmptyState(props EmptyStateProps) html.Node {
	actionNode := html.Node(nil)
	if props.ShowAction {
		actionNode = html.Div(
			attr.Class("flex w-full max-w-sm min-w-0 flex-col items-center gap-4 text-sm"),
			html.A(
				attr.Href(props.ActionHref),
				attr.Class("btn"),
				SVGIcon(IconPlus, "h-4 w-4"),
				html.Text(props.ActionText),
			),
		)
	}

	return html.Div(
		attr.Class("flex min-w-0 flex-1 flex-col items-center justify-center gap-6 rounded-lg border-dashed p-6 text-center md:p-12 min-h-[calc(100vh-61px-4rem)]"),
		html.Div(
			attr.Class("flex max-w-sm flex-col items-center gap-2 text-center"),
			html.Div(
				attr.Class("flex shrink-0 items-center justify-center mb-2 bg-white border text-foreground flex size-10 shrink-0 items-center justify-center rounded-lg"),
				props.Icon,
			),
			html.Div(
				attr.Class("text-lg font-medium tracking-tight"),
				html.Text(props.Title),
			),
			html.Div(
				attr.Class("text-muted-foreground text-sm"),
				html.Text(props.Description),
			),
		),
		actionNode,
	)
}

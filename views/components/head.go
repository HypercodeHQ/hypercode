package components

import (
	"github.com/hypercommithq/libhtml"
	"github.com/hypercommithq/libhtml/attr"
)

func Head(title string) html.Node {
	return html.Head(
		html.Meta(attr.Charset("utf-8")),
		html.Meta(
			attr.Name("viewport"),
			attr.Content("width=device-width, initial-scale=1"),
		),
		html.Title(html.Text(title)),
		html.Link(
			attr.Rel("icon"),
			attr.Href("/favicon.ico"),
		),
		html.Link(
			attr.Rel("preconnect"),
			attr.Href("https://fonts.bunny.net"),
			attr.Crossorigin("anonymous"),
		),
		html.Link(
			attr.Rel("stylesheet"),
			attr.Href("https://fonts.bunny.net/css?family=ibm-plex-sans:400,500,600,700|ibm-plex-mono:400,500,600,700"),
			attr.Crossorigin("anonymous"),
		),
		html.Link(
			attr.Rel("preload"),
			attr.Href("/logo.png"),
			attr.As("image"),
		),
		html.Link(
			attr.Rel("stylesheet"),
			attr.Href("/styles.css"),
		),
		html.Script(
			attr.Src("/dropdown.js"),
			attr.Defer(),
		),
		html.Script(
			attr.Src("https://cdn.jsdelivr.net/npm/basecoat-css@0.3.2/dist/js/basecoat.min.js"),
			attr.Defer(),
		),
		html.Script(
			attr.Src("https://cdn.jsdelivr.net/npm/basecoat-css@0.3.2/dist/js/toast.min.js"),
			attr.Defer(),
		),
	)
}

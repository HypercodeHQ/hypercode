package ui

import (
	"github.com/hyperstitieux/hypercode/html"
	"github.com/hyperstitieux/hypercode/html/attr"
)

// GitHubAuthButton renders a "Continue with GitHub" button
func GitHubAuthButton() html.Node {
	return html.A(
		attr.Href("/auth/github"),
		attr.Class("btn-outline w-full flex items-center justify-center gap-2"),
		SVGIcon(IconGitHub, "h-5 w-5"),
		html.Text("Continue with GitHub"),
	)
}

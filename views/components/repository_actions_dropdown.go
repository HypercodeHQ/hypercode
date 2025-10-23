package components

import (
	"fmt"
	"net/url"

	"github.com/hyperstitieux/hypercode/html"
	"github.com/hyperstitieux/hypercode/html/attr"
	"github.com/hyperstitieux/hypercode/views/components/ui"
)

type RepositoryActionsDropdownData struct {
	OwnerUsername string
	RepoName      string
	CloneURL      string
	RepositoryURL string
}

func ShareDropdown(data *RepositoryActionsDropdownData) html.Node {
	if data == nil {
		data = &RepositoryActionsDropdownData{}
	}

	return html.Div(
		attr.Class("dropdown-menu"),
		// Dropdown button - styled like Star button
		html.Element("button",
			attr.Type("button"),
			attr.Class("btn-outline inline-flex items-center gap-2"),
			attr.Id("share-trigger"),
			attr.AriaHaspopup("menu"),
			attr.AriaControls("share-menu"),
			attr.AriaExpanded("false"),
			ui.SVGIcon(ui.IconShare, "size-4"),
			html.Text("Share"),
		),
		// Dropdown menu
		html.Div(
			attr.Id("share-popover"),
			attr.DataPopover(""),
			attr.AriaHidden("true"),
			attr.Class("min-w-56 right-0 left-auto"),
			html.Div(
				attr.Role("menu"),
				attr.Id("share-menu"),
				attr.AriaLabelledby("share-trigger"),
				shareSection(data),
			),
		),
		// JavaScript for copy functionality
		copyScript(),
	)
}

func CloneDropdown(data *RepositoryActionsDropdownData) html.Node {
	if data == nil {
		data = &RepositoryActionsDropdownData{}
	}

	return html.Div(
		attr.Class("dropdown-menu"),
		// Dropdown button - styled like Star button
		html.Element("button",
			attr.Type("button"),
			attr.Class("btn-outline inline-flex items-center gap-2"),
			attr.Id("clone-trigger"),
			attr.AriaHaspopup("menu"),
			attr.AriaControls("clone-menu"),
			attr.AriaExpanded("false"),
			ui.SVGIcon(ui.IconCode, "size-4"),
			html.Text("Clone"),
		),
		// Dropdown menu
		html.Div(
			attr.Id("clone-popover"),
			attr.DataPopover(""),
			attr.AriaHidden("true"),
			attr.Class("min-w-80 right-0 left-auto"),
			html.Div(
				attr.Role("menu"),
				attr.Id("clone-menu"),
				attr.AriaLabelledby("clone-trigger"),
				cloneSection(data),
			),
		),
		// JavaScript for copy functionality
		copyScript(),
	)
}

func shareSection(data *RepositoryActionsDropdownData) html.Node {
	return html.Group(
		// Copy link button
		html.Element("button",
			attr.Type("button"),
			attr.Role("menuitem"),
			attr.Class("cursor-pointer"),
			attr.Onclick(fmt.Sprintf("copyToClipboard('%s', 'share-link-button')", data.RepositoryURL)),
			attr.Id("share-link-button"),
			html.Div(
				attr.Class("flex items-center gap-2"),
				ui.SVGIcon(ui.IconLink, ""),
				html.Text("Copy link"),
			),
		),
		html.Hr(attr.Role("separator")),
		// Share to X
		html.A(
			attr.Href(getTwitterShareURL(data.RepositoryURL, data.OwnerUsername, data.RepoName)),
			attr.Target("_blank"),
			attr.Rel("noopener noreferrer"),
			attr.Role("menuitem"),
			attr.Class("cursor-pointer"),
			ui.SVGIcon(ui.IconTwitter, ""),
			html.Text("Share to X"),
		),
		html.Hr(attr.Role("separator")),
		// Share to Bluesky
		html.A(
			attr.Href(getBlueskyShareURL(data.RepositoryURL)),
			attr.Target("_blank"),
			attr.Rel("noopener noreferrer"),
			attr.Role("menuitem"),
			attr.Class("cursor-pointer"),
			ui.SVGIcon(ui.IconBluesky, ""),
			html.Text("Share to Bluesky"),
		),
	)
}

func cloneSection(data *RepositoryActionsDropdownData) html.Node {
	cloneCommand := "git clone " + data.CloneURL

	return html.Div(
		attr.Class("p-3"),
		// Label
		html.Label(
			attr.For("clone-url-input"),
			attr.Class("label mb-2 block"),
			html.Text("Clone Command"),
		),
		// Clone command input and button
		html.Div(
			attr.Class("flex gap-2"),
			html.Input(
				attr.Type("text"),
				attr.Readonly(),
				attr.Value(cloneCommand),
				attr.Id("clone-url-input"),
				attr.Class("input font-mono text-sm flex-1"),
			),
			html.Element("button",
				attr.Type("button"),
				attr.Onclick("copyCloneURL()"),
				attr.Id("clone-copy-button"),
				attr.Class("btn-icon-outline cursor-pointer"),
				attr.DataTooltip("Copy to clipboard"),
				attr.DataSide("left"),
				ui.SVGIcon(ui.IconCopy, "size-4 copy-icon"),
				ui.SVGIcon(ui.IconCheck, "size-4 check-icon hidden"),
			),
		),
	)
}

func copyScript() html.Node {
	return html.Element("script",
		html.Text(`
// Copy to clipboard helper
if (typeof copyToClipboard === 'undefined') {
	window.copyToClipboard = function(text, buttonId) {
		const button = document.getElementById(buttonId);
		const copyIcon = button.querySelector('.copy-icon');
		const checkIcon = button.querySelector('.check-icon');

		if (navigator.clipboard && navigator.clipboard.writeText) {
			navigator.clipboard.writeText(text).then(() => {
				showCopySuccess(copyIcon, checkIcon);
			}).catch((err) => {
				console.error('Failed to copy:', err);
			});
		}
	};
}

// Copy clone URL
if (typeof copyCloneURL === 'undefined') {
	window.copyCloneURL = function() {
		const input = document.getElementById('clone-url-input');
		const button = document.getElementById('clone-copy-button');
		const copyIcon = button.querySelector('.copy-icon');
		const checkIcon = button.querySelector('.check-icon');

		const cloneCommand = input.value;

		if (navigator.clipboard && navigator.clipboard.writeText) {
			navigator.clipboard.writeText(cloneCommand).then(() => {
				showCloneSuccess(button, copyIcon, checkIcon);
			}).catch((err) => {
				console.error('Failed to copy:', err);
				fallbackCopy();
			});
		} else {
			fallbackCopy();
		}

		function fallbackCopy() {
			input.select();
			input.setSelectionRange(0, 99999);
			try {
				document.execCommand('copy');
				showCloneSuccess(button, copyIcon, checkIcon);
			} catch (err) {
				console.error('Fallback copy failed:', err);
			}
		}

		function showCloneSuccess(btn, copyIcn, checkIcn) {
			btn.setAttribute('data-tooltip', 'Copied!');
			copyIcn.classList.add('hidden');
			checkIcn.classList.remove('hidden');
			setTimeout(() => {
				btn.setAttribute('data-tooltip', 'Copy to clipboard');
				copyIcn.classList.remove('hidden');
				checkIcn.classList.add('hidden');
			}, 2000);

			// Show toast notification - ensure toaster exists
			const toaster = document.getElementById('toaster');
			if (toaster) {
				document.dispatchEvent(new CustomEvent('basecoat:toast', {
					detail: {
						config: {
							category: 'success',
							title: 'Clone command copied!',
							description: 'The clone command has been copied to your clipboard.',
							duration: 2000
						}
					}
				}));
			} else {
				console.error('Toaster container not found');
			}
		}
	};
}

if (typeof showCopySuccess === 'undefined') {
	window.showCopySuccess = function(copyIcon, checkIcon) {
		if (copyIcon && checkIcon) {
			copyIcon.classList.add('hidden');
			checkIcon.classList.remove('hidden');
			setTimeout(() => {
				copyIcon.classList.remove('hidden');
				checkIcon.classList.add('hidden');
			}, 2000);
		}

		// Show toast notification - ensure toaster exists
		const toaster = document.getElementById('toaster');
		if (toaster) {
			document.dispatchEvent(new CustomEvent('basecoat:toast', {
				detail: {
					config: {
						category: 'success',
						title: 'Link copied!',
						description: 'The link has been copied to your clipboard.',
						duration: 2000
					}
				}
			}));
		} else {
			console.error('Toaster container not found');
		}

		// Close the share dropdown
		const shareTrigger = document.getElementById('share-trigger');
		const sharePopover = document.getElementById('share-popover');
		if (shareTrigger && sharePopover) {
			shareTrigger.setAttribute('aria-expanded', 'false');
			sharePopover.setAttribute('aria-hidden', 'true');
		}
	};
}
		`),
	)
}

func getTwitterShareURL(repoURL, owner, repoName string) string {
	text := fmt.Sprintf("Check out %s/%s on Hypercode!", owner, repoName)
	return fmt.Sprintf("https://twitter.com/intent/tweet?text=%s&url=%s",
		url.QueryEscape(text),
		url.QueryEscape(repoURL))
}

func getBlueskyShareURL(repoURL string) string {
	text := "Check out this repository!"
	return fmt.Sprintf("https://bsky.app/intent/compose?text=%s%%20%s",
		url.QueryEscape(text),
		url.QueryEscape(repoURL))
}

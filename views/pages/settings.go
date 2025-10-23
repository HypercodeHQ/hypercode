package pages

import (
	"fmt"
	"net/http"
	"time"

	"github.com/hypercodehq/hypercode/database/models"
	"github.com/hypercodehq/libhtml"
	"github.com/hypercodehq/libhtml/attr"
	"github.com/hypercodehq/hypercode/views/components/layouts"
	"github.com/hypercodehq/hypercode/views/components/ui"
)

type SettingsData struct {
	User                 *models.User
	DisplayNameError     string
	UsernameError        string
	CurrentPasswordError string
	NewPasswordError     string
	ConfirmPasswordError string
	GeneralSuccess       string
	PasswordSuccess      string
	DisplayName          string
	Username             string
	AccessTokens         []*models.AccessToken
	NewAccessToken       string
	AccessTokenSuccess   string
	AccessTokenError     string
}

func Settings(r *http.Request, data *SettingsData) html.Node {
	if data == nil {
		data = &SettingsData{}
	}

	// Populate form values from user data if not set
	if data.User != nil {
		if data.DisplayName == "" {
			data.DisplayName = data.User.DisplayName
		}
		if data.Username == "" {
			data.Username = data.User.Username
		}
	}

	return layouts.Main(r,
		"Settings",
		html.Main(
			attr.Class("min-h-[calc(100vh-61px)] w-full mx-auto max-w-7xl space-y-6 py-8 px-4"),
			html.H1(
				attr.Class("font-semibold text-2xl mb-6"),
				html.Text("Settings"),
			),

			// General Settings Card
			ui.Card(ui.CardProps{
				Title:       "General",
				Description: "Update your account information",
				Content: html.Div(
					attr.Class("space-y-4"),
					html.If(data.GeneralSuccess != "", html.Div(
						attr.Class("p-3 rounded-lg bg-emerald-50 dark:bg-emerald-900/20 border border-emerald-200 dark:border-emerald-800 text-emerald-800 dark:text-emerald-200 text-sm"),
						html.Text(data.GeneralSuccess),
					)),
					html.Form(
						attr.Id("general-settings-form"),
						attr.Method("POST"),
						attr.Action("/settings/general"),
						attr.Class("space-y-4"),
						attr.Attribute{Key: "data-original-username", Value: data.User.Username},
						ui.FormField(ui.FormFieldProps{
							Label:       "Display Name",
							Id:          "display_name",
							Name:        "display_name",
							Type:        "text",
							Placeholder: "John Doe",
							Icon:        ui.IconUser,
							Required:    true,
							Value:       data.DisplayName,
							Error:       data.DisplayNameError,
						}),
						ui.FormField(ui.FormFieldProps{
							Label:       "Username",
							Id:          "username",
							Name:        "username",
							Type:        "text",
							Placeholder: "johndoe",
							Icon:        ui.IconAtSign,
							Required:    true,
							Value:       data.Username,
							Error:       data.UsernameError,
						}),
						html.Div(
							attr.Class("flex justify-end"),
							ui.Button(
								ui.ButtonProps{
									Variant: ui.ButtonPrimary,
									Type:    "submit",
								},
								html.Text("Save Changes"),
							),
						),
					),
				),
			}),

			// Password Settings Card
			ui.Card(ui.CardProps{
				Title:       "Password",
				Description: "Change your password",
				Content: html.Div(
					attr.Class("space-y-4"),
					html.If(data.PasswordSuccess != "", html.Div(
						attr.Class("p-3 rounded-lg bg-emerald-50 dark:bg-emerald-900/20 border border-emerald-200 dark:border-emerald-800 text-emerald-800 dark:text-emerald-200 text-sm"),
						html.Text(data.PasswordSuccess),
					)),
					html.Form(
						attr.Method("POST"),
						attr.Action("/settings/password"),
						attr.Class("space-y-4"),
						ui.FormField(ui.FormFieldProps{
							Label:       "Current Password",
							Id:          "current_password",
							Name:        "current_password",
							Type:        "password",
							Placeholder: "••••••••••••••••",
							Icon:        ui.IconLock,
							Required:    true,
							Error:       data.CurrentPasswordError,
						}),
						ui.FormField(ui.FormFieldProps{
							Label:       "New Password",
							Id:          "new_password",
							Name:        "new_password",
							Type:        "password",
							Placeholder: "••••••••••••••••",
							Icon:        ui.IconLock,
							Required:    true,
							Error:       data.NewPasswordError,
						}),
						ui.FormField(ui.FormFieldProps{
							Label:       "Confirm New Password",
							Id:          "confirm_password",
							Name:        "confirm_password",
							Type:        "password",
							Placeholder: "••••••••••••••••",
							Icon:        ui.IconLock,
							Required:    true,
							Error:       data.ConfirmPasswordError,
						}),
						html.Div(
							attr.Class("flex justify-end"),
							ui.Button(
								ui.ButtonProps{
									Variant: ui.ButtonPrimary,
									Type:    "submit",
								},
								html.Text("Update Password"),
							),
						),
					),
				),
			}),

			// Access Tokens Card
			html.Div(
				attr.Id("access-tokens"),
				ui.Card(ui.CardProps{
					Title:       "Access Tokens",
					Description: "Personal access tokens can be used as passwords for Git operations",
					Content: html.Div(
						attr.Class("space-y-4"),
						html.If(data.AccessTokenSuccess != "", html.Div(
							attr.Class("p-3 rounded-lg bg-emerald-50 dark:bg-emerald-900/20 border border-emerald-200 dark:border-emerald-800 text-emerald-800 dark:text-emerald-200 text-sm"),
							html.Text(data.AccessTokenSuccess),
						)),
						html.If(data.AccessTokenError != "", html.Div(
							attr.Class("p-3 rounded-lg bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-800 dark:text-red-200 text-sm"),
							html.Text(data.AccessTokenError),
						)),
						html.If(data.NewAccessToken != "", html.Div(
							attr.Class("p-4 rounded-lg bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800"),
							html.Div(
								attr.Class("flex items-start gap-3"),
								html.Div(
									attr.Class("flex-shrink-0 mt-0.5"),
									ui.SVGIcon(ui.IconAlertCircle, "text-yellow-600 dark:text-yellow-400"),
								),
								html.Div(
									attr.Class("flex-1 space-y-3"),
									html.P(
										attr.Class("text-sm font-medium text-yellow-800 dark:text-yellow-200"),
										html.Text("Make sure to copy your access token now. You won't be able to see it again!"),
									),
									html.Div(
										attr.Class("space-y-2"),
										html.Label(
											attr.For("new-token-value"),
											attr.Class("label block"),
											html.Text("Your New Access Token"),
										),
										html.Div(
											attr.Class("flex gap-2"),
											html.Input(
												attr.Type("text"),
												attr.Readonly(),
												attr.Value(data.NewAccessToken),
												attr.Id("new-token-value"),
												attr.Class("input font-mono text-sm flex-1"),
											),
											html.Element("button",
												attr.Type("button"),
												attr.Onclick("copyNewToken()"),
												attr.Id("copy-token-btn"),
												attr.Class("btn-icon-outline cursor-pointer"),
												attr.DataTooltip("Copy to clipboard"),
												attr.DataSide("top"),
												ui.SVGIcon(ui.IconCopy, "size-4"),
												ui.SVGIcon(ui.IconCheck, "size-4 hidden"),
											),
										),
									),
								),
							),
						)),
						html.Form(
							attr.Method("POST"),
							attr.Action("/settings/access-tokens"),
							attr.Class("space-y-4"),
							ui.FormField(ui.FormFieldProps{
								Label:       "Token Name",
								Id:          "token_name",
								Name:        "name",
								Type:        "text",
								Placeholder: "My Token",
								Icon:        ui.IconLock,
								Required:    true,
							}),
							html.Div(
								attr.Class("flex justify-end"),
								ui.Button(
									ui.ButtonProps{
										Variant: ui.ButtonPrimary,
										Type:    "submit",
									},
									html.Text("Generate Token"),
								),
							),
						),
						html.If(len(data.AccessTokens) > 0, html.Div(
							attr.Class("space-y-2 mt-6"),
							html.H3(
								attr.Class("text-sm font-semibold text-foreground mb-3"),
								html.Text("Active Tokens"),
							),
							html.Div(
								attr.Class("space-y-2"),
								html.Group(accessTokenList(data.AccessTokens)...),
							),
						)),
					),
				}),
			),

			// JavaScript for username change confirmation and token copying
			html.Element("script",
				html.Text(`
					(function() {
						const form = document.getElementById('general-settings-form');
						if (form) {
							form.addEventListener('submit', function(e) {
								const originalUsername = form.getAttribute('data-original-username');
								const currentUsername = document.getElementById('username').value;

								if (originalUsername !== currentUsername) {
									const confirmed = window.confirm(
										'Are you sure you want to change your username from "@' + originalUsername + '" to "@' + currentUsername + '"?\n\n' +
										'This may affect your profile URL and repository access.'
									);

									if (!confirmed) {
										e.preventDefault();
										return false;
									}
								}
							});
						}

						// Copy token to clipboard
						window.copyNewToken = function() {
							const input = document.getElementById('new-token-value');
							const button = document.getElementById('copy-token-btn');
							const copyIcon = button.querySelector('svg:nth-child(1)');
							const checkIcon = button.querySelector('svg:nth-child(2)');
							const tokenValue = input.value;

							if (navigator.clipboard && navigator.clipboard.writeText) {
								navigator.clipboard.writeText(tokenValue).then(function() {
									showTokenCopySuccess();
								}).catch(function(err) {
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
									showTokenCopySuccess();
								} catch (err) {
									console.error('Fallback copy failed:', err);
								}
							}

							function showTokenCopySuccess() {
								button.setAttribute('data-tooltip', 'Copied!');
								copyIcon.classList.add('hidden');
								checkIcon.classList.remove('hidden');
								setTimeout(function() {
									button.setAttribute('data-tooltip', 'Copy to clipboard');
									copyIcon.classList.remove('hidden');
									checkIcon.classList.add('hidden');
								}, 2000);

								// Show toast notification
								document.dispatchEvent(new CustomEvent('basecoat:toast', {
									detail: {
										config: {
											category: 'success',
											title: 'Access token copied!',
											description: 'Your access token has been copied to your clipboard.',
											duration: 2000
										}
									}
								}));
							}
						};

						// Handle token deletion confirmation
						document.querySelectorAll('[data-delete-token]').forEach(function(btn) {
							btn.addEventListener('click', function(e) {
								const tokenName = btn.getAttribute('data-token-name');
								if (!confirm('Are you sure you want to delete the token "' + tokenName + '"? This action cannot be undone.')) {
									e.preventDefault();
								}
							});
						});
					})();
				`),
			),
		),
	)
}

func accessTokenList(tokens []*models.AccessToken) []html.Node {
	if tokens == nil || len(tokens) == 0 {
		return []html.Node{}
	}

	nodes := make([]html.Node, 0, len(tokens))
	for _, token := range tokens {
		if token != nil {
			nodes = append(nodes, accessTokenItem(token))
		}
	}
	return nodes
}

func accessTokenItem(token *models.AccessToken) html.Node {
	if token == nil {
		return html.Div()
	}

	return html.Div(
		attr.Class("flex items-center justify-between p-3 bg-muted rounded-lg"),
		html.Div(
			attr.Class("flex-1"),
			html.Div(
				attr.Class("font-medium text-sm text-foreground"),
				html.Text(token.Name),
			),
			html.Div(
				attr.Class("text-xs text-muted-foreground mt-1"),
				html.Text("Created "+formatTimestamp(token.CreatedAt)),
				html.IfElsef(
					token.LastUsedAt != nil,
					func() html.Node {
						return html.Text(" • Last used " + formatTimestamp(*token.LastUsedAt))
					},
					func() html.Node {
						return html.Group()
					},
				),
			),
		),
		html.Form(
			attr.Method("POST"),
			attr.Action("/settings/access-tokens/"+fmt.Sprintf("%d", token.ID)+"/delete"),
			attr.Class("inline"),
			html.Element("button",
				attr.Type("submit"),
				attr.Class("btn-icon-ghost text-destructive hover:text-destructive"),
				attr.Attribute{Key: "data-delete-token", Value: "true"},
				attr.Attribute{Key: "data-token-name", Value: token.Name},
				attr.DataTooltip("Delete token"),
				attr.DataSide("left"),
				ui.SVGIcon(ui.IconTrash, "text-destructive"),
			),
		),
	)
}

func formatTimestamp(timestamp int64) string {
	// Convert Unix timestamp to a human-readable format
	// For now, just return a simple representation
	t := time.Unix(timestamp, 0)
	now := time.Now()

	diff := now.Sub(t)

	if diff < time.Minute {
		return "just now"
	} else if diff < time.Hour {
		mins := int(diff.Minutes())
		if mins == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", mins)
	} else if diff < 24*time.Hour {
		hours := int(diff.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	} else if diff < 7*24*time.Hour {
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "yesterday"
		}
		return fmt.Sprintf("%d days ago", days)
	} else if diff < 30*24*time.Hour {
		weeks := int(diff.Hours() / 24 / 7)
		if weeks == 1 {
			return "1 week ago"
		}
		return fmt.Sprintf("%d weeks ago", weeks)
	} else if diff < 365*24*time.Hour {
		months := int(diff.Hours() / 24 / 30)
		if months == 1 {
			return "1 month ago"
		}
		return fmt.Sprintf("%d months ago", months)
	} else {
		years := int(diff.Hours() / 24 / 365)
		if years == 1 {
			return "1 year ago"
		}
		return fmt.Sprintf("%d years ago", years)
	}
}

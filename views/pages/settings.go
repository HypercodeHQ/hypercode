package pages

import (
	"net/http"

	"github.com/hyperstitieux/hypercode/database/models"
	"github.com/hyperstitieux/hypercode/html"
	"github.com/hyperstitieux/hypercode/html/attr"
	"github.com/hyperstitieux/hypercode/views/components/layouts"
	"github.com/hyperstitieux/hypercode/views/components/ui"
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
			attr.Class("min-h-[calc(100vh-61px)] w-full mx-auto max-w-2xl space-y-6 py-8 px-4"),
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

			// JavaScript for username change confirmation
			html.Element("script",
				html.Text(`
					(function() {
						const form = document.getElementById('general-settings-form');
						if (!form) return;

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
					})();
				`),
			),
		),
	)
}

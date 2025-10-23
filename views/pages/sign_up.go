package pages

import (
	"net/http"

	"github.com/hypercodehq/libhtml"
	"github.com/hypercodehq/libhtml/attr"
	"github.com/hypercodehq/hypercode/views/components/layouts"
	"github.com/hypercodehq/hypercode/views/components/ui"
)

type SignUpData struct {
	DisplayNameError string
	UsernameError    string
	EmailError       string
	PasswordError    string
	DisplayName      string
	Username         string
	Email            string
}

func SignUp(r *http.Request, data *SignUpData) html.Node {
	if data == nil {
		data = &SignUpData{}
	}

	return layouts.Main(r,
		"Create an account",
		html.Main(
			attr.Class("min-h-[calc(100vh-61px)] flex flex-col items-center justify-center w-full mx-auto max-w-xs space-y-6 py-6 px-4 sm:px-0"),
			html.H2(
				attr.Class("font-medium text-xl text-center"),
				html.Text("Create an account"),
			),
			ui.Alert(ui.AlertProps{
				Variant:     ui.AlertDefault,
				Icon:        ui.SVGIcon(ui.IconInfo, "h-4 w-4"),
				Title:       "Hypercode is in early development.",
				Description: "Please reach out to the team if you encounter any issues.",
			}),
			html.Div(
				attr.Class("w-full space-y-4"),
				ui.GitHubAuthButton(),
				html.Div(
					attr.Class("relative"),
					html.Div(
						attr.Class("absolute inset-0 flex items-center"),
						html.Div(attr.Class("w-full border-t border-border")),
					),
					html.Div(
						attr.Class("relative flex justify-center text-xs uppercase"),
						html.Element("span",
							attr.Class("bg-neutral-50 px-2 text-muted-foreground"),
							html.Text("Or continue with"),
						),
					),
				),
				html.Form(
					attr.Method("POST"),
					attr.Action("/auth/sign-up"),
					attr.Class("space-y-4"),
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
					ui.FormField(ui.FormFieldProps{
						Label:       "Email Address",
						Id:          "email",
						Name:        "email",
						Type:        "email",
						Placeholder: "john@doe.com",
						Icon:        ui.IconMail,
						Required:    true,
						Value:       data.Email,
						Error:       data.EmailError,
					}),
					ui.FormField(ui.FormFieldProps{
						Label:       "Password",
						Id:          "password",
						Name:        "password",
						Type:        "password",
						Placeholder: "••••••••••••••••",
						Icon:        ui.IconLock,
						Required:    true,
						Error:       data.PasswordError,
					}),
					ui.Button(
						ui.ButtonProps{
							Variant: ui.ButtonPrimary,
							Type:    "submit",
							Class:   "w-full",
						},
						html.Text("Submit"),
					),
				),
			),
		),
	)
}

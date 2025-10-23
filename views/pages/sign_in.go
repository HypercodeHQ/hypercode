package pages

import (
	"net/http"

	"github.com/hypercodehq/hypercode/html"
	"github.com/hypercodehq/hypercode/html/attr"
	"github.com/hypercodehq/hypercode/views/components/layouts"
	"github.com/hypercodehq/hypercode/views/components/ui"
)

type SignInData struct {
	Error string
}

func SignIn(r *http.Request, data *SignInData) html.Node {
	if data == nil {
		data = &SignInData{}
	}

	return layouts.Main(r,
		"Sign in",
		html.Main(
			attr.Class("min-h-[calc(100vh-61px)] flex flex-col items-center justify-center w-full mx-auto max-w-xs space-y-6 py-6 px-4 sm:px-0"),
			html.H2(
				attr.Class("font-medium text-xl text-center"),
				html.Text("Sign in"),
			),
			ui.Alert(ui.AlertProps{
				Variant:     ui.AlertDefault,
				Icon:        ui.SVGIcon(ui.IconInfo, "h-4 w-4"),
				Title:       "Hypercode is in early development.",
				Description: "Please reach out to the team if you encounter any issues.",
			}),
			html.If(data.Error != "", ui.Alert(ui.AlertProps{
				Variant:     ui.AlertDestructive,
				Icon:        ui.SVGIcon(ui.IconAlertCircle, "h-4 w-4"),
				Title:       "Error",
				Description: data.Error,
			})),
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
						html.Span(
							attr.Class("bg-neutral-50 px-2 text-muted-foreground"),
							html.Text("Or continue with"),
						),
					),
				),
				html.Form(
					attr.Method("POST"),
					attr.Action("/auth/sign-in"),
					attr.Class("space-y-4"),
					ui.FormField(ui.FormFieldProps{
						Label:       "Email Address",
						Id:          "email",
						Name:        "email",
						Type:        "email",
						Placeholder: "john@doe.com",
						Icon:        ui.IconMail,
						Required:    true,
					}),
					ui.FormField(ui.FormFieldProps{
						Label:       "Password",
						Id:          "password",
						Name:        "password",
						Type:        "password",
						Placeholder: "••••••••••••••••",
						Icon:        ui.IconLock,
						Required:    true,
					}),
					ui.Button(
						ui.ButtonProps{
							Variant: ui.ButtonPrimary,
							Type:    "submit",
							Class:   "w-full",
						},
						html.Text("Sign In"),
					),
				),
			),
		),
	)
}

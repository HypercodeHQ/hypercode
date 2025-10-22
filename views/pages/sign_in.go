package pages

import (
	"net/http"

	"github.com/hyperstitieux/hypercode/html"
	"github.com/hyperstitieux/hypercode/html/attr"
	"github.com/hyperstitieux/hypercode/views/components/layouts"
	"github.com/hyperstitieux/hypercode/views/components/ui"
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
			html.Form(
				attr.Method("POST"),
				attr.Action("/sign-in"),
				attr.Class("w-full space-y-4"),
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
	)
}

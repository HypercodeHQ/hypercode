package pages

import (
	"net/http"

	"github.com/hypercodehq/hypercode/database/models"
	"github.com/hypercodehq/libhtml"
	"github.com/hypercodehq/libhtml/attr"
	"github.com/hypercodehq/hypercode/views/components/layouts"
	"github.com/hypercodehq/hypercode/views/components/ui"
)

type DeviceAuthData struct {
	User    *models.User
	Code    string
	Success bool
	Error   string
}

func DeviceAuth(r *http.Request, data *DeviceAuthData) html.Node {
	if data == nil {
		data = &DeviceAuthData{}
	}

	var content html.Node

	if data.Success {
		// Show success message
		content = html.Div(
			attr.Class("max-w-md mx-auto text-center"),
			html.Div(
				attr.Class("mb-6 flex justify-center"),
				html.Div(
					attr.Class("rounded-full bg-emerald-100 p-4"),
					ui.SVGIcon(ui.IconCheck, "size-12 text-emerald-600"),
				),
			),
			html.H2(
				attr.Class("text-2xl font-semibold mb-3"),
				html.Text("Device Connected!"),
			),
			html.P(
				attr.Class("text-muted-foreground mb-6 text-pretty"),
				html.Text("Your device has been successfully authenticated. You can now close this window and return to your terminal."),
			),
		)
	} else if data.User != nil {
		// Show confirmation form
		content = html.Div(
			attr.Class("max-w-md mx-auto"),
			html.Div(
				attr.Class("text-center mb-8"),
				html.H2(
					attr.Class("text-2xl font-semibold mb-3"),
					html.Text("Authenticate Hypercode CLI"),
				),
				html.P(
					attr.Class("text-muted-foreground"),
					html.Text("You are signing in as "),
					html.Span(
						attr.Class("font-semibold text-foreground"),
						html.Text("@"+data.User.Username),
					),
				),
			),

			html.If(data.Error != "", html.Div(
				attr.Class("mb-6 p-4 rounded-lg bg-red-50 border border-red-200 text-red-800 text-sm"),
				html.Text(data.Error),
			)),

			html.Form(
				attr.Method("POST"),
				attr.Action("/auth/device/confirm"),
				attr.Class("space-y-6"),

				html.Div(
					attr.Class("space-y-2"),
					html.Label(
						attr.For("code"),
						attr.Class("label text-center block"),
						html.Text("Enter the code shown in your terminal:"),
					),
					html.Input(
						attr.Type("text"),
						attr.Id("code"),
						attr.Name("code"),
						attr.Value(data.Code),
						attr.Placeholder("XXXX-XXXX"),
						attr.Required(),
						attr.Class("input w-full text-center text-2xl font-mono tracking-wider uppercase"),
						attr.Attribute{Key: "maxlength", Value: "9"}, // 8 chars + 1 hyphen
						attr.Attribute{Key: "pattern", Value: "[A-Z0-9]{4}-[A-Z0-9]{4}"},
						attr.Attribute{Key: "autocomplete", Value: "off"},
					),
					html.P(
						attr.Class("text-xs text-muted-foreground text-center"),
						html.Text("The code should look like: ABCD-1234"),
					),
				),

				ui.Button(ui.ButtonProps{
					Variant: ui.ButtonPrimary,
					Type:    "submit",
					Class:   "flex-1 w-full",
				}, html.Text("Confirm")),
			),

			// JavaScript to format code input
			html.Script(
				html.Text(`
					const codeInput = document.getElementById('code');
					if (codeInput) {
						codeInput.addEventListener('input', (e) => {
							let value = e.target.value.toUpperCase().replace(/[^A-Z0-9]/g, '');
							if (value.length > 4) {
								value = value.slice(0, 4) + '-' + value.slice(4, 8);
							}
							e.target.value = value;
						});
						codeInput.focus();
					}
				`),
			),
		)
	} else {
		// Show sign-in prompt
		content = html.Div(
			attr.Class("max-w-md mx-auto text-center"),
			html.Div(
				attr.Class("mb-6"),
				ui.SVGIcon(ui.IconLock, "size-12 text-muted-foreground mx-auto"),
			),
			html.H2(
				attr.Class("text-2xl font-semibold mb-3"),
				html.Text("Sign In Required"),
			),
			html.P(
				attr.Class("text-muted-foreground mb-6"),
				html.Text("You need to be signed in to authenticate your device."),
			),
			ui.Button(ui.ButtonProps{
				Variant: ui.ButtonPrimary,
			},
				html.Element("a",
					attr.Href("/auth/sign-in"),
					html.Text("Sign In"),
				),
			),
		)
	}

	return layouts.Main(r,
		"Device Authentication",
		html.Main(
			attr.Class("min-h-[calc(100vh-61px)] flex items-center justify-center px-4 py-8"),
			content,
		),
	)
}

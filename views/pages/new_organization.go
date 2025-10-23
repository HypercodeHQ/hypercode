package pages

import (
	"net/http"

	"github.com/hypercodehq/libhtml"
	"github.com/hypercodehq/libhtml/attr"
	"github.com/hypercodehq/hypercode/views/components/layouts"
	"github.com/hypercodehq/hypercode/views/components/ui"
)

type NewOrganizationData struct {
	Username         string
	DisplayName      string
	UsernameError    string
	DisplayNameError string
}

func NewOrganization(r *http.Request, data *NewOrganizationData) html.Node {
	if data == nil {
		data = &NewOrganizationData{}
	}

	return layouts.Main(r,
		"Create a new organization",
		html.Main(
			attr.Class("min-h-[calc(100vh-61px)] flex flex-col items-center justify-center w-full mx-auto max-w-sm space-y-8 py-6 px-4 sm:px-0"),
			html.H2(
				attr.Class("font-medium text-xl text-center"),
				html.Text("Create a new organization"),
			),
			html.Form(
				attr.Method("POST"),
				attr.Action("/organizations/new"),
				attr.Class("space-y-6 w-full"),
				ui.FormField(ui.FormFieldProps{
					Label:       "Organization username",
					Id:          "username",
					Name:        "username",
					Type:        "text",
					Placeholder: "acme-corp",
					Icon:        ui.IconBuilding,
					Required:    true,
					Value:       data.Username,
					Error:       data.UsernameError,
				}),
				ui.FormField(ui.FormFieldProps{
					Label:       "Display name",
					Id:          "display_name",
					Name:        "display_name",
					Type:        "text",
					Placeholder: "ACME Corporation",
					Icon:        ui.IconUser,
					Required:    true,
					Value:       data.DisplayName,
					Error:       data.DisplayNameError,
				}),
				ui.Button(ui.ButtonProps{
					Variant: ui.ButtonPrimary,
					Type:    "submit",
					Class:   "w-full",
				},
					html.Text("Create organization"),
				),
			),
			html.Script(
				html.Text(`
						const usernameInput = document.getElementById("username");
						usernameInput.addEventListener("input", (e) => {
							const cursorPos = e.target.selectionStart;
							const originalLength = e.target.value.length;
							e.target.value = e.target.value
								.toLowerCase()
								.replace(/\s+/g, "-")
								.replace(/[^a-z0-9\-_.]/g, "");
							const newLength = e.target.value.length;
							const diff = newLength - originalLength;
							e.target.setSelectionRange(cursorPos + diff, cursorPos + diff);
						});
					`),
			),
		),
	)
}

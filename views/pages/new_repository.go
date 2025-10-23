package pages

import (
	"net/http"

	"github.com/hypercodehq/hypercode/database/models"
	"github.com/hypercodehq/libhtml"
	"github.com/hypercodehq/libhtml/attr"
	"github.com/hypercodehq/hypercode/views/components/layouts"
	"github.com/hypercodehq/hypercode/views/components/ui"
)

type NewRepositoryData struct {
	Name               string
	DefaultBranch      string
	Visibility         string
	Owner              string
	NameError          string
	DefaultBranchError string
	User               *models.User
	Organizations      []*models.Organization
}

func NewRepository(r *http.Request, data *NewRepositoryData) html.Node {
	if data == nil {
		data = &NewRepositoryData{
			DefaultBranch: "main",
			Visibility:    "public",
		}
	}

	// Determine default owner
	defaultOwner := ""
	if data.Owner != "" {
		defaultOwner = data.Owner
	} else if data.User != nil {
		defaultOwner = data.User.Username
	}

	return layouts.Main(r,
		"Create a new repository",
		html.Main(
			attr.Class("min-h-[calc(100vh-61px)] flex flex-col items-center justify-center w-full mx-auto max-w-sm space-y-8 py-6 px-4 sm:px-0"),
			html.H2(
				attr.Class("font-medium text-xl text-center"),
				html.Text("Create a new repository"),
			),
			html.Form(
				attr.Method("POST"),
				attr.Action("/repositories/new"),
				attr.Class("space-y-4 w-full"),

				ownerSelector(data.User, data.Organizations, defaultOwner),
				ui.FormField(ui.FormFieldProps{
					Label:        "Repository name",
					Id:           "name",
					Name:         "name",
					Type:         "text",
					Placeholder:  "my-awesome-project",
					Icon:         ui.IconRepository,
					Required:     true,
					Value:        data.Name,
					Error:        data.NameError,
					WrapperClass: "sm:col-span-2",
				}),
				ui.FormField(ui.FormFieldProps{
					Label:       "Default branch",
					Id:          "default_branch",
					Name:        "default_branch",
					Type:        "text",
					Placeholder: "main",
					Icon:        ui.IconGitBranch,
					Required:    true,
					Value:       data.DefaultBranch,
					Error:       data.DefaultBranchError,
				}),
				html.Div(
					attr.Class("space-y-2"),
					html.Label(
						attr.Class("label"),
						html.Text("Visibility"),
					),
					html.Div(
						attr.Class("space-y-2"),
						html.Div(
							attr.Class("flex items-center space-x-2 p-3 border rounded-lg bg-white hover:bg-muted/50 transition-all cursor-pointer"),
							html.Input(
								attr.Type("radio"),
								attr.Id("visibility-public"),
								attr.Name("visibility"),
								attr.Value("public"),
								html.If(data.Visibility == "public" || data.Visibility == "", attr.Checked()),
								attr.Class("h-4 w-4"),
							),
							html.Label(
								attr.For("visibility-public"),
								attr.Class("flex-1 cursor-pointer flex items-center gap-2"),
								ui.SVGIcon(ui.IconGlobe, "h-4 w-4 text-muted-foreground"),
								html.Div(
									attr.Class("flex flex-col"),
									html.Span(
										attr.Class("font-medium text-sm"),
										html.Text("Public"),
									),
									html.Span(
										attr.Class("text-xs text-muted-foreground"),
										html.Text("Anyone can view this repository"),
									),
								),
							),
						),
						html.Div(
							attr.Class("flex items-center space-x-2 p-3 border rounded-lg bg-white hover:bg-muted/50 transition-all cursor-pointer"),
							html.Input(
								attr.Type("radio"),
								attr.Id("visibility-private"),
								attr.Name("visibility"),
								attr.Value("private"),
								html.If(data.Visibility == "private", attr.Checked()),
								attr.Class("h-4 w-4"),
							),
							html.Label(
								attr.For("visibility-private"),
								attr.Class("flex-1 cursor-pointer flex items-center gap-2"),
								ui.SVGIcon(ui.IconLock, "h-4 w-4 text-muted-foreground"),
								html.Div(
									attr.Class("flex flex-col"),
									html.Span(
										attr.Class("font-medium text-sm"),
										html.Text("Private"),
									),
									html.Span(
										attr.Class("text-xs text-muted-foreground"),
										html.Text("Only you and collaborators can access"),
									),
								),
							),
						),
					),
				),
				ui.Button(ui.ButtonProps{
					Variant: ui.ButtonPrimary,
					Type:    "submit",
					Class:   "w-full",
				},
					html.Text("Create repository"),
				),
			),
			html.Script(
				html.Text(`
						const nameInput = document.getElementById("name");
						nameInput.addEventListener("input", (e) => {
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

						const branchInput = document.getElementById("default_branch");
						branchInput.addEventListener("input", (e) => {
							const cursorPos = e.target.selectionStart;
							const originalLength = e.target.value.length;
							e.target.value = e.target.value
								.replace(/\s+/g, "-")
								.replace(/\.\./g, ".")
								.replace(/\/\//g, "/")
								.replace(/[~^:?*[\\\]@{}]/g, "")
								.toLowerCase();
							const newLength = e.target.value.length;
							const diff = newLength - originalLength;
							e.target.setSelectionRange(cursorPos + diff, cursorPos + diff);
						});
						branchInput.addEventListener("blur", (e) => {
							e.target.value = e.target.value
								.replace(/^[/.\-]+/, "")
								.replace(/[/.\-]+$/, "");
						});
					`),
			),
		),
	)
}

func ownerSelector(user *models.User, organizations []*models.Organization, selectedOwner string) html.Node {
	if user == nil {
		return html.Group()
	}

	// Determine if user is selected
	userSelected := selectedOwner == user.Username || selectedOwner == ""

	// Build select options
	selectOptions := []html.Node{
		// Current user option
		html.Element("option",
			attr.Value(user.Username),
			html.If(userSelected, attr.Selected(true)),
			html.Text(user.Username+" (You)"),
		),
	}

	// Add organizations under an optgroup if any exist
	if len(organizations) > 0 {
		orgOptions := []html.Node{}
		for _, org := range organizations {
			orgOptions = append(orgOptions, html.Element("option",
				attr.Value(org.Username),
				html.If(selectedOwner == org.Username, attr.Selected(true)),
				html.Text(org.Username+" ("+org.DisplayName+")"),
			))
		}

		selectOptions = append(selectOptions, html.Element("optgroup",
			attr.Label("Organizations"),
			html.Group(orgOptions...),
		))
	}

	return html.Div(
		attr.Class("space-y-2"),
		html.Label(
			attr.For("owner"),
			attr.Class("label"),
			html.Text("Owner"),
		),
		html.Element("select",
			attr.Id("owner"),
			attr.Name("owner"),
			attr.Required(),
			attr.Class("select w-full"),
			html.Group(selectOptions...),
		),
	)
}

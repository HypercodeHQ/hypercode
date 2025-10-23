package pages

import (
	"net/http"

	"github.com/hyperstitieux/hypercode/database/models"
	"github.com/hyperstitieux/hypercode/html"
	"github.com/hyperstitieux/hypercode/html/attr"
	"github.com/hyperstitieux/hypercode/views/components/layouts"
	"github.com/hyperstitieux/hypercode/views/components/ui"
)

type RepositorySettingsData struct {
	User               *models.User
	Repository         *models.Repository
	OwnerUsername      string
	Name               string
	DefaultBranch      string
	Visibility         string
	NameError          string
	DefaultBranchError string
	VisibilityError    string
	GeneralSuccess     string
	DangerZoneSuccess  string
}

func RepositorySettings(r *http.Request, data *RepositorySettingsData) html.Node {
	if data == nil {
		data = &RepositorySettingsData{}
	}

	// Populate form values from repository data if not set
	if data.Repository != nil {
		if data.Name == "" {
			data.Name = data.Repository.Name
		}
		if data.DefaultBranch == "" {
			data.DefaultBranch = data.Repository.DefaultBranch
		}
		if data.Visibility == "" {
			data.Visibility = data.Repository.Visibility
		}
	}

	return layouts.Repository(r,
		"Settings - "+data.OwnerUsername+"/"+data.Repository.Name,
		layouts.RepositoryLayoutOptions{
			OwnerUsername: data.OwnerUsername,
			RepoName:      data.Repository.Name,
			CurrentTab:    "settings",
			IsPublic:      data.Repository.Visibility == "public",
			ShowSettings:  true,
			StarCount:     0,
			HasStarred:    false,
		},
		html.Main(
			attr.Class("w-full mx-auto max-w-6xl space-y-6 py-8 px-4"),
			html.H1(
				attr.Class("font-semibold text-2xl mb-6"),
				html.Text("Repository Settings"),
			),

			// General Settings Card
			ui.Card(ui.CardProps{
				Title:       "General",
				Description: "Update repository information",
				Content: html.Div(
					attr.Class("space-y-4"),
					html.If(data.GeneralSuccess != "", html.Div(
						attr.Class("p-3 rounded-lg bg-emerald-50 dark:bg-emerald-900/20 border border-emerald-200 dark:border-emerald-800 text-emerald-800 dark:text-emerald-200 text-sm"),
						html.Text(data.GeneralSuccess),
					)),
					html.Form(
						attr.Id("general-settings-form"),
						attr.Method("POST"),
						attr.Action("/"+data.OwnerUsername+"/"+data.Repository.Name+"/settings/general"),
						attr.Class("space-y-4"),
						attr.Attribute{Key: "data-original-name", Value: data.Repository.Name},
						ui.FormField(ui.FormFieldProps{
							Label:       "Repository Name",
							Id:          "name",
							Name:        "name",
							Type:        "text",
							Placeholder: "my-awesome-project",
							Icon:        ui.IconRepository,
							Required:    true,
							Value:       data.Name,
							Error:       data.NameError,
						}),
						ui.FormField(ui.FormFieldProps{
							Label:       "Default Branch",
							Id:          "default_branch",
							Name:        "default_branch",
							Type:        "text",
							Placeholder: "main",
							Icon:        ui.IconGitBranch,
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
										html.If(data.Visibility == "public", attr.Checked()),
										attr.Class("h-4 w-4"),
									),
									html.Label(
										attr.For("visibility-public"),
										attr.Class("flex-1 cursor-pointer flex items-center gap-2"),
										ui.SVGIcon(ui.IconGlobe, "h-4 w-4 text-muted-foreground"),
										html.Div(
											attr.Class("flex flex-col"),
											html.Element("span",
												attr.Class("font-medium text-sm"),
												html.Text("Public"),
											),
											html.Element("span",
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
											html.Element("span",
												attr.Class("font-medium text-sm"),
												html.Text("Private"),
											),
											html.Element("span",
												attr.Class("text-xs text-muted-foreground"),
												html.Text("Only you and collaborators can access"),
											),
										),
									),
								),
							),
						),
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

			// Danger Zone Card
			ui.Card(ui.CardProps{
				Title:       "Danger Zone",
				Description: "Irreversible and destructive actions",
				Content: html.Div(
					attr.Class("space-y-4"),
					html.Div(
						attr.Class("flex items-center justify-between p-4 border border-destructive/50 rounded-lg bg-destructive/5"),
						html.Div(
							attr.Class("flex-1"),
							html.Element("h3",
								attr.Class("font-medium text-sm"),
								html.Text("Delete this repository"),
							),
							html.P(
								attr.Class("text-xs text-muted-foreground mt-1"),
								html.Text("Once you delete a repository, there is no going back. Please be certain."),
							),
						),
						ui.Button(
							ui.ButtonProps{
								Variant: ui.ButtonDestructive,
								OnClick: "confirmDeleteRepository()",
							},
							html.Text("Delete Repository"),
						),
					),
				),
			}),

			// JavaScript for name change confirmation and delete confirmation
			html.Element("script",
				html.Text(`
					(function() {
						const form = document.getElementById('general-settings-form');
						if (form) {
							form.addEventListener('submit', function(e) {
								const originalName = form.getAttribute('data-original-name');
								const currentName = document.getElementById('name').value;

								if (originalName !== currentName) {
									const confirmed = window.confirm(
										'Are you sure you want to rename this repository from "' + originalName + '" to "' + currentName + '"?\n\n' +
										'This will change the repository URL and may break existing clones.'
									);

									if (!confirmed) {
										e.preventDefault();
										return false;
									}
								}
							});
						}
					})();

					function confirmDeleteRepository() {
						const confirmed = window.confirm(
							'Are you ABSOLUTELY sure you want to delete this repository?\n\n' +
							'This action CANNOT be undone. This will permanently delete the repository, all commits, and all collaborators will lose access.\n\n' +
							'Type DELETE in the next prompt to confirm.'
						);

						if (confirmed) {
							const confirmation = window.prompt('Please type DELETE to confirm:');
							if (confirmation === 'DELETE') {
								const form = document.createElement('form');
								form.method = 'POST';
								form.action = window.location.pathname + '/delete';
								document.body.appendChild(form);
								form.submit();
							}
						}
					}
				`),
			),
		),
	)
}

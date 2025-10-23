package pages

import (
	"net/http"

	"github.com/hypercodehq/hypercode/database/models"
	"github.com/hypercodehq/hypercode/html"
	"github.com/hypercodehq/hypercode/html/attr"
	"github.com/hypercodehq/hypercode/views/components/layouts"
	"github.com/hypercodehq/hypercode/views/components/ui"
)

type NewTicketData struct {
	User          *models.User
	Repository    *models.Repository
	OwnerUsername string
	Title         string
	Body          string
	TitleError    string
	BodyError     string
	CanManage     bool
	StarCount     int64
	HasStarred    bool
	CloneURL      string
	RepositoryURL string
}

func NewTicket(r *http.Request, data *NewTicketData) html.Node {
	if data == nil {
		data = &NewTicketData{}
	}

	return layouts.Repository(r,
		"New ticket - "+data.OwnerUsername+"/"+data.Repository.Name,
		layouts.RepositoryLayoutOptions{
			OwnerUsername: data.OwnerUsername,
			RepoName:      data.Repository.Name,
			CurrentTab:    "tickets",
			IsPublic:      data.Repository.Visibility == "public",
			ShowSettings:  data.CanManage,
			StarCount:     data.StarCount,
			HasStarred:    data.HasStarred,
			DefaultBranch: data.Repository.DefaultBranch,
			CloneURL:      data.CloneURL,
			RepositoryURL: data.RepositoryURL,
		},
		html.Main(
			attr.Class("container mx-auto px-4 py-8 max-w-7xl"),
			html.Div(
				attr.Class("space-y-6"),
				html.H1(
					attr.Class("text-2xl font-semibold"),
					html.Text("New ticket"),
				),

				html.Form(
					attr.Method("post"),
					attr.Action("/"+data.OwnerUsername+"/"+data.Repository.Name+"/tickets/new"),
					attr.Class("space-y-6"),

					// Title field
					ui.FormField(ui.FormFieldProps{
						Label:       "Title",
						Id:          "title",
						Name:        "title",
						Type:        "text",
						Placeholder: "Brief description of the issue",
						Icon:        ui.IconCircle,
						Required:    true,
						Value:       data.Title,
						Error:       data.TitleError,
					}),

					// Body field
					html.Div(
						attr.Class("space-y-2"),
						html.Label(
							attr.For("body"),
							attr.Class("label"),
							html.Text("Description"),
						),
						html.Textarea(
							attr.Id("body"),
							attr.Name("body"),
							attr.Class("textarea min-h-[200px]"),
							attr.Placeholder("Provide more details about the issue..."),
							html.Text(data.Body),
						),
						html.If(data.BodyError != "", html.P(
							attr.Class("text-sm text-destructive"),
							html.Text(data.BodyError),
						)),
					),

					// Submit button
					html.Div(
						attr.Class("flex gap-3"),
						html.Button(
							attr.Type("submit"),
							attr.Class("btn-primary"),
							html.Text("Create ticket"),
						),
						html.A(
							attr.Href("/"+data.OwnerUsername+"/"+data.Repository.Name+"/tickets"),
							attr.Class("btn-outline"),
							html.Text("Cancel"),
						),
					),
				),
			),
		),
	)
}

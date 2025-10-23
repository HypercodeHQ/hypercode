package pages

import (
	"fmt"
	"net/http"

	"github.com/hypercodehq/hypercode/database/models"
	"github.com/hypercodehq/hypercode/html"
	"github.com/hypercodehq/hypercode/html/attr"
	"github.com/hypercodehq/hypercode/views/components/layouts"
	"github.com/hypercodehq/hypercode/views/components/ui"
)

type ShowTicketData struct {
	User           *models.User
	Repository     *models.Repository
	OwnerUsername  string
	Ticket         *models.Ticket
	Author         *models.User
	Comments       []*models.TicketComment
	CommentAuthors map[int64]*models.User
	CanManage      bool
	StarCount      int64
	HasStarred     bool
	CloneURL       string
	RepositoryURL  string
}

func ShowTicket(r *http.Request, data *ShowTicketData) html.Node {
	if data == nil {
		data = &ShowTicketData{}
	}

	return layouts.Repository(r,
		fmt.Sprintf("#%d %s - Tickets", data.Ticket.Number, data.Ticket.Title),
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
			attr.Class("container mx-auto px-4 py-8 max-w-4xl"),
			html.Div(
				attr.Class("space-y-6"),
				// Ticket header
				html.Div(
					attr.Class("flex items-start justify-between gap-4"),
					html.Div(
						attr.Class("flex-1"),
						html.H1(
							attr.Class("text-2xl font-semibold"),
							html.Text(data.Ticket.Title),
							html.Span(
								attr.Class("text-muted-foreground font-normal ml-2"),
								html.Text(fmt.Sprintf("#%d", data.Ticket.Number)),
							),
						),
						html.Div(
							attr.Class("mt-2 flex items-center gap-2"),
							statusBadge(data.Ticket.Status),
							html.Text(fmt.Sprintf("opened %s", formatTime(data.Ticket.CreatedAt))),
						),
					),
					html.If(
						data.User != nil,
						closeReopenButton(data),
					),
				),

				// Ticket body
				html.If(
					data.Ticket.Body != nil && *data.Ticket.Body != "",
					html.Div(
						attr.Class("border rounded-sm p-6 bg-card"),
						html.Div(
							attr.Class("flex items-center gap-3 mb-4 pb-4 border-b"),
							html.Div(
								attr.Class("p-2 rounded-full bg-muted"),
								ui.SVGIcon(ui.IconUser, "size-5"),
							),
							html.If(
								data.Author != nil,
								html.Span(
									attr.Class("font-medium"),
									html.Text(data.Author.DisplayName),
								),
							),
						),
						html.Div(
							attr.Class("prose prose-sm max-w-none"),
							html.Text(*data.Ticket.Body),
						),
					),
				),

				// Comments
				renderComments(data),

				// Comment form
				html.If(
					data.User != nil,
					commentForm(data),
				),
			),
		),
	)
}

func statusBadge(status string) html.Node {
	icon := ui.IconCircle
	classes := "inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-xs font-medium"

	if status == "closed" {
		icon = ui.IconCheck
		classes += " bg-purple-100 text-purple-800 dark:bg-purple-900/20 dark:text-purple-300"
	} else {
		classes += " bg-green-100 text-green-800 dark:bg-green-900/20 dark:text-green-300"
	}

	return html.Span(
		attr.Class(classes),
		ui.SVGIcon(icon, "size-3"),
		html.Text(capitalizeFirst(status)),
	)
}

func closeReopenButton(data *ShowTicketData) html.Node {
	if data.Ticket.Status == "open" {
		return html.Form(
			attr.Method("post"),
			attr.Action(fmt.Sprintf("/%s/%s/tickets/%d/close", data.OwnerUsername, data.Repository.Name, data.Ticket.Number)),
			html.Button(
				attr.Type("submit"),
				attr.Class("btn-outline"),
				html.Text("Close ticket"),
			),
		)
	}

	return html.Form(
		attr.Method("post"),
		attr.Action(fmt.Sprintf("/%s/%s/tickets/%d/reopen", data.OwnerUsername, data.Repository.Name, data.Ticket.Number)),
		html.Button(
			attr.Type("submit"),
			attr.Class("btn-outline"),
			html.Text("Reopen ticket"),
		),
	)
}

func renderComments(data *ShowTicketData) html.Node {
	if len(data.Comments) == 0 {
		return html.Div()
	}

	commentNodes := make([]html.Node, len(data.Comments))
	for i, comment := range data.Comments {
		author := data.CommentAuthors[comment.AuthorID]
		commentNodes[i] = renderComment(comment, author)
	}

	return html.Div(
		attr.Class("space-y-4"),
		html.Group(commentNodes...),
	)
}

func renderComment(comment *models.TicketComment, author *models.User) html.Node {
	return html.Div(
		attr.Class("border rounded-sm p-6 bg-card"),
		html.Div(
			attr.Class("flex items-center gap-3 mb-4 pb-4 border-b"),
			html.Div(
				attr.Class("p-2 rounded-full bg-muted"),
				ui.SVGIcon(ui.IconUser, "size-5"),
			),
			html.If(
				author != nil,
				html.Span(
					attr.Class("font-medium"),
					html.Text(author.DisplayName),
				),
			),
			html.Span(
				attr.Class("text-sm text-muted-foreground ml-auto"),
				html.Text(formatTime(comment.CreatedAt)),
			),
		),
		html.Div(
			attr.Class("prose prose-sm max-w-none"),
			html.Text(comment.Body),
		),
	)
}

func commentForm(data *ShowTicketData) html.Node {
	return html.Div(
		attr.Class("border rounded-sm p-6 bg-card"),
		html.Form(
			attr.Method("post"),
			attr.Action(fmt.Sprintf("/%s/%s/tickets/%d/comments", data.OwnerUsername, data.Repository.Name, data.Ticket.Number)),
			attr.Class("space-y-4"),
			html.Textarea(
				attr.Name("body"),
				attr.Id("comment-body"),
				attr.Class("input min-h-[100px]"),
				attr.Placeholder("Leave a comment..."),
				attr.Required(),
			),
			html.Div(
				attr.Class("flex justify-end"),
				html.Button(
					attr.Type("submit"),
					attr.Class("btn-primary"),
					html.Text("Comment"),
				),
			),
		),
	)
}

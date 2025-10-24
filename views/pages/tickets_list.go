package pages

import (
	"fmt"
	"net/http"
	"time"

	html "github.com/hypercommithq/libhtml"
	"github.com/hypercommithq/hypercommit/database/models"
	"github.com/hypercommithq/hypercommit/views/components/layouts"
	"github.com/hypercommithq/hypercommit/views/components/ui"
	"github.com/hypercommithq/libhtml/attr"
)

type TicketsListData struct {
	User          *models.User
	Repository    *models.Repository
	OwnerUsername string
	Tickets       []*models.Ticket
	StatusFilter  string
	OpenCount     int64
	ClosedCount   int64
	CanManage     bool
	StarCount     int64
	HasStarred    bool
	CloneURL      string
	RepositoryURL string
}

func TicketsList(r *http.Request, data *TicketsListData) html.Node {
	if data == nil {
		data = &TicketsListData{}
	}

	return layouts.Repository(r,
		"Tickets - "+data.OwnerUsername+"/"+data.Repository.Name,
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
				// Header with New Ticket button
				html.Div(
					attr.Class("flex justify-between items-center"),
					html.H1(
						attr.Class("text-2xl font-semibold"),
						html.Text("Tickets"),
					),
					html.If(
						data.User != nil,
						html.A(
							attr.Href("/"+data.OwnerUsername+"/"+data.Repository.Name+"/tickets/new"),
							attr.Class("btn-primary inline-flex items-center gap-2"),
							ui.SVGIcon(ui.IconPlus, "size-4"),
							html.Text("New ticket"),
						),
					),
				),

				// Tickets card
				ui.Card(ui.CardProps{
					Class: "!pt-1",
					Content: html.Div(
						attr.Class("space-y-4"),
						// Filter tabs
						html.Div(
							attr.Class("flex flex-wrap items-center gap-4 -mx-6 px-6 border-b"),
							filterTab("open", data.StatusFilter, data.OpenCount, data.OwnerUsername, data.Repository.Name),
							filterTab("closed", data.StatusFilter, data.ClosedCount, data.OwnerUsername, data.Repository.Name),
						),

						// Tickets list
						html.Div(
							attr.Class("-mx-6 -mb-6"),
							renderTicketsList(data),
						),
					),
				}),
			),
		),
	)
}

func filterTab(status, currentStatus string, count int64, owner, repo string) html.Node {
	isActive := status == currentStatus
	href := fmt.Sprintf("/%s/%s/tickets?status=%s", owner, repo, status)

	icon := ui.IconCircle
	if status == "closed" {
		icon = ui.IconCheck
	}

	spanClasses := "btn-ghost inline-flex items-center gap-2"
	if isActive {
		spanClasses += " font-medium"
	} else {
		spanClasses += " text-muted-foreground"
	}

	borderClass := "border-transparent"
	if isActive {
		borderClass = "border-zinc-900"
	}

	return html.A(
		attr.Href(href),
		attr.Class("inline-flex mt-2 pb-2 border-b-2 transition-colors "+borderClass),
		html.Span(
			attr.Class(spanClasses),
			ui.SVGIcon(icon, "size-4"),
			html.Text(fmt.Sprintf("%s (%d)", capitalizeFirst(status), count)),
		),
	)
}

func renderTicketsList(data *TicketsListData) html.Node {
	if len(data.Tickets) == 0 {
		return html.Div(
			attr.Class("py-8"),
			ui.EmptyState(ui.EmptyStateProps{
				Icon:        ui.SVGIcon(ui.IconCircle, "size-6"),
				Title:       fmt.Sprintf("No %s tickets", data.StatusFilter),
				Description: fmt.Sprintf("There are no %s tickets for this repository.", data.StatusFilter),
				ShowAction:  false,
			}),
		)
	}

	ticketItems := make([]html.Node, len(data.Tickets))
	for i, ticket := range data.Tickets {
		ticketItems[i] = renderTicketItem(data.OwnerUsername, data.Repository.Name, ticket)
	}

	return html.Div(
		attr.Class("divide-y"),
		html.Group(ticketItems...),
	)
}

func renderTicketItem(owner, repo string, ticket *models.Ticket) html.Node {
	ticketURL := fmt.Sprintf("/%s/%s/tickets/%d", owner, repo, ticket.Number)

	statusIcon := ui.IconCircle
	statusColor := "text-green-600"
	if ticket.Status == "closed" {
		statusIcon = ui.IconCheck
		statusColor = "text-purple-600"
	}

	return html.Div(
		attr.Class("p-4 hover:bg-muted/50 transition-colors"),
		html.A(
			attr.Href(ticketURL),
			attr.Class("flex items-start gap-3"),
			html.Div(
				attr.Class(statusColor+" flex-shrink-0 mt-1"),
				ui.SVGIcon(statusIcon, "size-5"),
			),
			html.Div(
				attr.Class("flex-1 min-w-0"),
				html.Div(
					attr.Class("flex items-start gap-2"),
					html.H3(
						attr.Class("font-medium text-foreground hover:text-primary"),
						html.Text(ticket.Title),
					),
				),
				html.Div(
					attr.Class("mt-1 text-sm text-muted-foreground"),
					html.Text(fmt.Sprintf("#%d opened %s", ticket.Number, formatTime(ticket.CreatedAt))),
				),
			),
		),
	)
}

func formatTime(unixTimestamp int64) string {
	t := time.Unix(unixTimestamp, 0)
	now := time.Now()
	diff := now.Sub(t)

	if diff < time.Minute {
		return "just now"
	} else if diff < time.Hour {
		minutes := int(diff.Minutes())
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	} else if diff < 24*time.Hour {
		hours := int(diff.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	} else if diff < 30*24*time.Hour {
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	} else {
		return t.Format("Jan 2, 2006")
	}
}

func capitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return string(s[0]-32) + s[1:]
}

package layouts

import (
	"net/http"

	"github.com/hypercommithq/hypercommit/database/models"
	"github.com/hypercommithq/libhtml"
	"github.com/hypercommithq/libhtml/attr"
	"github.com/hypercommithq/hypercommit/middleware"
	"github.com/hypercommithq/hypercommit/services"
	"github.com/hypercommithq/hypercommit/views/components"
	"github.com/hypercommithq/hypercommit/views/components/ui"
)

type exploreLayout struct {
	title      string
	children   []html.Node
	user       *models.User
	currentTab string
}

type ExploreLayoutOptions struct {
	CurrentTab string
}

func Explore(r *http.Request, title string, opts ExploreLayoutOptions, children ...html.Node) exploreLayout {
	return exploreLayout{
		title:      title,
		children:   children,
		user:       middleware.GetUserFromContext(r),
		currentTab: opts.CurrentTab,
	}
}

func (b exploreLayout) Render(w http.ResponseWriter, r *http.Request) error {
	bodyChildren := []html.Node{
		attr.Class("bg-neutral-50 text-neutral-900"),
		// Add toaster container for toast notifications early in DOM
		html.Div(
			attr.Id("toaster"),
			attr.Class("toaster"),
		),
		components.MainHeader(&components.MainHeaderData{User: b.user, Class: "!bg-accent"}),
		html.Div(
			attr.Class("bg-background border-b px-4 pt-2 flex flex-wrap items-center gap-4"),
			ui.ExploreTabs(ui.ExploreTabsProps{
				CurrentTab: b.currentTab,
			}),
		),
	}
	bodyChildren = append(bodyChildren, b.children...)

	flash := middleware.GetFlashFromContext(r)
	if flash != nil && flash.Type == services.FlashCelebration {
		bodyChildren = append(bodyChildren, components.Celebration())
	}

	doc := html.Document(
		html.HTML(
			attr.Lang("en"),
			components.Head(b.title),
			html.Body(bodyChildren...),
		),
	)
	return doc.Render(w, r)
}

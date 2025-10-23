package layouts

import (
	"net/http"

	"github.com/hyperstitieux/hypercode/database/models"
	"github.com/hyperstitieux/hypercode/html"
	"github.com/hyperstitieux/hypercode/html/attr"
	"github.com/hyperstitieux/hypercode/middleware"
	"github.com/hyperstitieux/hypercode/services"
	"github.com/hyperstitieux/hypercode/views/components"
)

type profileLayout struct {
	title        string
	children     []html.Node
	user         *models.User
	username     string
	displayName  string
	isOrg        bool
	currentTab   string
	showSettings bool
}

type ProfileLayoutOptions struct {
	Username     string
	DisplayName  string
	IsOrg        bool
	CurrentTab   string
	ShowSettings bool
}

func Profile(r *http.Request, title string, opts ProfileLayoutOptions, children ...html.Node) profileLayout {
	return profileLayout{
		title:        title,
		children:     children,
		user:         middleware.GetUserFromContext(r),
		username:     opts.Username,
		displayName:  opts.DisplayName,
		isOrg:        opts.IsOrg,
		currentTab:   opts.CurrentTab,
		showSettings: opts.ShowSettings,
	}
}

func (b profileLayout) Render(w http.ResponseWriter, r *http.Request) error {
	bodyChildren := []html.Node{
		attr.Class("bg-neutral-50 text-neutral-900"),
		// Add toaster container for toast notifications early in DOM
		html.Div(
			attr.Id("toaster"),
			attr.Class("toaster"),
		),
		components.ProfileHeader(&components.ProfileHeaderData{
			User:         b.user,
			Username:     b.username,
			DisplayName:  b.displayName,
			IsOrg:        b.isOrg,
			CurrentTab:   b.currentTab,
			ShowSettings: b.showSettings,
		}),
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

package layouts

import (
	"net/http"

	"github.com/hypercodehq/hypercode/database/models"
	"github.com/hypercodehq/libhtml"
	"github.com/hypercodehq/libhtml/attr"
	"github.com/hypercodehq/hypercode/middleware"
	"github.com/hypercodehq/hypercode/services"
	"github.com/hypercodehq/hypercode/views/components"
)

type mainLayout struct {
	title    string
	children []html.Node
	user     *models.User
}

func Main(r *http.Request, title string, children ...html.Node) mainLayout {
	return mainLayout{
		title:    title,
		children: children,
		user:     middleware.GetUserFromContext(r),
	}
}

func (b mainLayout) Render(w http.ResponseWriter, r *http.Request) error {
	bodyChildren := []html.Node{
		attr.Class("bg-neutral-50 text-neutral-900"),
		// Add toaster container for toast notifications early in DOM
		html.Div(
			attr.Id("toaster"),
			attr.Class("toaster"),
		),
		components.MainHeader(&components.MainHeaderData{User: b.user}),
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

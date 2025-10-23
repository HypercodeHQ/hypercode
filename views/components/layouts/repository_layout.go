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

type repositoryLayout struct {
	title         string
	children      []html.Node
	user          *models.User
	ownerUsername string
	repoName      string
	currentTab    string
	isPublic      bool
	showSettings  bool
	starCount     int64
	hasStarred    bool
	defaultBranch string
	cloneURL      string
	repositoryURL string
}

type RepositoryLayoutOptions struct {
	OwnerUsername string
	RepoName      string
	CurrentTab    string
	IsPublic      bool
	ShowSettings  bool
	StarCount     int64
	HasStarred    bool
	DefaultBranch string
	CloneURL      string
	RepositoryURL string
}

func Repository(r *http.Request, title string, opts RepositoryLayoutOptions, children ...html.Node) repositoryLayout {
	return repositoryLayout{
		title:         title,
		children:      children,
		user:          middleware.GetUserFromContext(r),
		ownerUsername: opts.OwnerUsername,
		repoName:      opts.RepoName,
		currentTab:    opts.CurrentTab,
		isPublic:      opts.IsPublic,
		showSettings:  opts.ShowSettings,
		starCount:     opts.StarCount,
		hasStarred:    opts.HasStarred,
		defaultBranch: opts.DefaultBranch,
		cloneURL:      opts.CloneURL,
		repositoryURL: opts.RepositoryURL,
	}
}

func (b repositoryLayout) Render(w http.ResponseWriter, r *http.Request) error {
	bodyChildren := []html.Node{
		attr.Class("bg-neutral-50 text-neutral-900"),
		components.RepositoryHeader(&components.RepositoryHeaderData{
			User:          b.user,
			OwnerUsername: b.ownerUsername,
			RepoName:      b.repoName,
			IsPublic:      b.isPublic,
			CurrentTab:    b.currentTab,
			ShowSettings:  b.showSettings,
			StarCount:     b.starCount,
			HasStarred:    b.hasStarred,
			DefaultBranch: b.defaultBranch,
			CloneURL:      b.cloneURL,
			RepositoryURL: b.repositoryURL,
		}),
	}
	bodyChildren = append(bodyChildren, b.children...)

	flash := middleware.GetFlashFromContext(r)
	if flash != nil && flash.Type == services.FlashCelebration {
		bodyChildren = append(bodyChildren, components.Celebration())
	}

	// Add toaster container for toast notifications
	bodyChildren = append(bodyChildren, html.Div(
		attr.Id("toaster"),
		attr.Class("toaster"),
	))

	doc := html.Document(
		html.HTML(
			attr.Lang("en"),
			components.Head(b.title),
			html.Body(bodyChildren...),
		),
	)
	return doc.Render(w, r)
}

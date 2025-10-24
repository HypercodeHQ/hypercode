package pages

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	html "github.com/hypercommithq/libhtml"
	"github.com/hypercommithq/hypercommit/database/models"
	"github.com/hypercommithq/hypercommit/views/components/layouts"
	"github.com/hypercommithq/hypercommit/views/components/ui"
	"github.com/hypercommithq/libhtml/attr"
)

type RepositoryFileData struct {
	User          *models.User
	Repository    *models.Repository
	OwnerUsername string
	CanManage     bool
	StarCount     int64
	HasStarred    bool
	Branches      []string
	CurrentBranch string
	CurrentPath   string
	FileContent   string
}

func RepositoryFile(r *http.Request, data *RepositoryFileData) html.Node {
	if data == nil {
		data = &RepositoryFileData{}
	}

	cloneURL := "https://" + r.Host + "/" + data.OwnerUsername + "/" + data.Repository.Name
	repositoryURL := cloneURL

	return layouts.Repository(r,
		"Code - "+data.OwnerUsername+"/"+data.Repository.Name,
		layouts.RepositoryLayoutOptions{
			OwnerUsername: data.OwnerUsername,
			RepoName:      data.Repository.Name,
			CurrentTab:    "tree",
			IsPublic:      data.Repository.Visibility == "public",
			ShowSettings:  data.CanManage,
			StarCount:     data.StarCount,
			HasStarred:    data.HasStarred,
			DefaultBranch: data.Repository.DefaultBranch,
			CloneURL:      cloneURL,
			RepositoryURL: repositoryURL,
		},
		html.Main(
			attr.Class("container mx-auto px-4 py-8 max-w-7xl"),
			html.H1(
				attr.Class("font-semibold text-2xl mb-6"),
				html.Text("Code"),
			),
			renderFileView(data),
		),
	)
}

func renderFileView(data *RepositoryFileData) html.Node {
	// Build branch selector
	branchSelector := renderFileBranchSelector(data)

	// Build breadcrumb navigation
	breadcrumb := renderFilePathBreadcrumb(data)

	// Get filename from path
	filename := filepath.Base(data.CurrentPath)

	// Build file content view
	fileContent := renderFileContent(data.FileContent, filename)

	return html.Div(
		attr.Class("space-y-4"),
		// Branch selector
		branchSelector,
		// Breadcrumb
		breadcrumb,
		// File content
		fileContent,
	)
}

func renderFileBranchSelector(data *RepositoryFileData) html.Node {
	selectOptions := []ui.SelectOption{}

	// Add default branch first
	defaultBranch := data.Repository.DefaultBranch
	if defaultBranch != "" {
		isSelected := data.CurrentBranch == defaultBranch
		selectOptions = append(selectOptions, ui.SelectOption{
			Value:    defaultBranch,
			Label:    defaultBranch + " (default)",
			Selected: isSelected,
			Icon:     ui.IconGitBranch,
		})
	}

	// Add other branches
	for _, branch := range data.Branches {
		if branch == defaultBranch {
			continue
		}
		isSelected := data.CurrentBranch == branch
		selectOptions = append(selectOptions, ui.SelectOption{
			Value:    branch,
			Label:    branch,
			Selected: isSelected,
			Icon:     ui.IconGitBranch,
		})
	}

	return html.Div(
		attr.Class("flex items-center gap-2"),
		ui.Select(ui.SelectProps{
			Id:      "branch-selector",
			Name:    "branch",
			Class:   "!mb-0 min-w-48",
			Options: selectOptions,
		}),
		html.Script(
			html.Text(fmt.Sprintf(`
(function() {
	const selector = document.getElementById('branch-selector');
	if (selector) {
		selector.addEventListener('change', function() {
			const branch = this.value;
			const owner = %q;
			const repo = %q;
			const path = %q;

			let url = "/" + owner + "/" + repo + "/tree/" + branch;
			if (path) {
				url += "/" + path;
			}

			window.location.href = url;
		});
	}
})();
			`, data.OwnerUsername, data.Repository.Name, data.CurrentPath)),
		),
	)
}

func renderFilePathBreadcrumb(data *RepositoryFileData) html.Node {
	if data.CurrentPath == "" {
		return html.Div()
	}

	parts := strings.Split(data.CurrentPath, "/")
	breadcrumbItems := []html.Node{}

	// Root
	breadcrumbItems = append(breadcrumbItems,
		html.Element("li",
			attr.Class("inline-flex items-center gap-1.5"),
			html.Element("a",
				attr.Href(fmt.Sprintf("/%s/%s/tree/%s", data.OwnerUsername, data.Repository.Name, data.CurrentBranch)),
				attr.Class("hover:text-foreground transition-colors"),
				html.Text(data.Repository.Name),
			),
		),
	)

	// Path parts
	currentPath := ""
	for i, part := range parts {
		if currentPath != "" {
			currentPath += "/"
		}
		currentPath += part

		// Add separator
		breadcrumbItems = append(breadcrumbItems,
			html.Element("li",
				ui.SVGIcon(ui.IconChevronRight, "size-3.5"),
			),
		)

		if i == len(parts)-1 {
			// Last part (filename) - not a link
			breadcrumbItems = append(breadcrumbItems, html.Element("li",
				attr.Class("inline-flex items-center gap-1.5"),
				html.Element("span",
					attr.Class("text-foreground font-normal"),
					html.Text(part),
				),
			))
		} else {
			// Intermediate part (directory) - link
			breadcrumbItems = append(breadcrumbItems,
				html.Element("li",
					attr.Class("inline-flex items-center gap-1.5"),
					html.Element("a",
						attr.Href(fmt.Sprintf("/%s/%s/tree/%s/%s", data.OwnerUsername, data.Repository.Name, data.CurrentBranch, currentPath)),
						attr.Class("hover:text-foreground transition-colors"),
						html.Text(part),
					),
				),
			)
		}
	}

	return html.Element("ol",
		attr.Class("text-muted-foreground flex flex-wrap items-center gap-1.5 text-sm break-words sm:gap-2.5"),
		html.Group(breadcrumbItems...),
	)
}

func renderFileContent(content string, filename string) html.Node {
	lines := strings.Split(content, "\n")

	return html.Div(
		attr.Class("border rounded-sm bg-card overflow-hidden"),
		// File header
		html.Div(
			attr.Class("px-4 py-3 border-b bg-muted/30 flex items-center gap-2"),
			ui.SVGIcon(ui.IconFile, "size-4 text-muted-foreground"),
			html.Span(
				attr.Class("text-sm font-medium"),
				html.Text(filename),
			),
			html.Span(
				attr.Class("text-sm text-muted-foreground ml-auto"),
				html.Text(fmt.Sprintf("%d lines", len(lines))),
			),
		),
		// File content
		html.Div(
			attr.Class("overflow-x-auto bg-white p-4"),
			html.Element("pre",
				attr.Class("text-sm"),
				html.Element("code",
					html.Text(content),
				),
			),
		),
	)
}

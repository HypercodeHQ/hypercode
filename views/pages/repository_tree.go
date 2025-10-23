package pages

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/hypercodehq/hypercode/database/models"
	"github.com/hypercodehq/libhtml"
	"github.com/hypercodehq/libhtml/attr"
	"github.com/hypercodehq/hypercode/services"
	"github.com/hypercodehq/hypercode/views/components/layouts"
	"github.com/hypercodehq/hypercode/views/components/ui"
)

type RepositoryTreeData struct {
	User          *models.User
	Repository    *models.Repository
	OwnerUsername string
	CanManage     bool
	StarCount     int64
	HasStarred    bool
	Branches      []string
	CurrentBranch string
	CurrentPath   string
	Entries       []services.TreeEntry
	IsEmpty       bool
}

func RepositoryTree(r *http.Request, data *RepositoryTreeData) html.Node {
	if data == nil {
		data = &RepositoryTreeData{}
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
			renderTreeContent(data),
		),
	)
}

func renderTreeContent(data *RepositoryTreeData) html.Node {
	if data.IsEmpty {
		return html.Div(
			attr.Class("border rounded-sm p-8 bg-card text-center"),
			ui.EmptyState(
				ui.EmptyStateProps{
					Icon:        ui.SVGIcon(ui.IconRepository, "size-6"),
					Title:       "This repository is empty",
					Description: "Get started by pushing code to this repository.",
				},
			),
		)
	}

	// Build branch selector
	branchSelector := renderBranchSelector(data)

	// Build breadcrumb navigation
	breadcrumb := renderPathBreadcrumb(data)

	// Build file/folder list
	fileList := renderFileList(data)

	return html.Div(
		attr.Class("space-y-4"),
		// Branch selector
		branchSelector,
		// Breadcrumb
		breadcrumb,
		// File list
		fileList,
	)
}

func renderBranchSelector(data *RepositoryTreeData) html.Node {
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

func renderPathBreadcrumb(data *RepositoryTreeData) html.Node {
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
			// Last part - not a link
			breadcrumbItems = append(breadcrumbItems, html.Element("li",
				attr.Class("inline-flex items-center gap-1.5"),
				html.Element("span",
					attr.Class("text-foreground font-normal"),
					html.Text(part),
				),
			))
		} else {
			// Intermediate part - link
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

func renderFileList(data *RepositoryTreeData) html.Node {
	if len(data.Entries) == 0 {
		return html.Div(
			attr.Class("border rounded-sm p-8 bg-card text-center"),
			html.P(
				attr.Class("text-muted-foreground"),
				html.Text("This directory is empty."),
			),
		)
	}

	rows := []html.Node{}

	for _, entry := range data.Entries {
		rows = append(rows, renderFileListItem(data, entry))
	}

	return html.Div(
		attr.Class("border rounded-sm bg-card overflow-hidden"),
		html.Table(
			attr.Class("w-full"),
			html.Tbody(
				rows...,
			),
		),
	)
}

func renderFileListItem(data *RepositoryTreeData, entry services.TreeEntry) html.Node {
	var icon ui.Icon
	var entryURL string
	isFolder := entry.Type == "tree"

	if isFolder {
		icon = ui.IconFolder
		entryURL = fmt.Sprintf("/%s/%s/tree/%s/%s", data.OwnerUsername, data.Repository.Name, data.CurrentBranch, entry.Path)
	} else {
		icon = ui.IconFile
		// For now, files also link to tree (in future, they should show file content)
		entryURL = fmt.Sprintf("/%s/%s/tree/%s/%s", data.OwnerUsername, data.Repository.Name, data.CurrentBranch, entry.Path)
	}

	return html.Tr(
		attr.Class("border-b last:border-b-0 hover:bg-muted/50 transition-colors"),
		html.Td(
			attr.Class("p-3"),
			html.Element("a",
				attr.Href(entryURL),
				attr.Class("flex items-center gap-3 text-foreground hover:text-primary transition-colors"),
				html.Div(
					attr.Class("flex-shrink-0 text-muted-foreground"),
					ui.SVGIcon(icon, "size-4"),
				),
				html.Span(
					attr.Class("text-sm"),
					html.Text(entry.Name),
				),
			),
		),
	)
}

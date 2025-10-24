package pages

import (
	"net/http"

	html "github.com/hypercommithq/libhtml"
	"github.com/hypercommithq/hypercommit/database/models"
	"github.com/hypercommithq/hypercommit/views/components/layouts"
	"github.com/hypercommithq/hypercommit/views/components/ui"
	"github.com/hypercommithq/libhtml/attr"
)

type ShowRepositoryData struct {
	User          *models.User
	Repository    *models.Repository
	OwnerUsername string
	CloneURL      string
	IsPublic      bool
	CanManage     bool
	StarCount     int64
	HasStarred    bool
}

func ShowRepository(r *http.Request, data *ShowRepositoryData) html.Node {
	if data == nil {
		data = &ShowRepositoryData{}
	}

	repositoryURL := "https://" + r.Host + "/" + data.OwnerUsername + "/" + data.Repository.Name

	return layouts.Repository(r,
		"Overview - "+data.OwnerUsername+"/"+data.Repository.Name,
		layouts.RepositoryLayoutOptions{
			OwnerUsername: data.OwnerUsername,
			RepoName:      data.Repository.Name,
			CurrentTab:    "overview",
			IsPublic:      data.IsPublic,
			ShowSettings:  data.CanManage,
			StarCount:     data.StarCount,
			HasStarred:    data.HasStarred,
			DefaultBranch: data.Repository.DefaultBranch,
			CloneURL:      data.CloneURL,
			RepositoryURL: repositoryURL,
		},
		html.Main(
			attr.Class("container mx-auto px-4 py-8 max-w-7xl"),
			html.H1(
				attr.Class("font-semibold text-2xl mb-6"),
				html.Text("Overview"),
			),
			html.Div(
				attr.Class("border rounded-sm p-6 bg-card"),
				html.Label(
					attr.For("clone-url"),
					attr.Class("label mb-2"),
					html.Text("Clone Command"),
				),
				html.Div(
					attr.Class("flex flex-wrap gap-2"),
					html.Input(
						attr.Type("text"),
						attr.Readonly(),
						attr.Value("git clone "+data.CloneURL),
						attr.Id("clone-url"),
						attr.Class("input font-mono text-sm flex-1"),
					),
					html.Button(
						attr.Onclick("copyCloneURL()"),
						attr.Id("copy-button"),
						attr.Class("btn-icon-outline cursor-pointer"),
						attr.DataTooltip("Copy to clipboard"),
						attr.DataSide("top"),
						ui.SVGIcon(ui.IconCopy, "size-4"),
						ui.SVGIcon(ui.IconCheck, "size-4 hidden"),
					),
				),
			),
			html.Script(
				html.Text(`function copyCloneURL() {
	const input = document.getElementById("clone-url");
	const cloneCommand = input.value;
	const copyButton = document.getElementById("copy-button");
	const copyIcon = copyButton.querySelector("svg:nth-child(1)");
	const checkIcon = copyButton.querySelector("svg:nth-child(2)");

	if (navigator.clipboard && navigator.clipboard.writeText) {
		navigator.clipboard.writeText(cloneCommand).then(() => {
			showSuccess();
		}).catch((err) => {
			console.error("Failed to copy:", err);
			fallbackCopy();
		});
	} else {
		fallbackCopy();
	}

	function fallbackCopy() {
		input.select();
		input.setSelectionRange(0, 99999);
		try {
			document.execCommand("copy");
			showSuccess();
		} catch (err) {
			console.error("Fallback copy failed:", err);
		}
	}

	function showSuccess() {
		copyButton.setAttribute("data-tooltip", "Copied!");
		copyIcon.classList.add("hidden");
		checkIcon.classList.remove("hidden");
		setTimeout(() => {
			copyButton.setAttribute("data-tooltip", "Copy to clipboard");
			copyIcon.classList.remove("hidden");
			checkIcon.classList.add("hidden");
		}, 2000);

		// Show toast notification
		document.dispatchEvent(new CustomEvent('basecoat:toast', {
			detail: {
				config: {
					category: 'success',
					title: 'Clone command copied!',
					description: 'The clone command has been copied to your clipboard.',
					duration: 2000
				}
			}
		}));
	}
}`),
			),
		),
	)
}

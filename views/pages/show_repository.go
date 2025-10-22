package pages

import (
	"net/http"

	"github.com/hyperstitieux/hypercode/database/models"
	"github.com/hyperstitieux/hypercode/html"
	"github.com/hyperstitieux/hypercode/html/attr"
	"github.com/hyperstitieux/hypercode/views/components/layouts"
	"github.com/hyperstitieux/hypercode/views/components/ui"
)

type ShowRepositoryData struct {
	User          *models.User
	Repository    *models.Repository
	OwnerUsername string
	CloneURL      string
	IsPublic      bool
}

func ShowRepository(r *http.Request, data *ShowRepositoryData) html.Node {
	if data == nil {
		data = &ShowRepositoryData{}
	}

	return layouts.Main(r, data.OwnerUsername+"/"+data.Repository.Name,
		html.Main(
			attr.Class("container mx-auto px-4 py-8 max-w-6xl"),
			html.Div(
				attr.Class("flex flex-wrap gap-4 items-center justify-between"),
				html.Div(
					attr.Class("flex flex-wrap gap-4 items-center"),
					html.H2(
						attr.Class("sm:text-xl font-medium"),
						html.Text(data.OwnerUsername+"/"+data.Repository.Name),
					),
					html.Div(
						ui.Badge(
							ui.BadgeProps{
								Variant: ui.BadgeOutline,
								Class:   "bg-card",
							},
							html.Text(visibilityText(data.IsPublic)),
						),
					),
				),
			),
			html.Div(
				attr.Class("border rounded-sm p-6 bg-card mt-6"),
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
					html.Element("button",
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
			html.Element("script",
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
	}
}`),
			),
		),
	)
}

func visibilityText(isPublic bool) string {
	if isPublic {
		return "Public"
	}
	return "Private"
}

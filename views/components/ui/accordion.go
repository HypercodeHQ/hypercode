package ui

import (
	"github.com/hypercodehq/hypercode/html"
	"github.com/hypercodehq/hypercode/html/attr"
)

type AccordionItemProps struct {
	Title   string
	Content html.Node
}

type AccordionProps struct {
	Items []AccordionItemProps
	Class string
}

func Accordion(props AccordionProps) html.Node {
	className := "accordion"
	if props.Class != "" {
		className += " " + props.Class
	}

	return html.Element("section",
		attr.Class(className),
		html.For(props.Items, func(item AccordionItemProps) html.Node {
			return html.Element("details",
				attr.Class("group border-b last:border-b-0"),
				html.Element("summary",
					attr.Class("w-full focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px] transition-all outline-none rounded-md"),
					html.H2(
						attr.Class("flex flex-1 items-start justify-between gap-4 py-4 text-left text-sm font-medium hover:underline"),
						html.Text(item.Title),
						SVGIcon(IconChevronDown, "text-muted-foreground pointer-events-none size-4 shrink-0 translate-y-0.5 transition-transform duration-200 group-open:rotate-180"),
					),
				),
				html.Element("section",
					attr.Class("pb-4"),
					item.Content,
				),
			)
		}),
	)
}

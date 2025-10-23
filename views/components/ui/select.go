package ui

import (
	"github.com/hyperstitieux/hypercode/html"
	"github.com/hyperstitieux/hypercode/html/attr"
)

type SelectOption struct {
	Value    string
	Label    string
	Selected bool
	Icon     Icon
}

type SelectProps struct {
	Id       string
	Name     string
	Label    string
	Options  []SelectOption
	Required bool
	Error    string
	Class    string
}

func Select(props SelectProps) html.Node {
	// Find selected option
	selectedLabel := ""
	selectedValue := ""
	selectedIcon := Icon("")
	for _, opt := range props.Options {
		if opt.Selected {
			selectedLabel = opt.Label
			selectedValue = opt.Value
			selectedIcon = opt.Icon
			break
		}
	}
	if selectedLabel == "" && len(props.Options) > 0 {
		selectedLabel = props.Options[0].Label
		selectedValue = props.Options[0].Value
		selectedIcon = props.Options[0].Icon
	}

	selectId := props.Id
	popoverId := props.Id + "-popover"
	listboxId := props.Id + "-listbox"

	labelAttrs := []html.Node{
		attr.For(selectId),
		attr.Class("label"),
		html.Text(props.Label),
	}

	wrapperClass := "space-y-2"
	if props.Class != "" {
		wrapperClass += " " + props.Class
	}

	return html.Div(
		attr.Class(wrapperClass),
		html.Label(labelAttrs...),
		html.Div(
			attr.Class("select !mb-0"),
			// Hidden input for form submission
			html.Input(
				attr.Type("hidden"),
				attr.Id(selectId),
				attr.Name(props.Name),
				attr.Value(selectedValue),
				html.If(props.Required, attr.Required()),
			),
			// Trigger button
			html.Element("button",
				attr.Type("button"),
				attr.Class("btn-outline w-full justify-between"),
				attr.AriaHaspopup("listbox"),
				attr.AriaControls(listboxId),
				attr.AriaExpanded("false"),
				html.Div(
					attr.Class("flex items-center gap-2"),
					html.Element("span",
						attr.Id(selectId+"-icon"),
						attr.Class("flex items-center"),
						html.If(selectedIcon != "", SVGIcon(selectedIcon, "h-4 w-4")),
					),
					html.Element("span",
						attr.Id(selectId+"-value"),
						html.Text(selectedLabel),
					),
				),
				SVGIcon(IconChevronDown, "h-4 w-4 opacity-50"),
			),
			// Popover
			html.Div(
				attr.Id(popoverId),
				attr.DataPopover(""),
				attr.AriaHidden("true"),
				attr.Class("w-full"),
				html.Div(
					attr.Role("listbox"),
					attr.Id(listboxId),
					attr.Class("max-h-64 overflow-y-auto scrollbar"),
					html.For(props.Options, func(option SelectOption) html.Node {
						return html.Div(
							attr.Role("option"),
							attr.Attribute{Key: "data-value", Value: option.Value},
							html.If(option.Icon != "", attr.Attribute{Key: "data-icon", Value: string(option.Icon)}),
							html.If(option.Selected, attr.Attribute{Key: "aria-selected", Value: "true"}),
							attr.Class("cursor-pointer flex items-center gap-2"),
							html.If(option.Icon != "", SVGIcon(option.Icon, "h-4 w-4")),
							html.Text(option.Label),
						)
					}),
				),
			),
		),
		html.If(props.Error != "", html.P(
			attr.Class("text-sm text-destructive"),
			html.Text(props.Error),
		)),
		// JavaScript for select functionality
		html.Element("script",
			html.Text(`
				(function() {
					const selectId = '`+selectId+`';
					const select = document.querySelector('.select:has(#' + selectId + ')');
					if (!select || select.hasAttribute('data-select-initialized')) return;
					select.setAttribute('data-select-initialized', 'true');

					const trigger = select.querySelector('[aria-haspopup="listbox"]');
					const popover = select.querySelector('[data-popover]');
					const listbox = select.querySelector('[role="listbox"]');
					const hiddenInput = select.querySelector('#' + selectId);
					const valueSpan = select.querySelector('#' + selectId + '-value');
					const iconSpan = select.querySelector('#' + selectId + '-icon');

					if (!trigger || !popover || !listbox || !hiddenInput || !valueSpan) return;

					trigger.addEventListener('click', () => {
						const isExpanded = trigger.getAttribute('aria-expanded') === 'true';
						trigger.setAttribute('aria-expanded', !isExpanded);
						popover.setAttribute('aria-hidden', isExpanded);
					});

					listbox.querySelectorAll('[role="option"]').forEach(option => {
						option.addEventListener('click', () => {
							const value = option.getAttribute('data-value');
							const optionIcon = option.querySelector('svg');
							const textContent = Array.from(option.childNodes)
								.filter(node => node.nodeType === Node.TEXT_NODE)
								.map(node => node.textContent.trim())
								.join('');

							// Update hidden input
							hiddenInput.value = value;

							// Update display
							valueSpan.textContent = textContent;

							// Update icon
							if (iconSpan) {
								if (optionIcon) {
									const clonedIcon = optionIcon.cloneNode(true);
									iconSpan.innerHTML = '';
									iconSpan.appendChild(clonedIcon);
								} else {
									iconSpan.innerHTML = '';
								}
							}

							// Update aria-selected
							listbox.querySelectorAll('[role="option"]').forEach(opt => {
								opt.removeAttribute('aria-selected');
							});
							option.setAttribute('aria-selected', 'true');

							// Close popover
							trigger.setAttribute('aria-expanded', 'false');
							popover.setAttribute('aria-hidden', 'true');
						});
					});

					// Close on outside click
					document.addEventListener('click', (e) => {
						if (!select.contains(e.target)) {
							trigger.setAttribute('aria-expanded', 'false');
							popover.setAttribute('aria-hidden', 'true');
						}
					});
				})();
			`),
		),
	)
}

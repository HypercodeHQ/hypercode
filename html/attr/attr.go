package attr

import (
	"net/http"

	"github.com/hyperstitieux/hypercode/html"
)

type Attribute struct {
	Key   string
	Value string
}

func (a Attribute) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (a Attribute) GetAttribute() (key, value string) {
	return a.Key, a.Value
}

func Lang(value string) html.Node {
	return Attribute{Key: "lang", Value: value}
}

func For(value string) html.Node {
	return Attribute{Key: "for", Value: value}
}

func Charset(value string) html.Node {
	return Attribute{Key: "charset", Value: value}
}

func Type(value string) html.Node {
	return Attribute{Key: "type", Value: value}
}

func Name(value string) html.Node {
	return Attribute{Key: "name", Value: value}
}

func Id(value string) html.Node {
	return Attribute{Key: "id", Value: value}
}

func Href(value string) html.Node {
	return Attribute{Key: "href", Value: value}
}

func Src(value string) html.Node {
	return Attribute{Key: "src", Value: value}
}

func Alt(value string) html.Node {
	return Attribute{Key: "alt", Value: value}
}

func Class(value string) html.Node {
	return Attribute{Key: "class", Value: value}
}

func Rel(value string) html.Node {
	return Attribute{Key: "rel", Value: value}
}

func Content(value string) html.Node {
	return Attribute{Key: "content", Value: value}
}

func Method(value string) html.Node {
	return Attribute{Key: "method", Value: value}
}

func Action(value string) html.Node {
	return Attribute{Key: "action", Value: value}
}

func Placeholder(value string) html.Node {
	return Attribute{Key: "placeholder", Value: value}
}

func Required() html.Node {
	return Attribute{Key: "required", Value: ""}
}

func Checked() html.Node {
	return Attribute{Key: "checked", Value: ""}
}

func Value(value string) html.Node {
	return Attribute{Key: "value", Value: value}
}

func Readonly() html.Node {
	return Attribute{Key: "readonly", Value: ""}
}

func Disabled() html.Node {
	return Attribute{Key: "disabled", Value: ""}
}

func Xmlns(value string) html.Node {
	return Attribute{Key: "xmlns", Value: value}
}

func Width(value string) html.Node {
	return Attribute{Key: "width", Value: value}
}

func Height(value string) html.Node {
	return Attribute{Key: "height", Value: value}
}

func ViewBox(value string) html.Node {
	return Attribute{Key: "viewBox", Value: value}
}

func Fill(value string) html.Node {
	return Attribute{Key: "fill", Value: value}
}

func Stroke(value string) html.Node {
	return Attribute{Key: "stroke", Value: value}
}

func StrokeWidth(value string) html.Node {
	return Attribute{Key: "stroke-width", Value: value}
}

func StrokeLinecap(value string) html.Node {
	return Attribute{Key: "stroke-linecap", Value: value}
}

func StrokeLinejoin(value string) html.Node {
	return Attribute{Key: "stroke-linejoin", Value: value}
}

func D(value string) html.Node {
	return Attribute{Key: "d", Value: value}
}

func Cx(value string) html.Node {
	return Attribute{Key: "cx", Value: value}
}

func Cy(value string) html.Node {
	return Attribute{Key: "cy", Value: value}
}

func R(value string) html.Node {
	return Attribute{Key: "r", Value: value}
}

func X(value string) html.Node {
	return Attribute{Key: "x", Value: value}
}

func Y(value string) html.Node {
	return Attribute{Key: "y", Value: value}
}

func Rx(value string) html.Node {
	return Attribute{Key: "rx", Value: value}
}

func Ry(value string) html.Node {
	return Attribute{Key: "ry", Value: value}
}

func X1(value string) html.Node {
	return Attribute{Key: "x1", Value: value}
}

func X2(value string) html.Node {
	return Attribute{Key: "x2", Value: value}
}

func Y1(value string) html.Node {
	return Attribute{Key: "y1", Value: value}
}

func Y2(value string) html.Node {
	return Attribute{Key: "y2", Value: value}
}

func Points(value string) html.Node {
	return Attribute{Key: "points", Value: value}
}

func As(value string) html.Node {
	return Attribute{Key: "as", Value: value}
}

func Defer() html.Node {
	return Attribute{Key: "defer", Value: ""}
}

func Onclick(value string) html.Node {
	return Attribute{Key: "onclick", Value: value}
}

func Target(value string) html.Node {
	return Attribute{Key: "target", Value: value}
}

func AriaLabel(value string) html.Node {
	return Attribute{Key: "aria-label", Value: value}
}

func AriaHaspopup(value string) html.Node {
	return Attribute{Key: "aria-haspopup", Value: value}
}

func AriaControls(value string) html.Node {
	return Attribute{Key: "aria-controls", Value: value}
}

func AriaExpanded(value string) html.Node {
	return Attribute{Key: "aria-expanded", Value: value}
}

func AriaHidden(value string) html.Node {
	return Attribute{Key: "aria-hidden", Value: value}
}

func AriaLabelledby(value string) html.Node {
	return Attribute{Key: "aria-labelledby", Value: value}
}

func Role(value string) html.Node {
	return Attribute{Key: "role", Value: value}
}

func DataTooltip(value string) html.Node {
	return Attribute{Key: "data-tooltip", Value: value}
}

func DataSide(value string) html.Node {
	return Attribute{Key: "data-side", Value: value}
}

func DataPopover(value string) html.Node {
	return Attribute{Key: "data-popover", Value: value}
}

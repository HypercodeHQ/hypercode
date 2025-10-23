package html

import (
	"net/http"
	"strings"
)

type element struct {
	tag         string
	children    []Node
	attrs       map[string]string
	selfClosing bool
}

func (e element) Render(w http.ResponseWriter, r *http.Request) error {
	var sb strings.Builder

	sb.WriteString("<")
	sb.WriteString(e.tag)

	for key, value := range e.attrs {
		sb.WriteString(" ")
		sb.WriteString(key)
		sb.WriteString(`="`)
		sb.WriteString(value)
		sb.WriteString(`"`)
	}

	if e.selfClosing {
		sb.WriteString(">")
		_, err := w.Write([]byte(sb.String()))
		return err
	}

	sb.WriteString(">")

	if _, err := w.Write([]byte(sb.String())); err != nil {
		return err
	}

	for _, child := range e.children {
		if err := child.Render(w, r); err != nil {
			return err
		}
	}

	closeTag := "</" + e.tag + ">"
	_, err := w.Write([]byte(closeTag))
	return err
}

func createElement(tag string, children ...Node) Node {
	attrs := make(map[string]string)
	var actualChildren []Node

	for _, child := range children {
		if attr, ok := child.(Attributer); ok {
			key, value := attr.GetAttribute()
			attrs[key] = value
		} else {
			actualChildren = append(actualChildren, child)
		}
	}

	return element{
		tag:      tag,
		children: actualChildren,
		attrs:    attrs,
	}
}

func HTML(children ...Node) Node {
	return createElement("html", children...)
}

func Head(children ...Node) Node {
	return createElement("head", children...)
}

func Body(children ...Node) Node {
	return createElement("body", children...)
}

func Title(children ...Node) Node {
	return createElement("title", children...)
}

func P(children ...Node) Node {
	return createElement("p", children...)
}

func Div(children ...Node) Node {
	return createElement("div", children...)
}

func Main(children ...Node) Node {
	return createElement("main", children...)
}

func Form(children ...Node) Node {
	return createElement("form", children...)
}

func Header(children ...Node) Node {
	return createElement("header", children...)
}

func Footer(children ...Node) Node {
	return createElement("footer", children...)
}

func Section(children ...Node) Node {
	return createElement("section", children...)
}

func H1(children ...Node) Node {
	return createElement("h1", children...)
}

func H2(children ...Node) Node {
	return createElement("h2", children...)
}

func H3(children ...Node) Node {
	return createElement("h3", children...)
}

func Ul(children ...Node) Node {
	return createElement("ul", children...)
}

func Li(children ...Node) Node {
	return createElement("li", children...)
}

func A(children ...Node) Node {
	return createElement("a", children...)
}

func Button(children ...Node) Node {
	return createElement("button", children...)
}

func Span(children ...Node) Node {
	return createElement("span", children...)
}

func Select(children ...Node) Node {
	return createElement("select", children...)
}

func Option(children ...Node) Node {
	return createElement("option", children...)
}

func Textarea(children ...Node) Node {
	return createElement("textarea", children...)
}

func Script(children ...Node) Node {
	return createElement("script", children...)
}

func Table(children ...Node) Node {
	return createElement("table", children...)
}

func Tbody(children ...Node) Node {
	return createElement("tbody", children...)
}

func Tr(children ...Node) Node {
	return createElement("tr", children...)
}

func Td(children ...Node) Node {
	return createElement("td", children...)
}

func Details(children ...Node) Node {
	return createElement("details", children...)
}

func Summary(children ...Node) Node {
	return createElement("summary", children...)
}

func Nav(children ...Node) Node {
	return createElement("nav", children...)
}

func Hr(children ...Node) Node {
	return createElement("hr", children...)
}

func Label(children ...Node) Node {
	return createElement("label", children...)
}

func Element(tag string, children ...Node) Node {
	return createElement(tag, children...)
}

func Input(children ...Node) Node {
	attrs := make(map[string]string)

	for _, child := range children {
		if attr, ok := child.(Attributer); ok {
			key, value := attr.GetAttribute()
			attrs[key] = value
		} else if cond, ok := child.(conditional); ok {
			if cond.condition && cond.node != nil {
				if attr, ok := cond.node.(Attributer); ok {
					key, value := attr.GetAttribute()
					attrs[key] = value
				}
			}
		}
	}

	return element{
		tag:         "input",
		attrs:       attrs,
		selfClosing: true,
	}
}

func Img(children ...Node) Node {
	attrs := make(map[string]string)

	for _, child := range children {
		if attr, ok := child.(Attributer); ok {
			key, value := attr.GetAttribute()
			attrs[key] = value
		}
	}

	return element{
		tag:         "img",
		attrs:       attrs,
		selfClosing: true,
	}
}

func Meta(children ...Node) Node {
	attrs := make(map[string]string)

	for _, child := range children {
		if attr, ok := child.(Attributer); ok {
			key, value := attr.GetAttribute()
			attrs[key] = value
		}
	}

	return element{
		tag:         "meta",
		attrs:       attrs,
		selfClosing: true,
	}
}

func Link(children ...Node) Node {
	attrs := make(map[string]string)

	for _, child := range children {
		if attr, ok := child.(Attributer); ok {
			key, value := attr.GetAttribute()
			attrs[key] = value
		}
	}

	return element{
		tag:         "link",
		attrs:       attrs,
		selfClosing: true,
	}
}

type document struct {
	root Node
}

func (d document) Render(w http.ResponseWriter, r *http.Request) error {
	if _, err := w.Write([]byte("<!DOCTYPE html>")); err != nil {
		return err
	}
	return d.root.Render(w, r)
}

func Document(root Node) document {
	return document{root: root}
}

type conditional struct {
	condition bool
	node      Node
}

func (c conditional) Render(w http.ResponseWriter, r *http.Request) error {
	if c.condition {
		return c.node.Render(w, r)
	}
	return nil
}

func If(condition bool, node Node) Node {
	return conditional{condition: condition, node: node}
}

func IfElse(condition bool, thenNode Node, elseNode Node) Node {
	if condition {
		return thenNode
	}
	return elseNode
}

func Iff(condition bool, fn func() Node) Node {
	if condition {
		return fn()
	}
	return conditional{condition: false, node: nil}
}

func IfElsef(condition bool, thenFn func() Node, elseFn func() Node) Node {
	if condition {
		return thenFn()
	}
	return elseFn()
}

type forLoop struct {
	nodes []Node
}

func (f forLoop) Render(w http.ResponseWriter, r *http.Request) error {
	for _, node := range f.nodes {
		if err := node.Render(w, r); err != nil {
			return err
		}
	}
	return nil
}

func For[T any](items []T, fn func(T) Node) Node {
	nodes := make([]Node, 0, len(items))
	for _, item := range items {
		nodes = append(nodes, fn(item))
	}
	return forLoop{nodes: nodes}
}

type group struct {
	children []Node
}

func (g group) Render(w http.ResponseWriter, r *http.Request) error {
	for _, child := range g.children {
		if err := child.Render(w, r); err != nil {
			return err
		}
	}
	return nil
}

// Group renders multiple children without wrapping them in an HTML element.
// This is useful when you need to return multiple nodes as a single Node.
func Group(children ...Node) Node {
	return group{children: children}
}

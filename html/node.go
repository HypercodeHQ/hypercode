package html

import "net/http"

type Node interface {
	Render(w http.ResponseWriter, r *http.Request) error
}

type Attributer interface {
	Node
	GetAttribute() (key, value string)
}

type textNode struct {
	content string
}

func (t textNode) Render(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err := w.Write([]byte(t.content))
	return err
}

func Text(content string) Node {
	return textNode{content: content}
}

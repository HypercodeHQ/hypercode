package public

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed assets
var embedFS embed.FS

func FileServer() http.Handler {
	fsys, err := fs.Sub(embedFS, "assets")
	if err != nil {
		panic(err)
	}
	return http.FileServer(http.FS(fsys))
}

package middleware

import (
	"net/http"
	"strings"
)

var staticFileExtensions = []string{
	".css", ".js", ".svg", ".png", ".jpg", ".ico",
	".woff", ".woff2", ".ttf", ".webp", ".gif",
}

func StaticFileServer(fileServer http.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if isStaticFile(r.URL.Path) {
				fileServer.ServeHTTP(w, r)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func isStaticFile(path string) bool {
	if path == "/" || strings.Contains(path[1:], "/") {
		return false
	}

	for _, ext := range staticFileExtensions {
		if strings.HasSuffix(path, ext) {
			return true
		}
	}

	return false
}

package httputil

import "net/http"

// IsHTTPS checks if the request is over HTTPS, either directly or via a proxy
func IsHTTPS(r *http.Request) bool {
	// Check if request is directly HTTPS
	if r.TLS != nil {
		return true
	}

	// Check X-Forwarded-Proto header (set by reverse proxies like Caddy)
	proto := r.Header.Get("X-Forwarded-Proto")
	return proto == "https"
}

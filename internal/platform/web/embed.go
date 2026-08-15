package web

import (
	"embed"
	"io/fs"
	"net/http"
	"path"
	"strings"
)

//go:embed all:static
var static embed.FS

// Handler serves the embedded dashboard. Deep links fall back to index.html
// so client-side routing works (e.g. /files/abc123). API paths are never
// served here.
func Handler() (http.Handler, error) {
	sub, err := fs.Sub(static, "static")
	if err != nil {
		return nil, err
	}
	fileServer := http.FileServer(http.FS(sub))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
		if p == "" {
			p = "index.html"
		}
		if strings.HasPrefix(p, "v1/") || strings.HasPrefix(p, "s/") {
			http.NotFound(w, r)
			return
		}
		// Serve real files; fall back to index.html for SPA deep links.
		if _, err := fs.Stat(sub, p); err != nil {
			r2 := r.Clone(r.Context())
			r2.URL.Path = "/"
			fileServer.ServeHTTP(w, r2)
			return
		}
		fileServer.ServeHTTP(w, r)
	}), nil
}

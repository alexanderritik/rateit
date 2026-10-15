package web

import (
	"io/fs"
	"net/http"
)

// Handler returns an http.Handler that serves the embedded React build.
// All routes that don't match a file fall back to index.html (SPA support).
func Handler() http.Handler {
	dist, _ := fs.Sub(StaticFiles, "dist")
	fileServer := http.FileServer(http.FS(dist))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Try the file; if it doesn't exist, serve index.html
		f, err := dist.Open(r.URL.Path[1:])
		if err != nil {
			r.URL.Path = "/"
		} else {
			f.Close()
		}
		fileServer.ServeHTTP(w, r)
	})
}

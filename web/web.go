package web

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed static
var static embed.FS

func SetupStatic(mux *http.ServeMux) {
	// Serve static files from the embedded filesystem in production mode
	fsys := fs.FS(static)
	contentStatic, _ := fs.Sub(fsys, "static")
	staticFS := http.FS(contentStatic)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(staticFS)))
}

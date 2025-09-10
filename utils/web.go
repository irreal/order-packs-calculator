package utils

import (
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
)

func Render(w http.ResponseWriter, r *http.Request, component templ.Component) error {
	w.Header().Set("Content-Type", "text/html")
	err := component.Render(r.Context(), w)
	if err != nil {
		slog.Error("Error rendering component", "error", err)
		return err
	}
	return nil
}

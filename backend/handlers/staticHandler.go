package handlers

import (
	"net/http"
	"os"
)

// StaticHandler handles requests for static files...
func StaticHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		RenderError(w, http.StatusMethodNotAllowed)
		return
	}

	file, err := os.Stat(r.URL.Path[1:])
	if err != nil || file.IsDir() {
		RenderError(w, http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, r.URL.Path[1:])
}

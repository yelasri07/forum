package google

import (
	"database/sql"
	"fmt"
	"net/http"

	"forum/backend/handlers"
)

const (
	clientID      = "978621672489-dsv02jh60f8jbkqdd3khg2p7jtse3om5.apps.googleusercontent.com"
	client_secret = "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"
)

func RegisterGoogle(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodGet {
		handlers.RenderError(w, http.StatusMethodNotAllowed)
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		handlers.RenderError(w, http.StatusBadRequest)
		return
	}
	
	fmt.Println(code)
}

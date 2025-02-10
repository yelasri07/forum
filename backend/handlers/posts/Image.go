package posts

import (
	"database/sql"
	"net/http"
	"strconv"

	"forum/backend/handlers"
	"forum/backend/models"
)

func ImageHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodGet {
		handlers.RenderError(w, http.StatusMethodNotAllowed)
		return
	}

	idPost, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		handlers.RenderError(w, http.StatusBadRequest)
		return
	}

	image := models.GetImage(idPost, db)

	imageContentType := http.DetectContentType(image)

	w.Header().Set("Content-Type", imageContentType)
	w.Write(image)
}

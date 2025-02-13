package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"forum/backend/models"
)

func ImageHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodGet {
		RenderError(w, http.StatusMethodNotAllowed)
		return
	}

	idPost, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		RenderError(w, http.StatusBadRequest)
		return
	}

	image := models.GetImage(idPost, db)

	imageContentType := http.DetectContentType(image)

	if imageContentType == "text/xml; charset=utf-8" {
		imageContentType = "image/svg+xml"
	}

	w.Header().Set("Content-Type", imageContentType)
	w.Write(image)
}

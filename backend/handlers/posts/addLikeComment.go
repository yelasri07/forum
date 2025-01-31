package posts

import (
	"database/sql"
	"net/http"
	"strconv"

	"forum/backend/handlers"
	"forum/backend/models"
)

func AddlikeComment(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		handlers.RenderError(w, http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseForm()
	refer := r.Referer()
	if err != nil {
		handlers.RenderError(w, http.StatusBadRequest)
		return
	}
	userID := r.Context().Value("userId").(int)
	IdComment := r.FormValue("IDComment")

	NewIdComment, err := strconv.Atoi(IdComment)
	if err != nil || NewIdComment <= 0 {
		handlers.RenderError(w, http.StatusBadRequest)
		return
	}

	status := r.FormValue("status")

	boolen, err := models.CheckIdPost(db, NewIdComment, "Comment")
	if err != nil {
		handlers.RenderError(w, http.StatusInternalServerError)
		return
	}

	if !boolen {
		handlers.RenderError(w, http.StatusBadRequest)
		return
	}

	err = models.AddReactInTheComment(db, status, userID, NewIdComment)
	if err != nil {
		handlers.RenderError(w, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, refer, http.StatusSeeOther)
}

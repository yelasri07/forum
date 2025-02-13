package posts

import (
	"database/sql"
	"net/http"
	"strconv"

	"forum/backend/handlers"
	"forum/backend/models"
)

// AddlikeComment handles adding a "like" or "dislike" reaction to a comment.
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
	status := r.FormValue("status")

	NewIdComment, err := strconv.Atoi(IdComment)
	if err != nil || !models.CheckIdExists(db, NewIdComment, "Comment") || (status != "like" && status != "dislike") {
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

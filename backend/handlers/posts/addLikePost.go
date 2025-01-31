package posts

import (
	"database/sql"
	"net/http"
	"strconv"

	"forum/backend/handlers"
	"forum/backend/models"
)

func AddLikePost(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		handlers.RenderError(w, http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		handlers.RenderError(w, http.StatusBadRequest)
		return
	}

	status := r.FormValue("status")
	userID := r.Context().Value("userId").(int)
	postIDStr := r.FormValue("postID")
	refer := r.Referer()
	postID, err := strconv.Atoi(postIDStr)
	if err != nil || postID <= 0 {
		http.Error(w, "Invalid post ID format", http.StatusBadRequest)
		return
	}
	if status != "like" && status != "dislike" {
		handlers.RenderError(w, http.StatusBadRequest)
		return
	}

	if !models.CheckIdPost(db, postID, "Posts") {
		handlers.RenderError(w, http.StatusBadRequest)
		return
	}

	err = models.AddReact(db, status, userID, postID)
	if err != nil {
		handlers.RenderError(w, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, refer, http.StatusFound)
}

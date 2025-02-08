package posts

import (
	"database/sql"
	"net/http"
	"strconv"

	"forum/backend/handlers"
	"forum/backend/models"
)

// AddLikePost handles adding a "like" or "dislike" reaction to a post.
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
	refer := r.Referer()
	postID, err := strconv.Atoi(r.FormValue("postID"))
	if err != nil || !models.CheckIdExists(db, postID, "Posts") || (status != "like" && status != "dislike") {
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

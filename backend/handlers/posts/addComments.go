package posts

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"forum/backend/handlers"
	"forum/backend/models"
)

func AddComments(w http.ResponseWriter, r *http.Request, db *sql.DB) {
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

	comment := r.FormValue("comment")
	postIDStr := strings.Join(r.Form["postID"], "")

	postID, err := strconv.Atoi(postIDStr)
	if err != nil || postID <= 0 {
		handlers.RenderError(w, http.StatusBadRequest)
		return
	}

	if comment == "" {
		handlers.RenderError(w, http.StatusBadRequest)
		return
	}

	if !models.CheckIdPost(db, postID, "Posts") {
		handlers.RenderError(w, http.StatusBadRequest)
		return
	}

	err = models.AddCommentsInDB(db, userID, postID, comment)
	if err != nil {
		handlers.RenderError(w, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, refer, http.StatusSeeOther)
}

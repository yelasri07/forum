package posts

import (
	"database/sql"
	"net/http"

	"forum/backend/handlers"
	"forum/backend/models"
)

// FilterByLikedPosts handles the request to display posts liked by the user.
func FilterByLikedPosts(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodGet {
		handlers.RenderError(w, http.StatusMethodNotAllowed)
		return
	}
	var err error

	ID := r.Context().Value("userId").(int)
	userName := r.Context().Value("userName").(string)
	Homepage, err := handlers.GetDataHomePage(r, db, ID)
	if err != nil {
		handlers.RenderError(w, http.StatusInternalServerError)
		return
	}

	Homepage.UserName = userName

	Homepage.PostCat, err = models.GetPostsByLike(db, ID)
	if err != nil {
		handlers.RenderError(w, http.StatusInternalServerError)
		return
	}

	err = handlers.RenderTemplate(w, "index.html", Homepage, http.StatusOK)
	if err != nil {
		handlers.RenderError(w, http.StatusInternalServerError)
		return
	}
}

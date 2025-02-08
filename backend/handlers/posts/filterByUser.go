package posts

import (
	"database/sql"
	"net/http"

	"forum/backend/handlers"
	"forum/backend/models"
)

// FilterByUser handles the request to display all posts created by the current user.
func FilterByUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodGet {
		handlers.RenderError(w, http.StatusMethodNotAllowed)
		return
	}

	ID := r.Context().Value("userId").(int)
	userName := r.Context().Value("userName").(string)
	Homepage, err := handlers.GetDataHomePage(r, db, ID)
	if err != nil {
		handlers.RenderError(w, http.StatusInternalServerError)
		return
	}

	Homepage.PostCat, err = models.GetAllPostCatByUser(db, ID)
	if err != nil {
		handlers.RenderError(w, http.StatusInternalServerError)
		return
	}

	Homepage.UserName = userName

	err = handlers.RenderTemplate(w, "index.html", Homepage, http.StatusOK)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

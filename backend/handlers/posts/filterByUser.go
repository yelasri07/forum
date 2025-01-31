package posts

import (
	"database/sql"
	"net/http"

	"forum/backend/models"
	"forum/backend/handlers"
)

func FilterByUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	Homepage := handlers.GetDataHomePage(r, db)
	ID := r.Context().Value("userId").(int)
	var err error

	if !Homepage.IsLogged {
		handlers.RenderError(w, http.StatusUnauthorized)
		return
	}

	Homepage.PostCat, err = models.GetAllPostCatByUser(db, ID)

	if err != nil {
		handlers.RenderError(w, http.StatusInternalServerError)
		return
	}

	err = handlers.RenderTemplate(w, "index.html", Homepage, http.StatusOK)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

}

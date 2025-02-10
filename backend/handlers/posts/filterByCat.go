package posts

import (
	"database/sql"
	"net/http"

	"forum/backend/handlers"
	"forum/backend/models"
	"forum/middleware"
)

// FilterByCat handles the request to filter posts by selected categories.
func FilterByCat(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodGet {
		handlers.RenderError(w, http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		handlers.RenderError(w, http.StatusBadRequest)
		return
	}

	var idUser int
	var userName string
	token, err := middleware.VerifyCookie(r, db)
	if err == nil {
		idUser, userName = models.GetInfos(db, token.Value)
	}

	categories := r.Form["category"]
	if len(categories) == 0 {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	var post_ids []int
	for _, category := range categories {
		err = models.GetPostIDsByCategory(db, category, &post_ids)
		if err != nil {
			handlers.RenderError(w, http.StatusInternalServerError)
			return
		}
	}

	var filteredPosts []*models.PostCat
	for _, postID := range post_ids {
		p, err := models.GetPostByID(db, postID, idUser)
		if err != nil {
			handlers.RenderError(w, http.StatusInternalServerError)
			return
		}
		filteredPosts = append(filteredPosts, p)
	}

	Homepage, err := handlers.GetDataHomePage(r, db, idUser)
	if err != nil {
		handlers.RenderError(w, http.StatusInternalServerError)
		return
	}

	Homepage.PostCat = filteredPosts
	Homepage.UserName = userName

	err = handlers.RenderTemplate(w, "index.html", Homepage, http.StatusOK)
	if err != nil {
		handlers.RenderError(w, http.StatusInternalServerError)
		return
	}
}

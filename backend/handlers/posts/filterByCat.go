package posts

import (
	"database/sql"
	"net/http"

	"forum/backend/handlers"
	"forum/backend/models"
	"forum/middleware"
)

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
	Token, err := middleware.VerifyCookie(r, db)
	id_user := 0
	if err == nil {
		id_user, _ = models.GetInfos(db, Token.Value)
	}

	categories := r.Form["category"]
	if len(categories) == 0 {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	var post_ids []int
	for _, category := range categories {
		err := models.GetPostIDsByCategory(db, &post_ids, category)
		if err != nil {
			handlers.RenderError(w, http.StatusInternalServerError)
			return
		}
	}
	var filteredPosts []*models.PostCat
	for _, postID := range post_ids {
		p, err := models.GetPostByID(db, postID,id_user)
		if err != nil {
			handlers.RenderError(w, http.StatusInternalServerError)
			return
		}
		filteredPosts = append(filteredPosts, p)
	}
	Homepage := handlers.GetDataHomePage(r, db)
	Homepage.PostCat = filteredPosts
	err = handlers.RenderTemplate(w, "index.html", Homepage, http.StatusOK)
	if err != nil {
		handlers.RenderError(w, http.StatusInternalServerError)
		return
	}
}

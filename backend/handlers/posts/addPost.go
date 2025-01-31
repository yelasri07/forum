package posts

import (
	"database/sql"
	"net/http"
	"time"

	"forum/backend/handlers"
)

func AddPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		handlers.RenderError(w, http.StatusMethodNotAllowed)
		return
	}
	referer := r.Referer()
	err := r.ParseForm()
	if err != nil {
		handlers.RenderError(w, http.StatusBadRequest)
		return
	}

	ID := r.Context().Value("userId").(int)
	title := r.FormValue("title")
	content := r.FormValue("content")

	categories := r.Form["category"]

	if title == "" || content == "" || len(categories) == 0 {
		handlers.RenderError(w, http.StatusBadRequest)
		return
	}
	result, err := db.Exec("INSERT INTO Posts (Title, Content, DateCreation, ID_User) VALUES (?,?,?,?)", title, content, time.Now(), ID)
	if err != nil {
		handlers.RenderError(w, http.StatusInternalServerError)
		return
	}
	idPost, _ := result.LastInsertId()
	for _, categoryID := range categories {
		_, err := db.Exec("INSERT INTO PostCategory (ID_Post, ID_Category) VALUES (?, ?)", int(idPost), categoryID)
		if err != nil {
			handlers.RenderError(w, http.StatusInternalServerError)
			return
		}
	}
	http.Redirect(w, r, referer, http.StatusFound)
}

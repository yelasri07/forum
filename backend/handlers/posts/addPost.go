package posts

import (
	"database/sql"
	"io"
	"net/http"
	"strings"
	"time"

	"forum/backend/handlers"
	"forum/backend/models"
)

// AddPost handles the creation of a new post by a user.
func AddPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		handlers.RenderError(w, http.StatusMethodNotAllowed)
		return
	}

	r.ParseMultipartForm(10 << 20)
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
	file, fileHeader, err := r.FormFile("image")

	image := []byte(nil)
	if err == nil {
		defer file.Close()
		sizeInMB := float64(fileHeader.Size) / (1024 * 1024)
		if sizeInMB > 20 || (!strings.HasSuffix(fileHeader.Filename, ".jpeg") &&
			!strings.HasSuffix(fileHeader.Filename, ".svg") &&
			!strings.HasSuffix(fileHeader.Filename, ".png") &&
			!strings.HasSuffix(fileHeader.Filename, ".gif") &&
			!strings.HasSuffix(fileHeader.Filename, ".jpg")) {
			handlers.RenderError(w, http.StatusBadRequest)
			return
		}

		img, err := io.ReadAll(file)
		if err != nil {
			handlers.RenderError(w, http.StatusInternalServerError)
		}

		image = append(image, img...)
	}

	if title == "" || content == "" || len(categories) == 0 || len([]rune(content)) > 1000 || len([]rune(title)) > 50 || !models.CheckCatExists(categories, db) {
		handlers.RenderError(w, http.StatusBadRequest)
		return
	}

	result, err := db.Exec("INSERT INTO Posts (Title, Content, DateCreation, Image,ID_User) VALUES (?,?,?,?,?)", title, content, time.Now(), image, ID)
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

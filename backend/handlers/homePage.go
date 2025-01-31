package handlers

import (
	"database/sql"
	"net/http"

	"forum/backend/models"
	"forum/middleware"
)

func IndexPage(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.URL.Path != "/" {
		RenderError(w, http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		RenderError(w, http.StatusMethodNotAllowed)
		return
	}

	Homepage := GetDataHomePage(r, db)

	err := RenderTemplate(w, "index.html", Homepage, http.StatusOK)
	if err != nil {
		RenderError(w, http.StatusInternalServerError)
		return
	}
}

func GetDataHomePage(r *http.Request, db *sql.DB) *models.HomePage {
	Homepage := new(models.HomePage)

	Token, err := middleware.VerifyCookie(r, db)
	id_user := -1
	if err == nil {
		Homepage.IsLogged = true
		ID, userName := models.GetInfos(db, Token.Value)
		Homepage.UserName = userName
		totalLikes, _ := models.GetTotalLikesByUser(db, ID)
		Homepage.TotalLikes = totalLikes
		id_user = ID
	} else {
		Homepage.IsLogged = false
	}

	// Always fetch all posts
	Homepage.PostCat, _ = models.GetAllPostCat(db, id_user)
	Homepage.Categories, _ = models.GetAllCategories(db)

	return Homepage
}

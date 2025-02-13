package handlers

import (
	"database/sql"
	"net/http"

	"forum/backend/models"
	"forum/middleware"
)

// IndexPage handles the request for the homepage.
func IndexPage(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.URL.Path != "/" {
		RenderError(w, http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		RenderError(w, http.StatusMethodNotAllowed)
		return
	}

	var idUser int
	var userName string
	token, err := middleware.VerifyCookie(r, db)
	if err == nil {
		idUser, userName = models.GetInfos(db, token.Value)
	}

	Homepage, err := GetDataHomePage(r, db, idUser)
	if err != nil {
		RenderError(w, http.StatusInternalServerError)
		return
	}

	Homepage.UserName = userName

	Homepage.PostCat, err = models.GetAllPostCat(db, idUser)
	if err != nil {
		RenderError(w, http.StatusInternalServerError)
		return
	}

	err = RenderTemplate(w, "index.html", Homepage, http.StatusOK)
	if err != nil {
		RenderError(w, http.StatusInternalServerError)
		return
	}
}

// GetDataHomePage retrieves data for the homepage based on whether the user is logged in.
func GetDataHomePage(r *http.Request, db *sql.DB, idUser int) (*models.HomePage, error) {
	Homepage := new(models.HomePage)
	var err error

	if idUser != 0 {
		Homepage.IsLogged = true
		Homepage.TotalLikes, err = models.GetTotalLikesByUser(db, idUser)
		if err != nil {
			return nil, err
		}
	} else {
		Homepage.IsLogged = false
	}

	Homepage.Categories, err = models.GetAllCategories(db)
	if err != nil {
		return nil, err
	}

	return Homepage, nil
}

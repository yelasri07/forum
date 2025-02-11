package github

import (
	"database/sql"
	"net/http"

	"forum/backend/handlers"
	"forum/backend/models"
	"forum/utils"
)

func CallbackGithubLogin(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	code := r.URL.Query().Get("code")
	if code == "" {
		handlers.RenderError(w, http.StatusBadRequest)
		return
	}

	accessToken, err := getAccessToken(code)
	if err != nil {
		handlers.RenderError(w, http.StatusInternalServerError)
		return
	}
	user, err := getGitHubUser(accessToken)
	if err != nil {
		handlers.RenderError(w, http.StatusInternalServerError)
		return
	}

	if user.Email == "" {
		email, err := getPrimaryEmail(accessToken)
		if err != nil {
			handlers.RenderError(w, http.StatusInternalServerError)
			return
		}
		user.Email = email
	}
	ID, err := models.VerifyEmail(db, user.Email)
	if err != nil {
		handlers.RenderError(w, http.StatusServiceUnavailable)
		return
	}

	if ID == -1 || !utils.IsValidEmail(user.Email) {
		e := &models.ErrorRegister{AlreadyLogedWithGithub: "You don't have account try to Register."}
		handlers.RenderTemplate(w, "login.html", e, http.StatusConflict)
		return
	}
	token, err := models.GenerateToken(ID, db)
	if err != nil {
		handlers.RenderError(w, http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{Name: "Token", Value: token, MaxAge: 3600, HttpOnly: true}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

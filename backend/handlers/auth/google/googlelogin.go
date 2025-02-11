package google

import (
	"database/sql"
	"net/http"
	"os"

	"forum/backend/handlers"
	"forum/backend/models"
	"forum/utils"
)

func GoogleLogin(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	code := r.URL.Query().Get("code")
	if code == "" {
		handlers.RenderError(w, http.StatusBadRequest)
		return
	}
	accessToken, err := getAccesstoken(code, os.Getenv("redirect_uri_login"))
	if err != nil {
		handlers.RenderError(w, http.StatusInternalServerError)
		return
	}
	user, err := getGoogleUser(accessToken)
	if err != nil {
		handlers.RenderError(w, http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
		return
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

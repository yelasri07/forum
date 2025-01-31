package auth

import (
	"database/sql"
	"net/http"

	"forum/backend/handlers"
	"forum/backend/models"
	"forum/middleware"

	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		handlers.RenderError(w, http.StatusBadRequest)
		return
	}

	Email := r.FormValue("Email")
	Pass := r.FormValue("Password")

	ID, err := models.VerifyEmail(db, Email)
	if err != nil {
		handlers.RenderError(w, http.StatusInternalServerError)
		return
	}
	if ID == -1 {	
		e := models.ErrorRegister{ErrEmail: "Incorrect email"}
		handlers.RenderTemplate(w, "login.html", e, http.StatusConflict)
		return
	}

	PasswordDatabase, err := models.VerifyPassword(db, ID)
	if err != nil {
		handlers.RenderError(w, http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(PasswordDatabase), []byte(Pass))
	if err != nil {
		e := models.ErrorRegister{ErrPassword: "Incorrect Password"}
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

func LoginPage(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodGet {
		handlers.RenderError(w, http.StatusMethodNotAllowed)
		return
	}

	_, err := middleware.VerifyCookie(r, db)
	if err == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err = handlers.RenderTemplate(w, "login.html", nil, http.StatusOK)
	if err != nil {
		handlers.RenderError(w, http.StatusInternalServerError)
		return
	}
}

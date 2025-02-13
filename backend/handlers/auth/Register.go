package auth

import (
	"database/sql"
	"net/http"
	"time"

	"forum/backend/handlers"
	"forum/backend/models"
	"forum/middleware"
	"forum/utils"

	"golang.org/x/crypto/bcrypt"
)

// Register handles the user registration process, validates the inputs, checks for unique username and email,
// hashes the password, inserts the user into the database, generates a session token, and sets it as a secure cookie.
func Register(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/sign-up", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		handlers.RenderError(w, http.StatusBadRequest)
		return
	}
	UserName := r.FormValue("UserName")
	Email := r.FormValue("Email")

	isUniqueUserName, err := models.UserExists(db, UserName, " UserName ")
	if err != nil {
		handlers.RenderError(w, http.StatusServiceUnavailable)
		return
	}

	isUniqueEmail, err := models.UserExists(db, Email, " Email ")
	if err != nil {
		handlers.RenderError(w, http.StatusServiceUnavailable)
		return
	}

	if !Verify(w, isUniqueUserName, isUniqueEmail, Email, UserName, r.FormValue("Password")) {
		return
	}
	password, _ := bcrypt.GenerateFromPassword([]byte(r.FormValue("Password")), 10)

	result, err := db.Exec("INSERT INTO Users (UserName, Email, Password, Created_At, Session, Expared_At) VALUES ( ?,?,?,?,?,?)", UserName, Email, string(password), time.Now(), "", nil)
	if err != nil {
		handlers.RenderError(w, http.StatusInternalServerError)
		return
	}

	ID, err := result.LastInsertId()
	if err != nil {
		handlers.RenderError(w, http.StatusInternalServerError)
		return
	}

	token, err := models.GenerateToken(int(ID), db)
	if err != nil {
		handlers.RenderError(w, http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{Name: "Token", Value: token, MaxAge: 3600, HttpOnly: true}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Verify(w http.ResponseWriter, isUniqueUserName, isUniqueEmail bool, Email, UserName, Password string) bool {
	e := &models.ErrorRegister{}
	if Email == "" {
		e.ErrName = "Username cannot be emty"
	}
	if UserName == "" {
		e.ErrName = "Username cannot be emty"
	}
	if !utils.ValidName(UserName) {
		e.ErrName = "Username cannot conatains a special charachters like: \"@()-.,;...\" except: _"
	}
	if len([]rune(Password)) < 8 || len([]rune(Password)) > 20 {
		e.ErrPassword = "Password must be greater than 8 characters and less than 20 characters"
	}
	if !isUniqueUserName {
		e.ErrName = "Username Already taken please chose Another"
	}
	if !isUniqueEmail {
		e.ErrEmail = "Email Already taken please chose Another"
	}
	if !utils.IsValidEmail(Email) {
		e.ErrEmail = "Email must be in the format: example@example.example"
	}
	if e.ErrEmail != "" || e.ErrName != "" || e.ErrPassword != "" {
		handlers.RenderTemplate(w, "register.html", e, http.StatusConflict)
		return false
	}
	return true
}

func RegisterPage(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodGet {
		handlers.RenderError(w, http.StatusMethodNotAllowed)
		return
	}

	_, err := middleware.VerifyCookie(r, db)
	if err == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err = handlers.RenderTemplate(w, "register.html", nil, http.StatusOK)
	if err != nil {
		handlers.RenderError(w, http.StatusInternalServerError)
		return
	}
}

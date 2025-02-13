package auth

import (
	"database/sql"
	"net/http"
)

// Logout invalidates the user's session by clearing the session cookie and redirecting to the home page.
func Logout(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	cookie := &http.Cookie{Name: "Token", Value: "", MaxAge: -1, HttpOnly: true}

	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

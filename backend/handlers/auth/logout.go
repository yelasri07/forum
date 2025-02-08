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

	deleteCookie, err := r.Cookie("Token")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	deleteCookie.Value = ""
	deleteCookie.MaxAge = -1

	http.SetCookie(w, deleteCookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

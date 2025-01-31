package auth

import (
	"database/sql"
	"net/http"
)

// had lfonction dyal nmssho token men cookie
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

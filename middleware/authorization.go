package middleware

import (
	"context"
	"database/sql"
	"net/http"
	"time"
)
//middleware
func Authorization(next http.Handler, db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("Token")

		if err != nil {
			http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
			return
		}

		var userId int
		var expired time.Time
		db.QueryRow("SELECT ID, Expared_At FROM Users WHERE Session=?", cookie.Value).Scan(&userId, &expired)
		if userId == 0 {
			http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
			return
		}

		if time.Now().UTC().After(expired.UTC()) {
			db.Exec("UPDATE users set Session=? WHERE ID=?", "", userId)
			http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(r.Context(), "userId", userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

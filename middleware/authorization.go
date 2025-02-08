package middleware

import (
	"context"
	"database/sql"
	"net/http"
	"time"
)

// Authorization checks if the user is logged in by verifying the session token and expiration time 
// from the database. If valid, it attaches user details to the request context and proceeds to the next handler.
func Authorization(next http.Handler, db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("Token")
		if err != nil {
			http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
			return
		}

		var userId int
		var userName string
		var expired time.Time
		db.QueryRow("SELECT ID, UserName ,Expared_At FROM Users WHERE Session=?", cookie.Value).Scan(&userId, &userName, &expired)
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
		ctx = context.WithValue(ctx, "userName", userName)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

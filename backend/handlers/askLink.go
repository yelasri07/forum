package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"forum/backend/models"
)

func AskLink(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		RenderError(w, http.StatusMethodNotAllowed)
		return
	}

	token, err := r.Cookie("UserID")
	if err != nil {
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
		return
	}

	UserID := GetUserIdByToken(token.Value, db)

	UserIDFromBrowser, err := strconv.Atoi(r.FormValue("accept"))

	if err != nil || UserID != UserIDFromBrowser {
		RenderError(w, http.StatusBadRequest)
		return
	}

	cookie := &http.Cookie{Name: "UserID", Value: "", MaxAge: -1, HttpOnly: true}
	http.SetCookie(w, cookie)

	_, err = db.Exec("UPDATE users SET AuthType = ? WHERE ID = ?", 1, UserID)
	if err != nil {
		RenderError(w, http.StatusInternalServerError)
		return
	}

	token2, err := models.GenerateToken(UserID, db)
	if err != nil {
		RenderError(w, http.StatusInternalServerError)
		return
	}

	cookie = &http.Cookie{Name: "Token", Value: token2, MaxAge: 3600, HttpOnly: true}

	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func GetUserIdByToken(token string, db *sql.DB) int {
	query := `
		SELECT ID FROM users WHERE Session = ?
	`
	var id int
	db.QueryRow(query, token).Scan(&id)

	return id
}

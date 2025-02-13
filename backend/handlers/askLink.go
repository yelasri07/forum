package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
)

func AskLink(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		RenderError(w, http.StatusMethodNotAllowed)
		return
	}

	token, err := r.Cookie("Token")
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

	_, err = db.Exec("UPDATE users SET AuthType = ? WHERE ID = ?", 1, UserID)
	if err != nil {
		RenderError(w, http.StatusInternalServerError)
		return
	}

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

package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
)

func AskLink(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		RenderError(w, http.StatusMethodNotAllowed)
		return
	}

	token, _ := r.Cookie("UserID")
	UserID := GetUserIdByToken(token.Value, db)

	UserIDFromBrowser, _ := strconv.Atoi(r.FormValue("accept"))

	if UserID == UserIDFromBrowser {
		fmt.Println(UserID)
	} else {
		fmt.Println("err id ")
	}
}

func GetUserIdByToken(token string, db *sql.DB) int {
	query := `
		SELECT ID FROM users WHERE Session = ?
	`
	var id int
	db.QueryRow(query, token).Scan(&id)

	return id
}

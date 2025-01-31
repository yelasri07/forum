package models

import (
	"database/sql"
)

func CheckIdPost(db *sql.DB, postID int, table string) bool {
	query := `
	  SELECT ID FROM ` + table + `
	  WHERE ID = ?
	`
	var id int
	db.QueryRow(query, postID).Scan(&id)

	return id != 0
}

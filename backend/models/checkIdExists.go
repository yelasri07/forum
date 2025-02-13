package models

import (
	"database/sql"
)

// CheckIdExists checks if a given ID exists in the specified table.
func CheckIdExists(db *sql.DB, postID int, table string) (bool) {
	query := `
	  SELECT ID FROM ` + table + `
	  WHERE ID = ?
	`
	var id int
	db.QueryRow(query, postID).Scan(&id)

	return id != 0
}

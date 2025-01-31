package models

import (
	"database/sql"
)

func CheckIdPost(db *sql.DB, postID int, table string) (bool, error) {
	query := `
	  SELECT ID FROM ` + table + `
	  WHERE ID = ?
	`
	var id int
	err := db.QueryRow(query, postID).Scan(&id)
	if err != nil {
		return false , err
	}
	return id != 0 , nil
}

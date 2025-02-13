package models

import (
	"database/sql"
)

func CheckCatExists(categories []string, db *sql.DB) bool {
	for _, cat := range categories {
		query := `
			SELECT ID FROM Category
			WHERE ID = ?
		`
		var c int
		db.QueryRow(query, cat).Scan(&c)
		if c == 0 {
			return false
		}
	}
	return true
}

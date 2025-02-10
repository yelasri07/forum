package models

import "database/sql"

func GetImage(idPost int, db *sql.DB) []byte {
	query := `
		SELECT Image FROM Posts
		WHERE ID = ?
	`
	var image []byte
	db.QueryRow(query, idPost).Scan(&image)

	return image
}

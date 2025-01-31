package models

import "database/sql"

func NbOfReactInComct(CommentID string, db *sql.DB, status string) int {
	query := `
	SELECT count(ID) FROM Comment_Like
	WHERE ID_Comment = ? AND status = ?
`
	var a int
	db.QueryRow(query, CommentID, status).Scan(&a)

	return a
}

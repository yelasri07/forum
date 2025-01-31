package models

import "database/sql"

func NbOfReactInComct(CommentID string, db *sql.DB, status string) (int , error) {
	query := `
	SELECT count(ID) FROM Comment_Like
	WHERE ID_Comment = ? AND status = ?
`
	var a int
	err := db.QueryRow(query, CommentID, status).Scan(&a)
	if err != nil {
		return 0 , err
	}
	return a , nil
}

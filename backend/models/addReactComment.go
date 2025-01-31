package models

import "database/sql"

func AddReactInTheComment(db *sql.DB, status string, userID int, NewIdComment int) error {
	Query := `
	   SELECT cl.status
	   FROM Comment_like cl
	   WHERE ID_User = ? AND ID_Comment = ?; 
	`
	var statuss string
	db.QueryRow(Query, userID, NewIdComment).Scan(&statuss)

	if statuss == "" {
		db.Exec("INSERT INTO comment_like VALUES (?,?,?,?)", nil, status, userID, NewIdComment)
	} else {
		if (statuss == "like" && status == "dislike") || (statuss == "dislike" && status == "like") {
			db.Exec("UPDATE comment_like SET status = ? WHERE ID_User = ? AND ID_comment = ?", status, userID, NewIdComment)
		} else {
			db.Exec("DELETE FROM comment_like WHERE ID_User = ? AND ID_Comment = ?", userID, NewIdComment)
		}
	}
	return nil
}

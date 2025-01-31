package models

import (
	"database/sql"
)

func AddReact(db *sql.DB, status string, userID int, postID int) error {
	query := `
	   SELECT pl.status
	   FROM post_like pl
	   WHERE ID_User = ? AND ID_Post = ?; 
	`
	var statuss string
	db.QueryRow(query, userID, postID).Scan(&statuss)

	if statuss == "" {
		db.Exec("INSERT INTO post_like VALUES (?,?,?,?)", nil, status, userID, postID)
	} else {
		if (statuss == "like" && status == "dislike") || (statuss == "dislike" && status == "like") {
			db.Exec("UPDATE post_like SET status = ? WHERE ID_User = ? AND ID_Post = ?", status, userID, postID)
		} else {
			db.Exec("DELETE FROM post_like WHERE ID_User = ? AND ID_Post = ?", userID, postID)
		}
	}

	return nil
}

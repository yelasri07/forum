package models

import (
	"database/sql"
)

// AddReact manages the user's reaction (like/dislike) on a post, either inserting, updating, or deleting the reaction based on the current status.
func AddReact(db *sql.DB, status string, userID int, postID int) error {
	query := `
	   SELECT pl.status
	   FROM post_like pl
	   WHERE ID_User = ? AND ID_Post = ?; 
	`
	var statusFromDB string
	var err error
	db.QueryRow(query, userID, postID).Scan(&statusFromDB)

	if statusFromDB == "" {
		_, err = db.Exec("INSERT INTO post_like VALUES (?,?,?,?)", nil, status, userID, postID)
		if err != nil {
			return err
		}
	} else {
		if (statusFromDB == "like" && status == "dislike") || (statusFromDB == "dislike" && status == "like") {
			_, err = db.Exec("UPDATE post_like SET status = ? WHERE ID_User = ? AND ID_Post = ?", status, userID, postID)
			if err != nil {
				return err
			}
		} else {
			_, err = db.Exec("DELETE FROM post_like WHERE ID_User = ? AND ID_Post = ?", userID, postID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

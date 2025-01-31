package models

import (
	"database/sql"
	"fmt"
)

func AddReact(db *sql.DB, status string, userID int, postID int) error {
	query := `
	   SELECT pl.status
	   FROM post_like pl
	   WHERE ID_User = ? AND ID_Post = ?; 
	`
	var statuss string
	var err error
	err = db.QueryRow(query, userID, postID).Scan(&statuss)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err = db.Exec("INSERT INTO post_like VALUES (?,?,?,?)", nil, status, userID, postID)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}

	if statuss != "" {
		if (statuss == "like" && status == "dislike") || (statuss == "dislike" && status == "like") {
			_, err = db.Exec("UPDATE post_like SET status = ? WHERE ID_User = ? AND ID_Post = ?", status, userID, postID)
			if err != nil {
				fmt.Println("b7al walo")
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

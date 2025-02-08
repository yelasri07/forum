package models

import (
	"database/sql"
	"time"
)

// AddCommentsInDB inserts a new comment into the database associated with a specific user and post.
func AddCommentsInDB(db *sql.DB, userID int, postID int, comment string) error {
	query := `
	INSERT INTO comment VALUES (?,?,?,?,?)
	`
	_, err := db.Exec(query, nil, comment, time.Now(), userID, postID)
	if err != nil {
		return err
	}
	return nil
}

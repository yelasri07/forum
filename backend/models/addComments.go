package models

import (
	"database/sql"
	"time"
)

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

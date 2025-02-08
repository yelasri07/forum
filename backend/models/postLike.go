package models

import (
	"database/sql"
	"log"
)

// GetTotalLikesByUser retrieves the total number of 'like' reactions given by a specific user.
func GetTotalLikesByUser(db *sql.DB, userID int) (int, error) {
	var countLikes int
	err := db.QueryRow("SELECT COUNT(status) FROM Post_Like WHERE ID_User = ? AND status = 'like'", userID).Scan(&countLikes)
	if err != nil {
		log.Printf("Error counting likes for user %d: %v", userID, err)
		return 0, err
	}
	return countLikes, nil
}

// CountNbOfLikes retrieves the total number of 'like' reactions for a specific post.
func CountNbOfLikes(idPost int, db *sql.DB) (int, error) {
	query := `
        SELECT count(ID) FROM Post_Like
        WHERE ID_Post = ? AND status = ?
    `
	var a int
	err := db.QueryRow(query, idPost, "like").Scan(&a)
	if err != nil {
		return -1, err
	}
	return a, nil
}

// CountNbOfDislikes retrieves the total number of 'dislike' reactions for a specific post.
func CountNbOfDislikes(idPost int, db *sql.DB) (int, error) {
	query := `
        SELECT count(ID) FROM Post_Like
        WHERE ID_Post = ? AND status = ?
    `
	var count int
	err := db.QueryRow(query, idPost, "dislike").Scan(&count)
	if err != nil {
		return -1, nil
	}
	return count, nil
}

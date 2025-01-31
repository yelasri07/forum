package models

import (
	"database/sql"
	"log"
)

func GetTotalLikesByUser(db *sql.DB, userID int) (int, error) {
	var countLikes int
	err := db.QueryRow("SELECT COUNT(status) FROM Post_Like WHERE ID_User = ? AND status = 'like'", userID).Scan(&countLikes)
	if err != nil {
		log.Printf("Error counting likes for user %d: %v", userID, err)
		return 0, err
	}
	return countLikes, nil
}

func CountNbOfLikes(idPost int, db *sql.DB) int {
    query := `
        SELECT count(ID) FROM Post_Like
        WHERE ID_Post = ? AND status = ?
    `
    var a int
    db.QueryRow(query, idPost, "like").Scan(&a)

    return a
}

func CountNbOfDislikes(idPost int, db *sql.DB) int {
    query := `
        SELECT count(ID) FROM Post_Like
        WHERE ID_Post = ? AND status = ?
    `
    var count int
    db.QueryRow(query, idPost, "dislike").Scan(&count)

    return count
}

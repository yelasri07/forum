package models

import (
	"database/sql"
	"strconv"

	"forum/utils"
)

// SelectTheComment retrieves a list of comments for a specific post from the database, along with their like and dislike counts, 
// and the reaction status of the user (like/dislike). It also formats the comment's creation date and adds reaction status flags.
func SelectTheComment(postID, userID int, db *sql.DB) ([]*Comment, error) {
	query := `
	   SELECT c.ID , c.Content , u.UserName , c.DateCreation
	   FROM comment c
	   JOIN posts p ON p.ID = c.ID_Post
	   JOIN users u ON u.ID = c.ID_User
	   WHERE ID_Post = ?;
	`

	rows, err := db.Query(query, postID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var Comments []*Comment

	for rows.Next() {
		var eachComment Comment

		err := rows.Scan(&eachComment.ID, &eachComment.Content, &eachComment.CreatedBy, &eachComment.DateCreation)
		if err != nil {
			return nil , err
		}

		eachComment.Date = utils.DateFromat(eachComment.DateCreation)
		eachComment.Like , err = NbOfReactInComct(eachComment.ID, db, "like")
		if err != nil {
			return nil , err
		}

		eachComment.Dislike, err = NbOfReactInComct(eachComment.ID, db, "dislike")
		if err != nil {
			return nil , err
		}

		status := ""
		if userID != -1 {
			id, _ := strconv.Atoi(eachComment.ID)
			status, err = GetReaction(userID, "Comment_Like", "ID_Comment", id, db)
			if err != nil {
				return nil, err
			}
			
			if status == "like" {
				eachComment.CUserLiked = true
			} else if status == "dislike" {
				eachComment.CUserDisliked = true
			}
		}
		Comments = append(Comments, &eachComment)
	}
	return Comments, nil
}

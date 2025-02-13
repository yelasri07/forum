package models

import (
	"database/sql"

	"forum/utils"
)

// postService enriches a PostCat object with additional details such as the count of likes, dislikes, comments,
func postService(post *PostCat, db *sql.DB, UserId int) error {
	var err error
	post.Like, err = CountNbOfLikes(post.ID, db)
	if err != nil {
		return err
	}

	post.Dislike, err = CountNbOfDislikes(post.ID, db)
	if err != nil {
		return err
	}

	comment, err := SelectTheComment(post.ID, UserId, db)
	if err != nil {
		return err
	}

	post.Comments = append(post.Comments, comment...)
	post.NemberOfComment = len(post.Comments)

	post.Date = utils.DateFromat(post.DateCreation)
	status := ""
	if UserId != 0 {
		status, err = GetReaction(UserId, "Post_Like", "ID_Post", post.ID, db)
		if err != nil {
			return err
		}
		if status == "like" {
			post.UserLiked = true
		} else if status == "dislike" {
			post.UserDisliked = true
		}
	}

	return nil
}

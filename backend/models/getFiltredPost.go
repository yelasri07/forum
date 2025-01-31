package models

import (
	"database/sql"
	"slices"

	"forum/utils"
)

func GetPostByID(db *sql.DB, postID, UserId int) (*PostCat, error) {
	query := `
	SELECT 
	p.ID, 
	p.Title, 
	p.Content, 
	p.DateCreation, 
	GROUP_CONCAT(c.Name_Category, ' #') AS Categories, 
	u.UserName 
	FROM Posts p 
	JOIN 
	PostCategory pc ON p.ID = pc.ID_Post 
	JOIN Category c ON pc.ID_Category = c.ID 
	JOIN users u ON p.ID_User=u.ID 
	WHERE p.ID=?
	GROUP BY p.ID
	ORDER BY p.DateCreation DESC; `
	rows, err := db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var post PostCat
	for rows.Next() {

		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.DateCreation, &post.Categories, &post.CreatedBy)
		if err != nil {
			return nil, err
		}
		post.Like = CountNbOfLikes(post.ID, db)
		post.Dislike = CountNbOfDislikes(post.ID, db)

		comment, err := SelectTheComment(post.ID, UserId, db)
		if err != nil {
			return nil, err
		}

		post.Comments = append(post.Comments, comment...)
		post.NemberOfComment = len(post.Comments)

		post.Date = utils.DateFromat(post.DateCreation)
		status := ""
		if UserId != -1 {
			status, err = GetReaction(UserId, "Post_Like", "ID_Post", post.ID, db)
			if err != nil {
				return nil, err
			}
			if status == "like" {
				post.UserLiked = true
			} else if status == "dislike" {
				post.UserDisliked = true
			}
		}

	}
	return &post, nil
}

// filter >By cat
func GetPostIDsByCategory(db *sql.DB, ids *[]int, category string) error {
	query := `
		SELECT p.ID
		FROM Posts p
		INNER JOIN PostCategory pc on p.ID==pc.ID_Post
		WHERE ID_Category==?;
	`
	rows, err := db.Query(query, category)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		id := 0
		if err := rows.Scan(&id); err != nil {
			return err
		}
		if !slices.Contains(*ids, id) {
			*ids = append(*ids, id)
		}

	}
	if err := rows.Err(); err != nil {
		return err
	}
	return nil
}

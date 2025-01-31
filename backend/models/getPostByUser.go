package models

import (
	"database/sql"

	"forum/utils"
)

// filter >By  user
func GetAllPostCatByUser(db *sql.DB, userID int) ([]*PostCat, error) {
	query := `
	SELECT 
		p.ID, 
		p.Title, 
		p.Content, 
		p.DateCreation, 
		GROUP_CONCAT(c.Name_Category, ' #') AS Categories, 
		u.UserName
	FROM 
		Posts p
	JOIN 
		PostCategory pc ON p.ID = pc.ID_Post
	JOIN 
		Category c ON pc.ID_Category = c.ID
	JOIN 
		users u ON p.ID_User = u.ID
	WHERE 
		p.ID_User = ?
	GROUP BY 
		p.ID
	ORDER BY 
		p.DateCreation DESC;`

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*PostCat

	for rows.Next() {
		var post PostCat
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.DateCreation, &post.Categories, &post.CreatedBy)
		if err != nil {
			return nil, err
		}

		post.Date = utils.DateFromat(post.DateCreation)

		post.Like = CountNbOfLikes(post.ID, db)
		post.Dislike = CountNbOfDislikes(post.ID, db)
		comment, err := SelectTheComment(post.ID, userID, db)
		if err != nil {
			return nil, err
		}
		post.Comments = append(post.Comments, comment...)
		post.NemberOfComment = len(post.Comments)
		posts = append(posts, &post)
		status := ""
		if userID != -1 {
			status, err = GetReaction(userID, "Post_Like", "ID_Post", post.ID, db)
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

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

package models

import (
	"database/sql"
)

// GetPostsByLike retrieves posts that the given user has liked.
func GetPostsByLike(db *sql.DB, UserID int) ([]*PostCat, error) {
	query := `
		SELECT 
			p.ID, 
			p.Title, 
			p.Content, 
			p.DateCreation, 
			GROUP_CONCAT(c.Name_Category, ' #') AS Categories, 
			u.UserName
		FROM Posts p
		INNER JOIN Post_Like ON p.ID = Post_Like.ID_Post
		INNER JOIN PostCategory pc ON p.ID = pc.ID_Post
		INNER JOIN Category c ON pc.ID_Category = c.ID 
		INNER JOIN users u ON p.ID_User=u.ID
		WHERE Post_Like.ID_User = ? AND Post_Like.status = "like"
		GROUP BY 
		p.ID
		ORDER BY 
		p.DateCreation DESC;
	`

	rows, err := db.Query(query, UserID)
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

		err = postService(&post, db, UserID)
		if err != nil {
			return nil, err
		}

		posts = append(posts, &post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

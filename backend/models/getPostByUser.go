package models

import (
	"database/sql"
)

// GetAllPostCatByUser retrieves all posts created by a specific user, along with their categories and creator details.
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

		err = postService(&post, db, userID)
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

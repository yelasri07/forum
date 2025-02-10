package models

import (
	"database/sql"
	"slices"
)

// GetPostByID retrieves a post by its ID along with its categories and creator details.
func GetPostByID(db *sql.DB, postID, UserId int) (*PostCat, error) {
	query := `
	SELECT 
	p.ID, 
	p.Title, 
	p.Content, 
	p.DateCreation, 
	GROUP_CONCAT(c.Name_Category, ' #') AS Categories, 
	u.UserName,
	p.Image
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

		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.DateCreation, &post.Categories, &post.CreatedBy,&post.Image)
		if err != nil {
			return nil, err
		}

		err = postService(&post, db, UserId)
		if err != nil {
			return nil, err
		}

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &post, nil
}

// GetPostIDsByCategory retrieves post IDs that belong to a specific category and appends them to the provided slice if they are not already present.
func GetPostIDsByCategory(db *sql.DB, category string, ids *[]int) error {
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

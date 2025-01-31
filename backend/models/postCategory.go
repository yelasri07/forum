package models

import (
	"database/sql"

	"forum/utils"
)

func GetAllCategories(db *sql.DB) ([]*Category, error) {
	query := "SELECT ID, Name_Category FROM Category"

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*Category
	for rows.Next() {
		var category Category

		err := rows.Scan(&category.ID, &category.NameCategory)
		if err != nil {
			return nil, err
		}

		categories = append(categories, &category)
	}
	return categories, nil
}

func GetAllPostCat(db *sql.DB, UserId int) ([]*PostCat, error) {
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
	GROUP BY p.ID
	ORDER BY p.DateCreation DESC; `
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var PostCats []*PostCat
	for rows.Next() {
		var post PostCat

		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.DateCreation, &post.Categories, &post.CreatedBy)
		if err != nil {
			return nil, err
		}

		post.Like, err = CountNbOfLikes(post.ID, db)
		if err != nil {
			return nil, err
		}

		post.Dislike, err = CountNbOfDislikes(post.ID, db)
		if err != nil {
			return nil, err
		}

		comment, err := SelectTheComment(post.ID, UserId, db)
		if err != nil {
			return nil, err
		}

		post.Comments = append(post.Comments, comment...)
		post.NemberOfComment = len(post.Comments)

		// Correct date formatting
		post.Date = utils.DateFromat(post.DateCreation)
		PostCats = append(PostCats, &post)
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

	return PostCats, nil
}

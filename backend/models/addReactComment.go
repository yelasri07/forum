package models

import "database/sql"

// AddReactInTheComment handles a user's reaction (like/dislike) to a comment, 
// either inserting, updating, or deleting the reaction based on the current status.
func AddReactInTheComment(db *sql.DB, status string, userID int, NewIdComment int) error {
	Query := `
	   SELECT cl.status
	   FROM Comment_like cl
	   WHERE ID_User = ? AND ID_Comment = ?; 
	`
	var statusFormDB string
	var err error
	db.QueryRow(Query, userID, NewIdComment).Scan(&statusFormDB)

	if statusFormDB == "" {
		_, err = db.Exec("INSERT INTO comment_like VALUES (?,?,?,?)", nil, status, userID, NewIdComment)
		if err != nil {
			return err
		}
	} else {
		if (statusFormDB == "like" && status == "dislike") || (statusFormDB == "dislike" && status == "like") {
			_, err = db.Exec("UPDATE comment_like SET status = ? WHERE ID_User = ? AND ID_comment = ?", status, userID, NewIdComment)
			if err != nil {
				return err
			}
		} else {
			_, err = db.Exec("DELETE FROM comment_like WHERE ID_User = ? AND ID_Comment = ?", userID, NewIdComment)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

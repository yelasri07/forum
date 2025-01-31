package models

import (
	"database/sql"
	"fmt"
)

func GetReaction(idUser int, table string, NRow string, idPost int, db *sql.DB) (string, error) {
	query := fmt.Sprintf(`
        SELECT status 
        FROM %s 
        WHERE ID_User = ? AND %s = ?
    `, table, NRow)

	var status string
	err := db.QueryRow(query, idUser, idPost).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil 
		}
		return "", err
	}

	return status, nil
}

package models

import "database/sql"

func UserExists(db *sql.DB, value string, searchBy string) (bool, error) {
	rows, err := db.Query("SELECT ID from Users where" +searchBy+" =?", value)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if rows.Next() {
		return false, nil
	}

	return true, nil
}

func VerifyEmail(db *sql.DB, Email string) (int, error) {
	rows, err := db.Query("SELECT ID FROM users WHERE Email = ?", Email)
	if err != nil {
		return -1, err
	}
	defer rows.Close()

	if rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			return -1, err
		}
		return id, nil
	}

	return -1, nil
}

func VerifyPassword(db *sql.DB, id int) (string, error) {
	rows, err := db.Query("SELECT Password FROM users WHERE ID = ?", id)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	if rows.Next() {
		pass := ""
		err := rows.Scan(&pass)
		if err != nil {
			return "", err
		}
		return pass, nil
	}

	return "", nil
}

package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

func GenerateToken(id int, db *sql.DB) (string, error) {
	u2, err := uuid.NewV6()
	if err != nil {
		return "", err
	}

	token := u2.String()
	expirationTime := time.Now().UTC().Add(time.Hour)

	_, err = db.Exec("UPDATE users set Session=? , Expared_At=?  WHERE ID=?", token, expirationTime, id)
	if err != nil {
		return "", err
	}

	return token, nil
}

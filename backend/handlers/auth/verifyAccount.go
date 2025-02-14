package auth

import (
	"database/sql"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"forum/backend/handlers"
	"forum/backend/models"
)

func VerifyAccount(w http.ResponseWriter, r *http.Request, UserName, Email string, db *sql.DB) error {
	UserName = strings.ReplaceAll(UserName, " ", "_")

	isUniqueUserName, err := models.UserExists(db, UserName, " UserName ")
	if err != nil {
		return err
	}

	for !isUniqueUserName {
		UserName = strconv.Itoa(rand.Intn(999)) + UserName + strconv.Itoa(rand.Intn(999))
		isUniqueUserName, err = models.UserExists(db, UserName, " UserName ")
		if err != nil {
			return err
		}
	}

	ID, AuthType, err := models.VerifyEmail(db, Email)
	if err != nil {
		return err
	}

	if ID == -1 {
		result, err := db.Exec("INSERT INTO Users (UserName, Email, Password, Created_At, Session, Expared_At, AuthType) VALUES ( ?,?,?,?,?,?,?)", UserName, Email, "", time.Now(), "", nil, 1)
		if err != nil {
			return err
		}

		ID, err = result.LastInsertId()
		if err != nil {
			return err
		}
	} else {
		if AuthType == 0 {
			token, err := models.GenerateToken(int(ID), db)
			if err != nil {
				return err
			}

			cookie := &http.Cookie{Name: "UserID", Value: token, MaxAge: 3600, HttpOnly: true}

			http.SetCookie(w, cookie)

			a := models.AskLink{UserID: int(ID), AuthType: "github"}
			err = handlers.RenderTemplate(w, "askLink.html", a, http.StatusOK)
			if err != nil {
				return err
			}
			return err
		}
	}

	token, err := models.GenerateToken(int(ID), db)
	if err != nil {
		return err
	}

	cookie := &http.Cookie{Name: "Token", Value: token, MaxAge: 3600, HttpOnly: true}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)

	return nil
}

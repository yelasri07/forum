package github

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"forum/backend/handlers"
	"forum/backend/models"
)

const (
	client_id     = "Ov23li4QFMOofuba4DHv"
	client_secret = "5926c995e7b3bc0c8ec52d04ccfd24ef2205c8c1"
)

type GithubUser struct {
	UserName string `json:"login"`
	Email    string `json:"email"`
}

type GitHubEmail struct {
	Email      string `json:"email"`
	Primary    bool   `json:"primary"`
	Verified   bool   `json:"verified"`
	Visibility string `json:"visibility"`
}

func RegisterGithub(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	code := r.URL.Query().Get("code")
	if code == "" {
		handlers.RenderError(w, http.StatusBadRequest)
		return
	}

	accessToken, err := GetAccessToken(code)
	if err != nil {
		handlers.RenderError(w, http.StatusInternalServerError)
		return
	}

	user, err := GetDataUser(accessToken)

	if user.Email == "" {
		email, err := getPrimaryEmail(accessToken)
		if err != nil {
			handlers.RenderError(w, http.StatusInternalServerError)
			return
		}
		user.Email = email
	}

	user.UserName = strings.ReplaceAll(user.UserName, " ", "_")

	isUniqueUserName, err := models.UserExists(db, user.UserName, " UserName ")
	if err != nil {
		handlers.RenderError(w, http.StatusInternalServerError)
		return
	}

	for !isUniqueUserName {
		user.UserName = strconv.Itoa(rand.Intn(1000)) + user.UserName + strconv.Itoa(rand.Intn(1000))
		isUniqueUserName, err = models.UserExists(db, user.UserName, " UserName ")
		if err != nil {
			handlers.RenderError(w, http.StatusInternalServerError)
			return
		}
	}

	ID, AuthType, err := models.VerifyEmail(db, user.Email)
	if err != nil {
		handlers.RenderError(w, http.StatusInternalServerError)
		return
	}

	var id int64
	if ID == -1 {
		result, err := db.Exec("INSERT INTO Users (UserName, Email, Password, Created_At, Session, Expared_At) VALUES ( ?,?,?,?,?,?)", user.UserName, user.Email, "", time.Now(), "", nil)
		if err != nil {
			handlers.RenderError(w, http.StatusInternalServerError)
			return
		}

		id, err = result.LastInsertId()
		if err != nil {
			handlers.RenderError(w, http.StatusInternalServerError)
			return
		}
	} else {
		if AuthType == 0 {
			token, err := models.GenerateToken(int(ID), db)
			if err != nil {
				handlers.RenderError(w, http.StatusInternalServerError)
				return
			}

			cookie := &http.Cookie{Name: "Token", Value: token, MaxAge: 3600, HttpOnly: true}

			http.SetCookie(w, cookie)

			a := models.AskLink{UserID: ID, AuthType: "github"}
			err = handlers.RenderTemplate(w, "askLink.html", a, http.StatusOK)
			if err != nil {
				handlers.RenderError(w, http.StatusInternalServerError)
				return
			}

			return
		}
	}

	token, err := models.GenerateToken(int(id), db)
	if err != nil {
		handlers.RenderError(w, http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{Name: "Token", Value: token, MaxAge: 3600, HttpOnly: true}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func GetDataUser(accessToken string) (*GithubUser, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, err
	}

	var user GithubUser
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func GetAccessToken(code string) (string, error) {
	data := fmt.Sprintf("client_id=%s&client_secret=%s&code=%s", client_id, client_secret, code)
	req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bytes.NewBuffer([]byte(data)))
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", err
	}
	var tokenResp models.GetToken
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", err
	}
	if tokenResp.Token == "" {
		return "", fmt.Errorf("no access token in response: %s", string(body))
	}

	return tokenResp.Token, nil
}

func getPrimaryEmail(accessToken string) (string, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", err
	}

	var emails []GitHubEmail
	if err := json.Unmarshal(body, &emails); err != nil {
		return "", fmt.Errorf("parsing response: %w", err)
	}

	for _, email := range emails {
		if email.Primary {
			return email.Email, nil
		}
	}

	return "", fmt.Errorf("no primary email found")
}

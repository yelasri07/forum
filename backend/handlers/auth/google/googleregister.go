package google

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"forum/backend/handlers"
	"forum/backend/models"
)

const (
	ID_client     = "143686775395-8hvrk8e9bji1g6e2o2s7ogubeblhgrcp.apps.googleusercontent.com"
	secret_client = "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"
	redirect_uri  = "http://localhost:8080/callbackRegisterGoogle"
)

type GoogleUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func GoogleRegister(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	code := r.URL.Query().Get("code")
	if code == "" {
		handlers.RenderError(w, http.StatusBadRequest)
		return
	}

	accessToken, err := getAccesstoken(code)
	if err != nil {
		handlers.RenderError(w, http.StatusInternalServerError)
		return
	}
	user, err := getGoogleUser(accessToken)
	if err != nil {
		handlers.RenderError(w, http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Redirect(w, r, "/sign-up", http.StatusSeeOther)
		return
	}
	isUniqueEmail, err := models.UserExists(db, user.Email, " Email ")
	if err != nil {
		handlers.RenderError(w, http.StatusServiceUnavailable)
		return
	}

	isUniqueUserName, err := models.UserExists(db, user.Name, " UserName ")
	if err != nil {
		handlers.RenderError(w, http.StatusServiceUnavailable)
		return
	}
	if !isUniqueEmail || !isUniqueUserName {
		e := &models.ErrorRegister{AlreadyLogedWithGithub: "You have an account try to Login."}
		handlers.RenderTemplate(w, "register.html", e, http.StatusConflict)
		return
	}
	result, _ := db.Exec("INSERT INTO Users VALUES (?, ?, ?,?,?,?,?)", nil, user.Name, user.Email, "", time.Now().Format(time.DateTime), "", nil)
	ID, err := result.LastInsertId()
	if err != nil {
		handlers.RenderError(w, http.StatusInternalServerError)
		return
	}

	token, err := models.GenerateToken(int(ID), db)
	if err != nil {
		handlers.RenderError(w, http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{Name: "Token", Value: token, MaxAge: 3600, HttpOnly: true}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func getAccesstoken(code string) (string, error) {
	data := fmt.Sprintf("client_id=%s&client_secret=%s&code=%s&redirect_uri=%s&grant_type=authorization_code", ID_client, secret_client, code, redirect_uri)

	req, err := http.NewRequest("POST", "https://oauth2.googleapis.com/token", bytes.NewBuffer([]byte(data)))
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

	var tokenResp models.AccessTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", err
	}

	if tokenResp.AccessToken == "" {
		return "", err
	}
	return tokenResp.AccessToken, nil
}

func getGoogleUser(accessToken string) (*GoogleUser, error) {
	url := "https://www.googleapis.com/oauth2/v2/userinfo"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	var user GoogleUser
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

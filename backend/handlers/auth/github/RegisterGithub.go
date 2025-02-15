package github

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"forum/backend/handlers"
	"forum/backend/handlers/auth"
	"forum/backend/models"
)

const (
	client_id     = "Ov23li4QFMOofuba4DHv"
	client_secret = "5926c995e7b3bc0c8ec52d04ccfd24ef2205c8c1"
)

func RegisterGithub(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodGet {
		handlers.RenderError(w, http.StatusMethodNotAllowed)
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		handlers.RenderError(w, http.StatusBadRequest)
		return
	}

	accessToken, err := getAccessToken(code)
	if err != nil {
		cookie := &http.Cookie{Name: "UserID", Value: "", MaxAge: -1, HttpOnly: true}
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
		return
	}

	user, err := GetDataUser(accessToken)
	if err != nil {
		handlers.RenderError(w, http.StatusInternalServerError)
		return
	}

	if user.Email == "" {
		email, err := getPrimaryEmail(accessToken)
		if err != nil {
			handlers.RenderError(w, http.StatusInternalServerError)
			return
		}
		user.Email = email
	}

	err = auth.VerifyAccount(w, r, user.UserName, user.Email, "github",db)
	if err != nil {
		handlers.RenderError(w, http.StatusInternalServerError)
		return
	}
}

func GetDataUser(accessToken string) (*models.GithubUser, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch user data")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var user models.GithubUser
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func getAccessToken(code string) (string, error) {
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

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to fetch access token")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var tokenResp models.GetToken
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", err
	}

	if tokenResp.Token == "" {
		return "", errors.New("no access token in response")
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

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to fetch email user")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var emails []models.GitHubEmail
	if err := json.Unmarshal(body, &emails); err != nil {
		return "", err
	}

	for _, email := range emails {
		if email.Primary {
			return email.Email, nil
		}
	}

	return "", errors.New("no primary email found")
}

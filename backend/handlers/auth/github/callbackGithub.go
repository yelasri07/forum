package github

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"forum/backend/handlers"
	"forum/backend/models"
)

type GitHubUser struct {
	Login     string `json:"login"`
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

type GitHubEmail struct {
	Email      string `json:"email"`
	Primary    bool   `json:"primary"`
	Verified   bool   `json:"verified"`
	Visibility string `json:"visibility"`
}

func CallbackGithubRegister(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	code := r.URL.Query().Get("code")
	if code == "" {
		handlers.RenderError(w, http.StatusBadRequest)
		return
	}

	accessToken, err := getAccessToken(code)
	if err != nil {
		handlers.RenderError(w, http.StatusInternalServerError)
		return
	}
	user, err := getGitHubUser(accessToken)
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
	isUniqueEmail, err := models.UserExists(db, user.Email, " Email ")
	if err != nil {
		handlers.RenderError(w, http.StatusServiceUnavailable)
		return
	}

	isUniqueUserName, err := models.UserExists(db, user.Login, " UserName ")
	if err != nil {
		handlers.RenderError(w, http.StatusServiceUnavailable)
	}
	if !isUniqueUserName {
		user.Name = "0" + user.Name
	}
	if !isUniqueEmail {
		e := &models.ErrorRegister{AlreadyLogedWithGithub: "You have an account try to Login."}
		handlers.RenderTemplate(w, "register.html", e, http.StatusConflict)
		return
	}
	result, _ := db.Exec("INSERT INTO Users VALUES (?, ?, ?,?,?,?,?)", nil, user.Login, user.Email, "", time.Now().Format(time.DateTime), "", nil)

	ID, _ := result.LastInsertId()

	token, err := models.GenerateToken(int(ID), db)
	if err != nil {
		handlers.RenderError(w, http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{Name: "Token", Value: token, MaxAge: 3600, HttpOnly: true}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func getAccessToken(code string) (string, error) {
	data := fmt.Sprintf("client_id=%s&client_secret=%s&code=%s", os.Getenv("clientID"), os.Getenv("clientSecret"), code)
	req, err := http.NewRequest(
		"POST",
		"https://github.com/login/oauth/access_token",
		bytes.NewBuffer([]byte(data)),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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

	var tokenResp models.AccessTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", err
	}

	if tokenResp.AccessToken == "" {
		return "", fmt.Errorf("no access token in response: %s", string(body))
	}

	return tokenResp.AccessToken, nil
}

func getGitHubUser(accessToken string) (*GitHubUser, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("token %s", accessToken))
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d, body: %s", resp.StatusCode, string(body))
	}

	var user GitHubUser
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	return &user, nil
}

func getPrimaryEmail(accessToken string) (string, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("token %s", accessToken))
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status: %d, body: %s", resp.StatusCode, string(body))
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

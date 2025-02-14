package google

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
	client_id     = "978621672489-dsv02jh60f8jbkqdd3khg2p7jtse3om5.apps.googleusercontent.com"
	client_secret = "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"
	redirect_uri  = "http://localhost:8080/RegisterGoogle"
)

func RegisterGoogle(w http.ResponseWriter, r *http.Request, db *sql.DB) {
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

	user, err := getUserInfos(accessToken)
	if err != nil {
		handlers.RenderError(w, http.StatusInternalServerError)
		return
	}

	err = auth.VerifyAccount(w, r, user.UserName, user.Email, db)
	if err != nil {
		handlers.RenderError(w, http.StatusInternalServerError)
		return
	}
}

func getAccessToken(code string) (string, error) {
	data := fmt.Sprintf("client_id=%s&client_secret=%s&code=%s&redirect_uri=%s&grant_type=authorization_code",
		client_id, client_secret, code, redirect_uri)
	req, err := http.NewRequest("POST", "https://oauth2.googleapis.com/token", bytes.NewBuffer([]byte(data)))
	if err != nil {
		return "", err
	}

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

	var user models.GetToken
	if err := json.Unmarshal(body, &user); err != nil {
		return "", err
	}

	return user.Token, nil
}

func getUserInfos(accessToken string) (*models.GoogleUser, error) {
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch user info")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var user models.GoogleUser
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

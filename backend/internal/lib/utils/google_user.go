package utils

import (
	"encoding/json"
	"net/http"
)

type GoogleUser struct {
	Sub     string `json:"sub"`
	Name    string `json:"name"`
	Email   string `json:"email" `
	Picture string `json:"picture"`
}

func FetchGoogleUser(accessToken string) (*GoogleUser, error) {
	req, _ := http.NewRequest(
		"GET",
		"https://www.googleapis.com/oauth2/v3/userinfo",
		nil,
	)

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var user GoogleUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

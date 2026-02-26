package whoop

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const refreshBuffer = 60 * time.Second

// TokenData holds OAuth2 token information.
type TokenData struct {
	AccessToken  string  `json:"access_token"`
	RefreshToken string  `json:"refresh_token,omitempty"`
	ClientID     string  `json:"client_id"`
	ClientSecret string  `json:"client_secret"`
	ExpiresAt    float64 `json:"expires_at,omitempty"`
	CreatedAt    float64 `json:"created_at"`
}

// IsExpired returns true if the token has expired or will expire within the buffer.
func (t *TokenData) IsExpired() bool {
	if t.ExpiresAt == 0 {
		return true
	}
	return time.Now().Unix() >= int64(t.ExpiresAt)-int64(refreshBuffer.Seconds())
}

// tokenResponse is the raw JSON response from the token endpoint.
type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// RefreshToken uses the refresh token to obtain a new access token.
// It uses client_secret_post authentication (credentials in POST body).
func RefreshToken(tokenURL string, data *TokenData) (*TokenData, error) {
	if data.RefreshToken == "" {
		return nil, fmt.Errorf("no refresh token available")
	}

	form := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {data.RefreshToken},
		"client_id":     {data.ClientID},
		"client_secret": {data.ClientSecret},
	}

	resp, err := http.Post(tokenURL, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("refresh request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("refresh failed with status %d", resp.StatusCode)
	}

	var tr tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return nil, fmt.Errorf("decode refresh response: %w", err)
	}

	now := time.Now()
	expiresIn := tr.ExpiresIn
	if expiresIn == 0 {
		expiresIn = 3600
	}

	newData := &TokenData{
		AccessToken:  tr.AccessToken,
		RefreshToken: tr.RefreshToken,
		ClientID:     data.ClientID,
		ClientSecret: data.ClientSecret,
		ExpiresAt:    float64(now.Unix() + int64(expiresIn)),
		CreatedAt:    float64(now.Unix()),
	}

	// Keep old refresh token if new one wasn't provided
	if newData.RefreshToken == "" {
		newData.RefreshToken = data.RefreshToken
	}

	return newData, nil
}

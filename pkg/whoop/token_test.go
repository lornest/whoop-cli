package whoop

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenData_IsExpired(t *testing.T) {
	tests := []struct {
		name    string
		token   TokenData
		expired bool
	}{
		{
			name:    "zero expires_at is expired",
			token:   TokenData{ExpiresAt: 0},
			expired: true,
		},
		{
			name:    "past time is expired",
			token:   TokenData{ExpiresAt: float64(time.Now().Unix() - 100)},
			expired: true,
		},
		{
			name:    "within buffer is expired",
			token:   TokenData{ExpiresAt: float64(time.Now().Unix() + 30)},
			expired: true,
		},
		{
			name:    "far future is not expired",
			token:   TokenData{ExpiresAt: float64(time.Now().Unix() + 3600)},
			expired: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expired, tt.token.IsExpired())
		})
	}
}

func TestRefreshToken_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		require.NoError(t, r.ParseForm())

		assert.Equal(t, "refresh_token", r.FormValue("grant_type"))
		assert.Equal(t, "old-refresh", r.FormValue("refresh_token"))
		assert.Equal(t, "my-client-id", r.FormValue("client_id"))
		assert.Equal(t, "my-client-secret", r.FormValue("client_secret"))

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"access_token":  "new-access",
			"refresh_token": "new-refresh",
			"expires_in":    7200,
			"token_type":    "Bearer",
		})
	}))
	defer server.Close()

	data := &TokenData{
		AccessToken:  "old-access",
		RefreshToken: "old-refresh",
		ClientID:     "my-client-id",
		ClientSecret: "my-client-secret",
	}

	newData, err := RefreshToken(server.URL, data)
	require.NoError(t, err)

	assert.Equal(t, "new-access", newData.AccessToken)
	assert.Equal(t, "new-refresh", newData.RefreshToken)
	assert.Equal(t, "my-client-id", newData.ClientID)
	assert.Equal(t, "my-client-secret", newData.ClientSecret)
	assert.Greater(t, newData.ExpiresAt, float64(0))
}

func TestRefreshToken_NoRefreshToken(t *testing.T) {
	data := &TokenData{
		AccessToken: "old-access",
		ClientID:    "my-client-id",
	}

	_, err := RefreshToken("http://example.com", data)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no refresh token")
}

func TestRefreshToken_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	data := &TokenData{
		AccessToken:  "old-access",
		RefreshToken: "old-refresh",
		ClientID:     "my-client-id",
		ClientSecret: "my-secret",
	}

	_, err := RefreshToken(server.URL, data)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "400")
}

func TestRefreshToken_KeepsOldRefreshIfNotProvided(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"access_token": "new-access",
			"expires_in":   3600,
		})
	}))
	defer server.Close()

	data := &TokenData{
		AccessToken:  "old-access",
		RefreshToken: "keep-this-refresh",
		ClientID:     "my-client-id",
		ClientSecret: "my-secret",
	}

	newData, err := RefreshToken(server.URL, data)
	require.NoError(t, err)
	assert.Equal(t, "keep-this-refresh", newData.RefreshToken)
}

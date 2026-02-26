package whoop

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

const (
	oauthAuthURL     = "https://api.prod.whoop.com/oauth/oauth2/auth"
	oauthTokenURL    = "https://api.prod.whoop.com/oauth/oauth2/token"
	oauthRedirectURI = "http://localhost:8080/callback"
	oauthScopes      = "read:body_measurement read:cycles read:recovery read:sleep read:workout offline"
)

// CallbackResult holds the result from the OAuth callback.
type CallbackResult struct {
	Code  string
	State string
	Err   error
}

// Login performs the OAuth2 Authorization Code flow:
// 1. Starts local callback server
// 2. Opens browser to authorization URL
// 3. Waits for callback with auth code
// 4. Exchanges code for tokens
// 5. Saves tokens to storage
func Login(clientID, clientSecret string, storage Storage) (*TokenData, error) {
	return LoginWithTimeout(clientID, clientSecret, storage, 5*time.Minute)
}

// LoginWithTimeout performs the OAuth2 flow with a custom timeout.
func LoginWithTimeout(clientID, clientSecret string, storage Storage, timeout time.Duration) (*TokenData, error) {
	state, err := generateState()
	if err != nil {
		return nil, fmt.Errorf("generate state: %w", err)
	}

	resultCh := make(chan CallbackResult, 1)

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		callbackState := r.URL.Query().Get("state")
		errParam := r.URL.Query().Get("error")

		if errParam != "" {
			_, _ = fmt.Fprintf(w, "<html><body><h1>Authentication Failed</h1><p>%s</p></body></html>", errParam)
			resultCh <- CallbackResult{Err: fmt.Errorf("oauth error: %s", errParam)}
			return
		}

		_, _ = fmt.Fprint(w, "<html><body><h1>Authentication Successful!</h1><p>You can close this window.</p></body></html>")
		resultCh <- CallbackResult{Code: code, State: callbackState}
	})

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		return nil, fmt.Errorf("start callback server: %w", err)
	}

	server := &http.Server{Handler: mux}
	go func() { _ = server.Serve(listener) }()
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = server.Shutdown(ctx)
	}()

	authURL := buildAuthURL(clientID, state)
	if err := openBrowser(authURL); err != nil {
		fmt.Printf("Open this URL in your browser:\n%s\n", authURL)
	}

	select {
	case result := <-resultCh:
		if result.Err != nil {
			return nil, result.Err
		}
		if result.State != state {
			return nil, fmt.Errorf("state mismatch: expected %s, got %s", state, result.State)
		}
		return exchangeCode(clientID, clientSecret, result.Code, storage)
	case <-time.After(timeout):
		return nil, fmt.Errorf("login timed out after %s", timeout)
	}
}

func buildAuthURL(clientID, state string) string {
	params := url.Values{
		"response_type": {"code"},
		"client_id":     {clientID},
		"redirect_uri":  {oauthRedirectURI},
		"scope":         {oauthScopes},
		"state":         {state},
	}
	return oauthAuthURL + "?" + params.Encode()
}

func exchangeCode(clientID, clientSecret, code string, storage Storage) (*TokenData, error) {
	form := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"redirect_uri":  {oauthRedirectURI},
	}

	resp, err := http.Post(oauthTokenURL, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("exchange code: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("code exchange failed with status %d", resp.StatusCode)
	}

	var tr tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return nil, fmt.Errorf("decode token response: %w", err)
	}

	now := time.Now()
	expiresIn := tr.ExpiresIn
	if expiresIn == 0 {
		expiresIn = 3600
	}

	data := &TokenData{
		AccessToken:  tr.AccessToken,
		RefreshToken: tr.RefreshToken,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		ExpiresAt:    float64(now.Unix() + int64(expiresIn)),
		CreatedAt:    float64(now.Unix()),
	}

	if err := storage.Save(data); err != nil {
		return nil, fmt.Errorf("save token: %w", err)
	}

	return data, nil
}

func generateState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func openBrowser(url string) error {
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", url).Start()
	case "linux":
		return exec.Command("xdg-open", url).Start()
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

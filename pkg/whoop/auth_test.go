package whoop

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildAuthURL(t *testing.T) {
	result := buildAuthURL("my-client-id", "test-state")

	u, err := url.Parse(result)
	require.NoError(t, err)

	assert.Equal(t, "api.prod.whoop.com", u.Host)
	assert.Equal(t, "/oauth/oauth2/auth", u.Path)
	assert.Equal(t, "code", u.Query().Get("response_type"))
	assert.Equal(t, "my-client-id", u.Query().Get("client_id"))
	assert.Equal(t, "http://localhost:8080/callback", u.Query().Get("redirect_uri"))
	assert.Contains(t, u.Query().Get("scope"), "read:recovery")
	assert.Equal(t, "test-state", u.Query().Get("state"))
}

func TestGenerateState(t *testing.T) {
	state1, err := generateState()
	require.NoError(t, err)
	assert.Len(t, state1, 32) // 16 bytes hex encoded

	state2, err := generateState()
	require.NoError(t, err)
	assert.NotEqual(t, state1, state2) // Should be random
}

func TestCallbackHandler_Success(t *testing.T) {
	resultCh := make(chan CallbackResult, 1)

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state")
		resultCh <- CallbackResult{Code: code, State: state}
	})

	// Simulate a callback request
	req, _ := http.NewRequest("GET", "/callback?code=auth-code-123&state=test-state", nil)
	w := &fakeResponseWriter{}
	mux.ServeHTTP(w, req)

	result := <-resultCh
	assert.Equal(t, "auth-code-123", result.Code)
	assert.Equal(t, "test-state", result.State)
	assert.NoError(t, result.Err)
}

func TestCallbackHandler_Error(t *testing.T) {
	resultCh := make(chan CallbackResult, 1)

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		errParam := r.URL.Query().Get("error")
		if errParam != "" {
			resultCh <- CallbackResult{Err: assert.AnError}
			return
		}
	})

	req, _ := http.NewRequest("GET", "/callback?error=access_denied", nil)
	w := &fakeResponseWriter{}
	mux.ServeHTTP(w, req)

	result := <-resultCh
	assert.Error(t, result.Err)
}

// fakeResponseWriter is a minimal http.ResponseWriter for testing handlers.
type fakeResponseWriter struct {
	code int
	body []byte
}

func (f *fakeResponseWriter) Header() http.Header { return http.Header{} }
func (f *fakeResponseWriter) Write(b []byte) (int, error) {
	f.body = append(f.body, b...)
	return len(b), nil
}
func (f *fakeResponseWriter) WriteHeader(code int) { f.code = code }

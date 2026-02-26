package whoop

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type memoryStorage struct {
	data *TokenData
}

func (m *memoryStorage) Save(data *TokenData) error { m.data = data; return nil }
func (m *memoryStorage) Load() (*TokenData, error) {
	if m.data == nil {
		return nil, os.ErrNotExist
	}
	return m.data, nil
}
func (m *memoryStorage) Delete() error { m.data = nil; return nil }

func validToken() *TokenData {
	return &TokenData{
		AccessToken:  "valid-token",
		RefreshToken: "refresh-token",
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		ExpiresAt:    float64(9999999999), // far future
		CreatedAt:    float64(1700000000),
	}
}

func TestClient_GetBodyMeasurement(t *testing.T) {
	fixture, err := os.ReadFile("../../testdata/body.json")
	require.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, PathBodyMeasurement, r.URL.Path)
		assert.Equal(t, "Bearer valid-token", r.Header.Get("Authorization"))
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	storage := &memoryStorage{data: validToken()}
	client := NewClientWithBase(server.URL, server.URL, storage)

	body, err := client.GetBodyMeasurement()
	require.NoError(t, err)
	assert.Equal(t, int64(12345), body.UserID)
	assert.InDelta(t, 1.8288, body.HeightMeter, 0.001)
}

func TestClient_GetCycles(t *testing.T) {
	fixture, err := os.ReadFile("../../testdata/cycles.json")
	require.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, PathCycles, r.URL.Path)
		assert.Equal(t, "2026-02-20T00:00:00.000Z", r.URL.Query().Get("start"))
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	storage := &memoryStorage{data: validToken()}
	client := NewClientWithBase(server.URL, server.URL, storage)

	resp, err := client.GetCycles(&QueryParams{Start: "2026-02-20T00:00:00.000Z"})
	require.NoError(t, err)
	assert.Len(t, resp.Records, 3)
	assert.Equal(t, "SCORED", resp.Records[0].ScoreState)
}

func TestClient_GetRecovery(t *testing.T) {
	fixture, err := os.ReadFile("../../testdata/recovery.json")
	require.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, PathRecovery, r.URL.Path)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	storage := &memoryStorage{data: validToken()}
	client := NewClientWithBase(server.URL, server.URL, storage)

	resp, err := client.GetRecovery(nil)
	require.NoError(t, err)
	assert.Len(t, resp.Records, 3)
	assert.InDelta(t, 78.0, resp.Records[0].Score.RecoveryScore, 0.01)
}

func TestClient_GetSleep(t *testing.T) {
	fixture, err := os.ReadFile("../../testdata/sleep.json")
	require.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, PathSleep, r.URL.Path)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	storage := &memoryStorage{data: validToken()}
	client := NewClientWithBase(server.URL, server.URL, storage)

	resp, err := client.GetSleep(nil)
	require.NoError(t, err)
	assert.Len(t, resp.Records, 2)
}

func TestClient_GetWorkouts(t *testing.T) {
	fixture, err := os.ReadFile("../../testdata/workouts.json")
	require.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, PathWorkouts, r.URL.Path)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	storage := &memoryStorage{data: validToken()}
	client := NewClientWithBase(server.URL, server.URL, storage)

	resp, err := client.GetWorkouts(nil)
	require.NoError(t, err)
	assert.Len(t, resp.Records, 2)
	require.NotNil(t, resp.NextToken)
	assert.Equal(t, "abc123", *resp.NextToken)
}

func TestClient_401_RefreshAndRetry(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == PathBodyMeasurement {
			calls++
			if calls == 1 {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			// Second call should have new token
			assert.Equal(t, "Bearer new-access-token", r.Header.Get("Authorization"))
			_ = json.NewEncoder(w).Encode(BodyMeasurement{UserID: 1, MaxHeartRate: 180})
			return
		}
		// Token refresh endpoint
		_ = json.NewEncoder(w).Encode(map[string]any{
			"access_token":  "new-access-token",
			"refresh_token": "new-refresh-token",
			"expires_in":    3600,
		})
	}))
	defer server.Close()

	storage := &memoryStorage{data: validToken()}
	client := NewClientWithBase(server.URL, server.URL, storage)

	body, err := client.GetBodyMeasurement()
	require.NoError(t, err)
	assert.Equal(t, int64(1), body.UserID)
	assert.Equal(t, 2, calls)

	// Verify token was saved
	assert.Equal(t, "new-access-token", storage.data.AccessToken)
}

func TestClient_429_RateLimited(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	storage := &memoryStorage{data: validToken()}
	client := NewClientWithBase(server.URL, server.URL, storage)

	_, err := client.GetBodyMeasurement()
	assert.ErrorIs(t, err, ErrRateLimited)
}

func TestClient_500_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("internal server error"))
	}))
	defer server.Close()

	storage := &memoryStorage{data: validToken()}
	client := NewClientWithBase(server.URL, server.URL, storage)

	_, err := client.GetBodyMeasurement()
	require.Error(t, err)
	var apiErr *APIError
	assert.True(t, errors.As(err, &apiErr))
	assert.Equal(t, 500, apiErr.StatusCode)
}

func TestClient_NoToken(t *testing.T) {
	storage := &memoryStorage{}
	client := NewClientWithBase("http://localhost", "http://localhost", storage)

	_, err := client.GetBodyMeasurement()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "load token")
}

func TestClient_ProactiveRefresh(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == PathBodyMeasurement {
			assert.Equal(t, "Bearer refreshed-token", r.Header.Get("Authorization"))
			_ = json.NewEncoder(w).Encode(BodyMeasurement{UserID: 42, MaxHeartRate: 190})
			return
		}
		// Token endpoint
		_ = json.NewEncoder(w).Encode(map[string]any{
			"access_token":  "refreshed-token",
			"refresh_token": "new-refresh",
			"expires_in":    3600,
		})
	}))
	defer server.Close()

	// Token that's about to expire (within buffer)
	expiredToken := validToken()
	expiredToken.ExpiresAt = float64(0)

	storage := &memoryStorage{data: expiredToken}
	client := NewClientWithBase(server.URL, server.URL, storage)

	body, err := client.GetBodyMeasurement()
	require.NoError(t, err)
	assert.Equal(t, int64(42), body.UserID)
	assert.Equal(t, "refreshed-token", storage.data.AccessToken)
}

func TestBuildParams(t *testing.T) {
	params := buildParams(&QueryParams{
		Start:     "2026-02-20T00:00:00.000Z",
		End:       "2026-02-26T23:59:59.999Z",
		Limit:     25,
		NextToken: "abc123",
	})
	assert.Equal(t, "2026-02-20T00:00:00.000Z", params.Get("start"))
	assert.Equal(t, "2026-02-26T23:59:59.999Z", params.Get("end"))
	assert.Equal(t, "25", params.Get("limit"))
	assert.Equal(t, "abc123", params.Get("nextToken"))
}

func TestBuildParams_Nil(t *testing.T) {
	params := buildParams(nil)
	assert.Empty(t, params.Encode())
}

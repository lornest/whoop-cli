package whoop

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrRateLimited  = errors.New("rate limited by Whoop API")
)

// APIError represents an HTTP error from the Whoop API.
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error %d: %s", e.StatusCode, e.Message)
}

// Client is an authenticated HTTP client for the Whoop API.
type Client struct {
	baseURL    string
	tokenURL   string
	httpClient *http.Client
	storage    Storage
}

// NewClient creates a new API client.
func NewClient(storage Storage) *Client {
	return &Client{
		baseURL:    BaseURL,
		tokenURL:   TokenURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		storage:    storage,
	}
}

// NewClientWithBase creates a client with a custom base URL (for testing).
func NewClientWithBase(baseURL, tokenURL string, storage Storage) *Client {
	return &Client{
		baseURL:    baseURL,
		tokenURL:   tokenURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		storage:    storage,
	}
}

func (c *Client) doRequest(path string, params url.Values) ([]byte, error) {
	token, err := c.getValidToken()
	if err != nil {
		return nil, fmt.Errorf("get token: %w", err)
	}

	u := c.baseURL + path
	if len(params) > 0 {
		u += "?" + params.Encode()
	}

	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// Handle 401: try refresh and retry once
	if resp.StatusCode == http.StatusUnauthorized {
		newToken, refreshErr := RefreshToken(c.tokenURL, token)
		if refreshErr != nil {
			return nil, ErrUnauthorized
		}
		if err := c.storage.Save(newToken); err != nil {
			return nil, fmt.Errorf("save refreshed token: %w", err)
		}

		req2, _ := http.NewRequest(http.MethodGet, c.baseURL+path, nil)
		if len(params) > 0 {
			req2.URL.RawQuery = params.Encode()
		}
		req2.Header.Set("Authorization", "Bearer "+newToken.AccessToken)

		resp2, err := c.httpClient.Do(req2)
		if err != nil {
			return nil, fmt.Errorf("retry request failed: %w", err)
		}
		defer func() { _ = resp2.Body.Close() }()

		if resp2.StatusCode == http.StatusUnauthorized {
			return nil, ErrUnauthorized
		}
		return c.readResponse(resp2)
	}

	return c.readResponse(resp)
}

func (c *Client) readResponse(resp *http.Response) ([]byte, error) {
	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, ErrRateLimited
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, &APIError{StatusCode: resp.StatusCode, Message: string(body)}
	}

	return body, nil
}

func (c *Client) getValidToken() (*TokenData, error) {
	token, err := c.storage.Load()
	if err != nil {
		return nil, fmt.Errorf("load token: %w", err)
	}

	// Proactive refresh if expired or about to expire
	if token.IsExpired() {
		newToken, err := RefreshToken(c.tokenURL, token)
		if err != nil {
			return nil, fmt.Errorf("proactive refresh: %w", err)
		}
		if err := c.storage.Save(newToken); err != nil {
			return nil, fmt.Errorf("save refreshed token: %w", err)
		}
		return newToken, nil
	}

	return token, nil
}

func buildParams(opts *QueryParams) url.Values {
	params := url.Values{}
	if opts == nil {
		return params
	}
	if opts.Start != "" {
		params.Set("start", opts.Start)
	}
	if opts.End != "" {
		params.Set("end", opts.End)
	}
	if opts.Limit > 0 {
		params.Set("limit", strconv.Itoa(opts.Limit))
	}
	if opts.NextToken != "" {
		params.Set("nextToken", opts.NextToken)
	}
	return params
}

// GetBodyMeasurement returns the user's body measurements.
func (c *Client) GetBodyMeasurement() (*BodyMeasurement, error) {
	data, err := c.doRequest(PathBodyMeasurement, nil)
	if err != nil {
		return nil, err
	}
	var result BodyMeasurement
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("decode body measurement: %w", err)
	}
	return &result, nil
}

// GetCycles returns physiological cycles.
func (c *Client) GetCycles(opts *QueryParams) (*CycleResponse, error) {
	data, err := c.doRequest(PathCycles, buildParams(opts))
	if err != nil {
		return nil, err
	}
	var result CycleResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("decode cycles: %w", err)
	}
	return &result, nil
}

// GetRecovery returns recovery data.
func (c *Client) GetRecovery(opts *QueryParams) (*RecoveryResponse, error) {
	data, err := c.doRequest(PathRecovery, buildParams(opts))
	if err != nil {
		return nil, err
	}
	var result RecoveryResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("decode recovery: %w", err)
	}
	return &result, nil
}

// GetSleep returns sleep data.
func (c *Client) GetSleep(opts *QueryParams) (*SleepResponse, error) {
	data, err := c.doRequest(PathSleep, buildParams(opts))
	if err != nil {
		return nil, err
	}
	var result SleepResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("decode sleep: %w", err)
	}
	return &result, nil
}

// GetWorkouts returns workout data.
func (c *Client) GetWorkouts(opts *QueryParams) (*WorkoutResponse, error) {
	data, err := c.doRequest(PathWorkouts, buildParams(opts))
	if err != nil {
		return nil, err
	}
	var result WorkoutResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("decode workouts: %w", err)
	}
	return &result, nil
}

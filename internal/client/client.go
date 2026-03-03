package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/danilodrobac/asana-cli/internal/models"
)

const (
	defaultBaseURL = "https://app.asana.com/api/1.0"
	maxRetries     = 3
)

// TokenProvider returns a valid access token.
type TokenProvider func() (string, error)

// Client is an HTTP client for the Asana API.
type Client struct {
	httpClient    *http.Client
	tokenProvider TokenProvider
	baseURL       string
}

// New creates a new Asana API client.
func New(tp TokenProvider) *Client {
	return &Client{
		httpClient:    &http.Client{Timeout: 30 * time.Second},
		tokenProvider: tp,
		baseURL:       defaultBaseURL,
	}
}

// NewWithHTTPClient creates a client with a custom HTTP client and base URL (for testing).
func NewWithHTTPClient(tp TokenProvider, httpClient *http.Client, baseURL string) *Client {
	return &Client{
		httpClient:    httpClient,
		tokenProvider: tp,
		baseURL:       baseURL,
	}
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	token, err := c.tokenProvider()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			backoff := time.Duration(1<<uint(attempt-1))*time.Second + time.Duration(rand.Intn(500))*time.Millisecond
			time.Sleep(backoff)

			// Re-read body for retry if needed
			if req.GetBody != nil {
				body, err := req.GetBody()
				if err != nil {
					return nil, fmt.Errorf("failed to re-read body for retry: %w", err)
				}
				req.Body = body
			}
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("HTTP request failed: %w", err)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = fmt.Errorf("failed to read response body: %w", err)
			continue
		}

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return body, nil
		}

		// Parse error response
		var asanaErr models.AsanaErrorResponse
		msg := string(body)
		if json.Unmarshal(body, &asanaErr) == nil && len(asanaErr.Errors) > 0 {
			msg = asanaErr.Errors[0].Message
		}

		apiErr := &APIError{
			StatusCode: resp.StatusCode,
			Message:    msg,
			Code:       ErrorCode(resp.StatusCode),
		}

		// Only retry on 429 or 5xx
		if resp.StatusCode == 429 || resp.StatusCode >= 500 {
			lastErr = apiErr
			continue
		}

		return nil, apiErr
	}

	return nil, lastErr
}

// Get performs a GET request.
func (c *Client) Get(path string) ([]byte, error) {
	req, err := http.NewRequest("GET", c.baseURL+path, nil)
	if err != nil {
		return nil, err
	}
	return c.doRequest(req)
}

// Post performs a POST request with a JSON body.
func (c *Client) Post(path string, payload interface{}) ([]byte, error) {
	data, err := json.Marshal(map[string]interface{}{"data": payload})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", c.baseURL+path, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader(data)), nil
	}
	return c.doRequest(req)
}

// Put performs a PUT request with a JSON body.
func (c *Client) Put(path string, payload interface{}) ([]byte, error) {
	data, err := json.Marshal(map[string]interface{}{"data": payload})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PUT", c.baseURL+path, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader(data)), nil
	}
	return c.doRequest(req)
}

// Delete performs a DELETE request.
func (c *Client) Delete(path string) ([]byte, error) {
	req, err := http.NewRequest("DELETE", c.baseURL+path, nil)
	if err != nil {
		return nil, err
	}
	return c.doRequest(req)
}

// PostMultipart uploads a file via multipart/form-data.
func (c *Client) PostMultipart(path string, filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err := io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("failed to copy file data: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	bodyBytes := body.Bytes()
	req, err := http.NewRequest("POST", c.baseURL+path, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader(bodyBytes)), nil
	}
	return c.doRequest(req)
}

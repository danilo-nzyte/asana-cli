package client

import (
	"fmt"
)

// Exit codes for different error categories.
const (
	ExitAuthError   = 1
	ExitNotFound    = 2
	ExitValidation  = 3
	ExitRateLimited = 4
	ExitServerError = 5
	ExitUsageError  = 10
)

// APIError represents an error response from the Asana API.
type APIError struct {
	StatusCode int
	Message    string
	Code       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("asana API error (%d): %s", e.StatusCode, e.Message)
}

// ExitCode returns the appropriate CLI exit code for this error.
func (e *APIError) ExitCode() int {
	switch {
	case e.StatusCode == 401 || e.StatusCode == 403:
		return ExitAuthError
	case e.StatusCode == 404:
		return ExitNotFound
	case e.StatusCode == 400 || e.StatusCode == 422:
		return ExitValidation
	case e.StatusCode == 429:
		return ExitRateLimited
	case e.StatusCode >= 500:
		return ExitServerError
	default:
		return 1
	}
}

// ErrorCode returns a short error code string from the HTTP status.
func ErrorCode(statusCode int) string {
	switch {
	case statusCode == 401:
		return "unauthorized"
	case statusCode == 403:
		return "forbidden"
	case statusCode == 404:
		return "not_found"
	case statusCode == 400:
		return "bad_request"
	case statusCode == 422:
		return "unprocessable"
	case statusCode == 429:
		return "rate_limited"
	case statusCode >= 500:
		return "server_error"
	default:
		return "unknown_error"
	}
}

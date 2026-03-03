package output

import (
	"encoding/json"
	"fmt"
	"os"
)

// Response is the standard CLI output envelope.
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

// ErrorInfo holds structured error details.
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Success prints a success response and exits 0.
func Success(data interface{}, message string) {
	printJSON(Response{
		Success: true,
		Data:    data,
		Message: message,
	})
}

// Fail prints an error response and exits with the given code.
func Fail(code string, message string, exitCode int) {
	printJSON(Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
		},
	})
	os.Exit(exitCode)
}

func printJSON(v interface{}) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		fmt.Fprintf(os.Stderr, "failed to encode JSON: %v\n", err)
		os.Exit(1)
	}
}

package cmd

import (
	"github.com/danilodrobac/asana-cli/internal/client"
	"github.com/danilodrobac/asana-cli/internal/output"
)

// handleAPIError checks if err is an APIError and exits with the appropriate
// code and message. For non-API errors it exits with code 1.
func handleAPIError(err error) {
	if apiErr, ok := err.(*client.APIError); ok {
		output.Fail(apiErr.Code, apiErr.Message, apiErr.ExitCode())
	}
	output.Fail("unknown", err.Error(), 1)
}

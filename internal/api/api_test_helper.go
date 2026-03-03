package api

import (
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/danilodrobac/asana-cli/internal/client"
)

// newTestServer creates an httptest.Server and a client.Client pointing to it.
func newTestServer(handler http.HandlerFunc) (*client.Client, *httptest.Server) {
	server := httptest.NewServer(handler)
	c := client.NewWithHTTPClient(
		func() (string, error) { return "test-token", nil },
		&http.Client{Timeout: 5 * time.Second},
		server.URL,
	)
	return c, server
}

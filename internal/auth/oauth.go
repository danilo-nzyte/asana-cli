package auth

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"time"
)

const (
	authorizeEndpoint = "https://app.asana.com/-/oauth_authorize"
)

// RunOAuthFlow starts a local server, opens the browser for authorization,
// and exchanges the received code for tokens.
func RunOAuthFlow(cfg *Config) (*TokenData, error) {
	const port = 8931
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return nil, fmt.Errorf("failed to start local server on port %d (is it already in use?): %w", port, err)
	}
	redirectURI := fmt.Sprintf("http://localhost:%d/callback", port)

	codeCh := make(chan string, 1)
	errCh := make(chan error, 1)

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			errMsg := r.URL.Query().Get("error")
			if errMsg == "" {
				errMsg = "no authorization code received"
			}
			fmt.Fprintf(w, "<html><body><h2>Authorization failed</h2><p>%s</p><p>You can close this window.</p></body></html>", errMsg)
			errCh <- fmt.Errorf("authorization failed: %s", errMsg)
			return
		}
		fmt.Fprint(w, "<html><body><h2>Authorization successful!</h2><p>You can close this window and return to the terminal.</p></body></html>")
		codeCh <- code
	})

	server := &http.Server{Handler: mux}

	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("local server error: %w", err)
		}
	}()

	// Build authorization URL
	authURL := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=code",
		authorizeEndpoint, cfg.ClientID, redirectURI)

	fmt.Printf("Opening browser for Asana authorization...\n")
	fmt.Printf("If the browser doesn't open, visit:\n%s\n\n", authURL)

	openBrowser(authURL)

	// Wait for the callback
	var code string
	select {
	case code = <-codeCh:
	case err := <-errCh:
		server.Shutdown(context.Background())
		return nil, err
	case <-time.After(5 * time.Minute):
		server.Shutdown(context.Background())
		return nil, fmt.Errorf("authorization timed out (5 minutes)")
	}

	server.Shutdown(context.Background())

	// Exchange code for tokens
	token, err := ExchangeCode(cfg, code, redirectURI)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	}
	if cmd != nil {
		cmd.Start()
	}
}

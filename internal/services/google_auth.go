package services

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/oauth2"
	googleoauth "golang.org/x/oauth2/google"
	calendar "google.golang.org/api/calendar/v3"
)

const (
	GoogleAuthErrorMissingCredentials = "missing_credentials"
	GoogleAuthErrorCredentialsRead    = "credentials_read"
	GoogleAuthErrorInvalidCredentials = "invalid_credentials"
	GoogleAuthErrorAuthExchange       = "auth_exchange"
	GoogleAuthErrorTokenRead          = "token_read"
	GoogleAuthErrorInvalidToken       = "invalid_token"
	GoogleAuthErrorTokenPermissions   = "token_permissions"
	GoogleAuthErrorTokenSave          = "token_save"
)

var (
	errAuthCodeEmpty = errors.New("authorization code cannot be empty")

	parseGoogleCredentials = googleoauth.ConfigFromJSON
	authCodePrompter       = promptAuthCodeFromTerminal
	authCodeExchanger      = func(ctx context.Context, config *oauth2.Config, code string) (*oauth2.Token, error) {
		return config.Exchange(ctx, code)
	}
)

// GoogleAuthOptions configures terminal-friendly Google OAuth authentication.
type GoogleAuthOptions struct {
	CredentialsPath string
	TokenPath       string
	Scopes          []string
	BaseDir         string
	Input           io.Reader
	Output          io.Writer
}

// GoogleAuthResult contains the authenticated client and cache metadata.
type GoogleAuthResult struct {
	Client          *http.Client
	TokenPath       string
	UsedCachedToken bool
}

// GoogleAuthError represents typed auth lifecycle failures.
type GoogleAuthError struct {
	Kind string
	Path string
	Err  error
}

func (e *GoogleAuthError) Error() string {
	base := "google auth failed"
	if e.Kind != "" {
		base = base + ": " + e.Kind
	}
	if e.Path != "" {
		base = base + " (" + e.Path + ")"
	}
	if e.Err != nil {
		base = base + ": " + e.Err.Error()
	}
	return base
}

func (e *GoogleAuthError) Unwrap() error {
	return e.Err
}

// EnsureGoogleAuthClient returns an authenticated client and durable token cache state.
func EnsureGoogleAuthClient(ctx context.Context, opts GoogleAuthOptions) (GoogleAuthResult, error) {
	credentialsPath := strings.TrimSpace(opts.CredentialsPath)
	if credentialsPath == "" {
		credentialsPath = strings.TrimSpace(os.Getenv("GOOGLE_OAUTH_CREDENTIALS_PATH"))
	}
	if credentialsPath == "" {
		return GoogleAuthResult{}, &GoogleAuthError{Kind: GoogleAuthErrorMissingCredentials, Err: errors.New("set GOOGLE_OAUTH_CREDENTIALS_PATH or pass CredentialsPath")}
	}

	credentialsJSON, err := os.ReadFile(credentialsPath)
	if err != nil {
		kind := GoogleAuthErrorCredentialsRead
		if errors.Is(err, os.ErrNotExist) {
			kind = GoogleAuthErrorMissingCredentials
		}
		return GoogleAuthResult{}, &GoogleAuthError{Kind: kind, Path: credentialsPath, Err: err}
	}

	scopes := opts.Scopes
	if len(scopes) == 0 {
		scopes = []string{calendar.CalendarScope}
	}

	config, err := parseGoogleCredentials(credentialsJSON, scopes...)
	if err != nil {
		return GoogleAuthResult{}, &GoogleAuthError{Kind: GoogleAuthErrorInvalidCredentials, Path: credentialsPath, Err: err}
	}

	tokenPath, err := resolveGoogleTokenPath(opts)
	if err != nil {
		return GoogleAuthResult{}, err
	}

	token, usedCachedToken, err := loadOrAuthorizeToken(ctx, config, tokenPath, opts)
	if err != nil {
		return GoogleAuthResult{}, err
	}

	tokenSource := config.TokenSource(ctx, token)
	resolvedToken, err := tokenSource.Token()
	if err != nil {
		return GoogleAuthResult{}, &GoogleAuthError{Kind: GoogleAuthErrorAuthExchange, Path: tokenPath, Err: err}
	}

	if !tokensEquivalent(token, resolvedToken) {
		if err := saveOAuthToken(tokenPath, resolvedToken); err != nil {
			return GoogleAuthResult{}, err
		}
	}

	client := oauth2.NewClient(ctx, oauth2.ReuseTokenSource(resolvedToken, tokenSource))

	return GoogleAuthResult{
		Client:          client,
		TokenPath:       tokenPath,
		UsedCachedToken: usedCachedToken,
	}, nil
}

func loadOrAuthorizeToken(ctx context.Context, config *oauth2.Config, tokenPath string, opts GoogleAuthOptions) (*oauth2.Token, bool, error) {
	token, err := tokenFromFile(tokenPath)
	if err == nil {
		return token, true, nil
	}

	var authErr *GoogleAuthError
	if !errors.As(err, &authErr) || authErr.Kind != GoogleAuthErrorTokenRead {
		return nil, false, err
	}

	authCode, promptErr := authCodePrompter(opts.Input, opts.Output, config)
	if promptErr != nil {
		return nil, false, &GoogleAuthError{Kind: GoogleAuthErrorAuthExchange, Path: tokenPath, Err: promptErr}
	}

	token, exchangeErr := authCodeExchanger(ctx, config, authCode)
	if exchangeErr != nil {
		return nil, false, &GoogleAuthError{Kind: GoogleAuthErrorAuthExchange, Path: tokenPath, Err: exchangeErr}
	}

	if err := saveOAuthToken(tokenPath, token); err != nil {
		return nil, false, err
	}

	return token, false, nil
}

func resolveGoogleTokenPath(opts GoogleAuthOptions) (string, error) {
	if strings.TrimSpace(opts.TokenPath) != "" {
		return filepath.Clean(opts.TokenPath), nil
	}

	base := strings.TrimSpace(opts.BaseDir)
	if base == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return "", &GoogleAuthError{Kind: GoogleAuthErrorTokenSave, Err: fmt.Errorf("get working directory: %w", err)}
		}
		base = cwd
	}

	repoRoot, err := findRepoRoot(base)
	if err != nil {
		return "", &GoogleAuthError{Kind: GoogleAuthErrorTokenSave, Err: fmt.Errorf("find repository root: %w", err)}
	}

	return filepath.Join(repoRoot, "data", "google", "token.json"), nil
}

func tokenFromFile(path string) (*oauth2.Token, error) {
	info, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, &GoogleAuthError{Kind: GoogleAuthErrorTokenRead, Path: path, Err: err}
		}
		return nil, &GoogleAuthError{Kind: GoogleAuthErrorTokenRead, Path: path, Err: fmt.Errorf("stat token file: %w", err)}
	}

	if info.Mode().Perm() != 0o600 {
		return nil, &GoogleAuthError{Kind: GoogleAuthErrorTokenPermissions, Path: path, Err: fmt.Errorf("token file must use 0600 permissions, got %04o", info.Mode().Perm())}
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return nil, &GoogleAuthError{Kind: GoogleAuthErrorTokenRead, Path: path, Err: fmt.Errorf("read token file: %w", err)}
	}

	token := &oauth2.Token{}
	if err := json.Unmarshal(b, token); err != nil {
		return nil, &GoogleAuthError{Kind: GoogleAuthErrorInvalidToken, Path: path, Err: err}
	}

	return token, nil
}

func saveOAuthToken(path string, token *oauth2.Token) error {
	if token == nil {
		return &GoogleAuthError{Kind: GoogleAuthErrorTokenSave, Path: path, Err: errors.New("token is nil")}
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return &GoogleAuthError{Kind: GoogleAuthErrorTokenSave, Path: path, Err: fmt.Errorf("create token directory: %w", err)}
	}

	b, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		return &GoogleAuthError{Kind: GoogleAuthErrorTokenSave, Path: path, Err: fmt.Errorf("marshal token: %w", err)}
	}

	if err := os.WriteFile(path, b, 0o600); err != nil {
		return &GoogleAuthError{Kind: GoogleAuthErrorTokenSave, Path: path, Err: fmt.Errorf("write token file: %w", err)}
	}

	return nil
}

func promptAuthCodeFromTerminal(input io.Reader, output io.Writer, config *oauth2.Config) (string, error) {
	if config == nil {
		return "", errors.New("oauth config is nil")
	}

	if input == nil {
		input = os.Stdin
	}
	if output == nil {
		output = os.Stdout
	}

	authURL := config.AuthCodeURL("openppl-google-auth", oauth2.AccessTypeOffline)
	_, _ = fmt.Fprintf(output, "Open this URL in your browser and authorize access:\n%s\n", authURL)
	_, _ = fmt.Fprint(output, "Enter authorization code: ")

	reader := bufio.NewReader(input)
	code, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}

	code = strings.TrimSpace(code)
	if code == "" {
		return "", errAuthCodeEmpty
	}

	return code, nil
}

func tokensEquivalent(a, b *oauth2.Token) bool {
	if a == nil || b == nil {
		return false
	}

	return a.AccessToken == b.AccessToken &&
		a.TokenType == b.TokenType &&
		a.RefreshToken == b.RefreshToken &&
		a.Expiry.Equal(b.Expiry)
}

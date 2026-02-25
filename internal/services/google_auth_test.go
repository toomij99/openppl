package services

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"golang.org/x/oauth2"
)

func TestGoogleAuth_FirstRunLoadsCredentialsAndSavesToken(t *testing.T) {
	tmp := t.TempDir()
	credentialsPath := filepath.Join(tmp, "credentials.json")
	tokenPath := filepath.Join(tmp, "cache", "token.json")

	if err := os.WriteFile(credentialsPath, []byte(`{"installed":{}}`), 0o644); err != nil {
		t.Fatalf("write credentials file: %v", err)
	}

	originalParse := parseGoogleCredentials
	originalPrompt := authCodePrompter
	originalExchange := authCodeExchanger
	t.Cleanup(func() {
		parseGoogleCredentials = originalParse
		authCodePrompter = originalPrompt
		authCodeExchanger = originalExchange
	})

	parseGoogleCredentials = func(_ []byte, _ ...string) (*oauth2.Config, error) {
		return &oauth2.Config{Endpoint: oauth2.Endpoint{AuthURL: "https://example.test/auth", TokenURL: "https://example.test/token"}}, nil
	}
	authCodePrompter = func(_ io.Reader, _ io.Writer, _ *oauth2.Config) (string, error) {
		return "code-from-user", nil
	}
	authCodeExchanger = func(_ context.Context, _ *oauth2.Config, code string) (*oauth2.Token, error) {
		if code != "code-from-user" {
			t.Fatalf("unexpected auth code: %q", code)
		}
		return &oauth2.Token{AccessToken: "first-token", RefreshToken: "refresh", TokenType: "Bearer", Expiry: time.Now().Add(30 * time.Minute)}, nil
	}

	result, err := EnsureGoogleAuthClient(context.Background(), GoogleAuthOptions{
		CredentialsPath: credentialsPath,
		TokenPath:       tokenPath,
	})
	if err != nil {
		t.Fatalf("EnsureGoogleAuthClient failed: %v", err)
	}

	if result.Client == nil {
		t.Fatal("expected non-nil client")
	}
	if result.UsedCachedToken {
		t.Fatal("expected first-run auth, got cached token")
	}

	info, err := os.Stat(tokenPath)
	if err != nil {
		t.Fatalf("expected token file: %v", err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("expected token permissions 0600, got %04o", info.Mode().Perm())
	}
}

func TestGoogleAuth_UsesCachedTokenWithoutInteractiveFlow(t *testing.T) {
	tmp := t.TempDir()
	credentialsPath := filepath.Join(tmp, "credentials.json")
	tokenPath := filepath.Join(tmp, "token.json")

	if err := os.WriteFile(credentialsPath, []byte(`{"installed":{}}`), 0o644); err != nil {
		t.Fatalf("write credentials file: %v", err)
	}
	if err := saveOAuthToken(tokenPath, &oauth2.Token{AccessToken: "cached-token", RefreshToken: "refresh", TokenType: "Bearer", Expiry: time.Now().Add(time.Hour)}); err != nil {
		t.Fatalf("saveOAuthToken failed: %v", err)
	}

	originalParse := parseGoogleCredentials
	originalPrompt := authCodePrompter
	originalExchange := authCodeExchanger
	t.Cleanup(func() {
		parseGoogleCredentials = originalParse
		authCodePrompter = originalPrompt
		authCodeExchanger = originalExchange
	})

	parseGoogleCredentials = func(_ []byte, _ ...string) (*oauth2.Config, error) {
		return &oauth2.Config{Endpoint: oauth2.Endpoint{AuthURL: "https://example.test/auth", TokenURL: "https://example.test/token"}}, nil
	}
	authCodePrompter = func(_ io.Reader, _ io.Writer, _ *oauth2.Config) (string, error) {
		t.Fatal("did not expect auth prompt when token is cached")
		return "", nil
	}
	authCodeExchanger = func(_ context.Context, _ *oauth2.Config, _ string) (*oauth2.Token, error) {
		t.Fatal("did not expect token exchange when token is cached")
		return nil, nil
	}

	result, err := EnsureGoogleAuthClient(context.Background(), GoogleAuthOptions{
		CredentialsPath: credentialsPath,
		TokenPath:       tokenPath,
	})
	if err != nil {
		t.Fatalf("EnsureGoogleAuthClient failed: %v", err)
	}

	if !result.UsedCachedToken {
		t.Fatal("expected cached token to be used")
	}
}

func TestGoogleAuth_RejectsTokenFileWithLoosePermissions(t *testing.T) {
	tmp := t.TempDir()
	tokenPath := filepath.Join(tmp, "token.json")

	if err := os.WriteFile(tokenPath, []byte(`{"access_token":"abc"}`), 0o644); err != nil {
		t.Fatalf("write token file: %v", err)
	}

	_, err := tokenFromFile(tokenPath)
	if err == nil {
		t.Fatal("expected token permission error")
	}

	var authErr *GoogleAuthError
	if !errors.As(err, &authErr) {
		t.Fatalf("expected GoogleAuthError, got %T", err)
	}
	if authErr.Kind != GoogleAuthErrorTokenPermissions {
		t.Fatalf("expected kind %q, got %q", GoogleAuthErrorTokenPermissions, authErr.Kind)
	}
}

func TestGoogleAuth_MapsInvalidCredentialAndTokenErrors(t *testing.T) {
	t.Run("invalid credentials json", func(t *testing.T) {
		tmp := t.TempDir()
		credentialsPath := filepath.Join(tmp, "credentials.json")
		if err := os.WriteFile(credentialsPath, []byte(`{"bad":"json"}`), 0o644); err != nil {
			t.Fatalf("write credentials file: %v", err)
		}

		originalParse := parseGoogleCredentials
		t.Cleanup(func() { parseGoogleCredentials = originalParse })
		parseGoogleCredentials = func(_ []byte, _ ...string) (*oauth2.Config, error) {
			return nil, errors.New("invalid credentials payload")
		}

		_, err := EnsureGoogleAuthClient(context.Background(), GoogleAuthOptions{CredentialsPath: credentialsPath, TokenPath: filepath.Join(tmp, "token.json")})
		if err == nil {
			t.Fatal("expected invalid credentials error")
		}

		var authErr *GoogleAuthError
		if !errors.As(err, &authErr) {
			t.Fatalf("expected GoogleAuthError, got %T", err)
		}
		if authErr.Kind != GoogleAuthErrorInvalidCredentials {
			t.Fatalf("expected kind %q, got %q", GoogleAuthErrorInvalidCredentials, authErr.Kind)
		}
	})

	t.Run("invalid cached token json", func(t *testing.T) {
		tmp := t.TempDir()
		tokenPath := filepath.Join(tmp, "token.json")
		if err := os.WriteFile(tokenPath, []byte(`not-json`), 0o600); err != nil {
			t.Fatalf("write token file: %v", err)
		}

		_, err := tokenFromFile(tokenPath)
		if err == nil {
			t.Fatal("expected invalid token error")
		}

		var authErr *GoogleAuthError
		if !errors.As(err, &authErr) {
			t.Fatalf("expected GoogleAuthError, got %T", err)
		}
		if authErr.Kind != GoogleAuthErrorInvalidToken {
			t.Fatalf("expected kind %q, got %q", GoogleAuthErrorInvalidToken, authErr.Kind)
		}
	})
}

func TestGoogleAuth_ReturnsMissingCredentialError(t *testing.T) {
	_, err := EnsureGoogleAuthClient(context.Background(), GoogleAuthOptions{CredentialsPath: ""})
	if err == nil {
		t.Fatal("expected missing credentials error")
	}

	var authErr *GoogleAuthError
	if !errors.As(err, &authErr) {
		t.Fatalf("expected GoogleAuthError, got %T", err)
	}
	if authErr.Kind != GoogleAuthErrorMissingCredentials {
		t.Fatalf("expected kind %q, got %q", GoogleAuthErrorMissingCredentials, authErr.Kind)
	}
	if !strings.Contains(authErr.Error(), "GOOGLE_OAUTH_CREDENTIALS_PATH") {
		t.Fatalf("expected actionable message, got: %s", authErr.Error())
	}
}

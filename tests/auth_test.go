package tests

import (
	"context"
	"testing"
	"time"

	"github.com/tolstovrob/knocknock"
)

func TestAuth(t *testing.T) {
	store := knocknock.HandleMemoryStore()
	auth := knocknock.HandleAuth(store)
	ctx := context.Background()

	t.Run("CreateSession", func(t *testing.T) {
		userData := "test-user"
		session, err := auth.CreateSession(ctx, userData)
		if err != nil {
			t.Fatalf("CreateSession failed: %v", err)
		}

		if session.Token == "" {
			t.Error("Session token should not be empty")
		}

		if session.UserData != userData {
			t.Errorf("Expected userData %v, got %v", userData, session.UserData)
		}

		if session.IsExpired() {
			t.Error("New session should not be expired")
		}
	})

	t.Run("GetSession", func(t *testing.T) {
		userData := "test-user"
		session, err := auth.CreateSession(ctx, userData)
		if err != nil {
			t.Fatalf("CreateSession failed: %v", err)
		}

		retrieved, err := auth.GetSession(ctx, session.Token)
		if err != nil {
			t.Fatalf("GetSession failed: %v", err)
		}

		if retrieved.Token != session.Token {
			t.Errorf("Expected token %s, got %s", session.Token, retrieved.Token)
		}
	})

	t.Run("GetSession expired", func(t *testing.T) {
		store := knocknock.HandleMemoryStore()
		auth := knocknock.HandleAuth(store, knocknock.WithDefaultExpiry(-time.Hour))

		session, err := auth.CreateSession(ctx, "user")
		if err != nil {
			t.Fatalf("CreateSession failed: %v", err)
		}

		_, err = auth.GetSession(ctx, session.Token)
		if err == nil {
			t.Error("Expected error for expired session")
		}
	})

	t.Run("DeleteSession", func(t *testing.T) {
		session, err := auth.CreateSession(ctx, "user")
		if err != nil {
			t.Fatalf("CreateSession failed: %v", err)
		}

		err = auth.DeleteSession(ctx, session.Token)
		if err != nil {
			t.Fatalf("DeleteSession failed: %v", err)
		}

		_, err = auth.GetSession(ctx, session.Token)
		if err == nil {
			t.Error("Session should be deleted")
		}
	})

	t.Run("UpdateAuthOptions all options", func(t *testing.T) {
		auth := knocknock.HandleAuth(store)

		newTokenSize := 64
		newExpiry := 2 * time.Hour
		newCookieName := "new-cookie"
		newHeaderName := "X-Auth-Token"
		newQueryParamName := "auth-token"

		auth.UpdateAuthOptions(
			knocknock.WithTokenSize(newTokenSize),
			knocknock.WithDefaultExpiry(newExpiry),
			knocknock.WithCookieName(newCookieName),
			knocknock.WithHeaderName(newHeaderName),
			knocknock.WithQueryParamName(newQueryParamName),
		)

		if auth.AuthOptions.TokenSize != newTokenSize {
			t.Errorf("Expected TokenSize %d, got %d", newTokenSize, auth.AuthOptions.TokenSize)
		}

		if auth.AuthOptions.DefaultExpiry != newExpiry {
			t.Errorf("Expected DefaultExpiry %v, got %v", newExpiry, auth.AuthOptions.DefaultExpiry)
		}

		if auth.AuthOptions.CookieName != newCookieName {
			t.Errorf("Expected CookieName %s, got %s", newCookieName, auth.AuthOptions.CookieName)
		}

		if auth.AuthOptions.HeaderName != newHeaderName {
			t.Errorf("Expected HeaderName %s, got %s", newHeaderName, auth.AuthOptions.HeaderName)
		}

		if auth.AuthOptions.QueryParamName != newQueryParamName {
			t.Errorf("Expected QueryParamName %s, got %s", newQueryParamName, auth.AuthOptions.QueryParamName)
		}
	})

	t.Run("CreateSession with custom token size", func(t *testing.T) {
		customTokenSize := 16
		auth := knocknock.HandleAuth(store, knocknock.WithTokenSize(customTokenSize))

		session, err := auth.CreateSession(ctx, "test-user")
		if err != nil {
			t.Fatalf("CreateSession failed: %v", err)
		}

		expectedTokenLength := customTokenSize * 2
		if len(session.Token) != expectedTokenLength {
			t.Errorf("Expected token length %d, got %d", expectedTokenLength, len(session.Token))
		}
	})

	t.Run("CreateSession with custom expiry", func(t *testing.T) {
		customExpiry := 30 * time.Minute
		auth := knocknock.HandleAuth(store, knocknock.WithDefaultExpiry(customExpiry))

		session, err := auth.CreateSession(ctx, "test-user")
		if err != nil {
			t.Fatalf("CreateSession failed: %v", err)
		}

		expectedExpiry := session.CreatedAt.Add(customExpiry)
		if !session.ExpiresAt.Equal(expectedExpiry) {
			t.Errorf("Expected expiry %v, got %v", expectedExpiry, session.ExpiresAt)
		}
	})

	t.Run("Default options", func(t *testing.T) {
		auth := knocknock.HandleAuth(store)

		defaultOpts := auth.AuthOptions

		if defaultOpts.TokenSize != 32 {
			t.Errorf("Expected default TokenSize 32, got %d", defaultOpts.TokenSize)
		}

		if defaultOpts.DefaultExpiry != 24*time.Hour {
			t.Errorf("Expected default DefaultExpiry 24h, got %v", defaultOpts.DefaultExpiry)
		}

		if defaultOpts.CookieName != "session_token" {
			t.Errorf("Expected default CookieName 'session_token', got %s", defaultOpts.CookieName)
		}

		if defaultOpts.HeaderName != "Authorization" {
			t.Errorf("Expected default HeaderName 'Authorization', got %s", defaultOpts.HeaderName)
		}

		if defaultOpts.QueryParamName != "token" {
			t.Errorf("Expected default QueryParamName 'token', got %s", defaultOpts.QueryParamName)
		}
	})

	t.Run("Multiple updates", func(t *testing.T) {
		auth := knocknock.HandleAuth(store)

		auth.UpdateAuthOptions(
			knocknock.WithTokenSize(16),
			knocknock.WithCookieName("first-cookie"),
		)

		if auth.AuthOptions.TokenSize != 16 {
			t.Error("First update failed for TokenSize")
		}
		if auth.AuthOptions.CookieName != "first-cookie" {
			t.Error("First update failed for CookieName")
		}

		auth.UpdateAuthOptions(
			knocknock.WithTokenSize(64),
			knocknock.WithHeaderName("X-Custom-Auth"),
		)

		if auth.AuthOptions.TokenSize != 64 {
			t.Error("Second update failed for TokenSize")
		}
		if auth.AuthOptions.HeaderName != "X-Custom-Auth" {
			t.Error("Second update failed for HeaderName")
		}

		if auth.AuthOptions.CookieName != "first-cookie" {
			t.Error("CookieName should remain unchanged after second update")
		}
	})

	t.Run("Empty options", func(t *testing.T) {
		auth := knocknock.HandleAuth(store)
		originalOpts := *auth.AuthOptions

		auth.UpdateAuthOptions()

		if *auth.AuthOptions != originalOpts {
			t.Error("AuthOptions should not change with empty update")
		}
	})
}

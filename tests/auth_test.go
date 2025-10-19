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

	t.Run("UpdateAuthOptions", func(t *testing.T) {
		auth := knocknock.HandleAuth(store)

		newExpiry := 2 * time.Hour
		newCookieName := "new-cookie"

		auth.UpdateAuthOptions(
			knocknock.WithDefaultExpiry(newExpiry),
			knocknock.WithCookieName(newCookieName),
		)

		if auth.AuthOptions.DefaultExpiry != newExpiry {
			t.Errorf("Expected DefaultExpiry %v, got %v", newExpiry, auth.AuthOptions.DefaultExpiry)
		}

		if auth.AuthOptions.CookieName != newCookieName {
			t.Errorf("Expected CookieName %s, got %s", newCookieName, auth.AuthOptions.CookieName)
		}
	})
}

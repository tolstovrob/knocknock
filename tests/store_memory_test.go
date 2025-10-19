package tests

import (
	"context"
	"testing"
	"time"

	"github.com/tolstovrob/knocknock"
)

func TestMemoryStore(t *testing.T) {
	store := knocknock.HandleMemoryStore()
	ctx := context.Background()

	session := &knocknock.Session{
		Token:     "test-token",
		UserData:  "test-user",
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Hour),
	}

	t.Run("Save and Get", func(t *testing.T) {
		err := store.Save(ctx, session)
		if err != nil {
			t.Fatalf("Save failed: %v", err)
		}

		retrieved, err := store.Get(ctx, "test-token")
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}

		if retrieved.Token != session.Token {
			t.Errorf("Expected token %s, got %s", session.Token, retrieved.Token)
		}
	})

	t.Run("Save duplicate", func(t *testing.T) {
		err := store.Save(ctx, session)
		if err == nil {
			t.Error("Expected error for duplicate session")
		}
	})

	t.Run("Get non-existent", func(t *testing.T) {
		_, err := store.Get(ctx, "non-existent")
		if err == nil {
			t.Error("Expected error for non-existent session")
		}
	})

	t.Run("Delete", func(t *testing.T) {
		err := store.Delete(ctx, "test-token")
		if err != nil {
			t.Fatalf("Delete failed: %v", err)
		}

		_, err = store.Get(ctx, "test-token")
		if err == nil {
			t.Error("Session should be deleted")
		}
	})

	t.Run("Cleanup", func(t *testing.T) {
		expiredSession := &knocknock.Session{
			Token:     "expired-token",
			UserData:  "expired-user",
			CreatedAt: time.Now().Add(-2 * time.Hour),
			ExpiresAt: time.Now().Add(-time.Hour),
		}

		store.Save(ctx, expiredSession)
		store.Save(ctx, session)

		store.Cleanup()

		_, err := store.Get(ctx, "expired-token")
		if err == nil {
			t.Error("Expired session should be cleaned up")
		}

		_, err = store.Get(ctx, session.Token)
		if err != nil {
			t.Error("Valid session should not be cleaned up")
		}
	})
}

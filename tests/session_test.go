package tests

import (
	"testing"
	"time"

	"github.com/tolstovrob/knocknock"
)

func TestSession(t *testing.T) {
	t.Run("MakeSession", func(t *testing.T) {
		userData := "test-user"
		expiresIn := time.Hour
		session := knocknock.MakeSession("token123", userData, expiresIn)

		if session.Token != "token123" {
			t.Errorf("Expected token token123, got %s", session.Token)
		}

		if session.UserData != userData {
			t.Errorf("Expected userData %v, got %v", userData, session.UserData)
		}

		if time.Since(session.CreatedAt) > time.Second {
			t.Error("CreatedAt should be approximately now")
		}

		expectedExpiry := session.CreatedAt.Add(expiresIn)
		if !session.ExpiresAt.Equal(expectedExpiry) {
			t.Errorf("Expected ExpiresAt %v, got %v", expectedExpiry, session.ExpiresAt)
		}
	})

	t.Run("IsExpired", func(t *testing.T) {
		validSession := knocknock.MakeSession("valid", "user", time.Hour)
		if validSession.IsExpired() {
			t.Error("Session should not be expired")
		}

		expiredSession := &knocknock.Session{
			Token:     "expired",
			UserData:  "user",
			CreatedAt: time.Now().Add(-2 * time.Hour),
			ExpiresAt: time.Now().Add(-time.Hour),
		}
		if !expiredSession.IsExpired() {
			t.Error("Session should be expired")
		}
	})
}

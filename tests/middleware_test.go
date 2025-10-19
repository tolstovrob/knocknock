package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tolstovrob/knocknock"
)

func TestMiddleware(t *testing.T) {
	store := knocknock.HandleMemoryStore()
	auth := knocknock.HandleAuth(store)
	ctx := context.Background()

	t.Run("Middleware with valid token", func(t *testing.T) {
		session, _ := auth.CreateSession(ctx, "test-user")

		handler := auth.Middleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sess := knocknock.GetSession(r.Context())
			if sess == nil {
				t.Error("Session should be in context")
			}
			if sess.Token != session.Token {
				t.Errorf("Expected token %s, got %s", session.Token, sess.Token)
			}
			w.WriteHeader(http.StatusOK)
		}))

		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer "+session.Token)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}
	})

	t.Run("Middleware with cookie", func(t *testing.T) {
		session, _ := auth.CreateSession(ctx, "test-user")

		handler := auth.Middleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sess := knocknock.GetSession(r.Context())
			if sess == nil {
				t.Error("Session should be in context")
			}
			w.WriteHeader(http.StatusOK)
		}))

		req := httptest.NewRequest("GET", "/", nil)
		req.AddCookie(&http.Cookie{
			Name:  auth.AuthOptions.CookieName,
			Value: session.Token,
		})
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}
	})

	t.Run("Middleware with query param", func(t *testing.T) {
		session, _ := auth.CreateSession(ctx, "test-user")

		handler := auth.Middleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sess := knocknock.GetSession(r.Context())
			if sess == nil {
				t.Error("Session should be in context")
			}
			w.WriteHeader(http.StatusOK)
		}))

		req := httptest.NewRequest("GET", "/?token="+session.Token, nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}
	})

	t.Run("Middleware without token", func(t *testing.T) {
		handler := auth.Middleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sess := knocknock.GetSession(r.Context())
			if sess != nil {
				t.Error("Session should not be in context")
			}
			w.WriteHeader(http.StatusOK)
		}))

		req := httptest.NewRequest("GET", "/", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}
	})

	t.Run("GetSession from context", func(t *testing.T) {
		session := &knocknock.Session{Token: "test-token", UserData: "test-user"}

		ctx := context.WithValue(context.Background(), knocknock.SessionContextKey, session)

		retrieved := knocknock.GetSession(ctx)
		if retrieved == nil {
			t.Error("Should retrieve session from context")
		}
		if retrieved.Token != session.Token {
			t.Errorf("Expected token %s, got %s", session.Token, retrieved.Token)
		}
	})

	t.Run("GetSession from empty context", func(t *testing.T) {
		retrieved := knocknock.GetSession(context.Background())
		if retrieved != nil {
			t.Error("Should not retrieve session from empty context")
		}
	})

	t.Run("Extract token from header without Bearer", func(t *testing.T) {
		session, _ := auth.CreateSession(ctx, "test-user")

		handler := auth.Middleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sess := knocknock.GetSession(r.Context())
			if sess == nil {
				t.Error("Session should be in context")
			}
			w.WriteHeader(http.StatusOK)
		}))

		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", session.Token) // Без "Bearer "
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}
	})
}

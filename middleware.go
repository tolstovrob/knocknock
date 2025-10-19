package knocknock

import (
	"context"
	"net/http"
	"strings"

	"github.com/tolstovrob/knocknock/sessions"
)

type contextKey string

const sessionContextKey contextKey = "session"

func (a *Auth) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := a.extractToken(r)

			if session, err := a.GetSession(r.Context(), token); err == nil {
				ctx := context.WithValue(r.Context(), sessionContextKey, session)
				r = r.WithContext(ctx)
			}

			next.ServeHTTP(w, r)
		})
	}
}

func GetSession(ctx context.Context) *sessions.Session {
	if session, ok := ctx.Value(sessionContextKey).(*sessions.Session); ok {
		return session
	}
	return nil
}

func (a *Auth) extractToken(r *http.Request) string {
	if authHeader := r.Header.Get(a.options.HeaderName); authHeader != "" {
		if strings.HasPrefix(authHeader, "Bearer ") {
			return authHeader[7:]
		}
		return authHeader
	}

	if token := r.URL.Query().Get(a.options.QueryParamName); token != "" {
		return token
	}

	if cookie, err := r.Cookie(a.options.CookieName); err == nil {
		return cookie.Value
	}

	return ""
}

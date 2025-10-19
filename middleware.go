package knocknock

import (
	"context"
	"net/http"
	"strings"
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

func GetSession(ctx context.Context) *Session {
	if session, ok := ctx.Value(sessionContextKey).(*Session); ok {
		return session
	}
	return nil
}

func (a *Auth) extractToken(r *http.Request) string {
	if authHeader := r.Header.Get(a.authOptions.HeaderName); authHeader != "" {
		if strings.HasPrefix(authHeader, "Bearer ") {
			return authHeader[7:]
		}
		return authHeader
	}

	if token := r.URL.Query().Get(a.authOptions.QueryParamName); token != "" {
		return token
	}

	if cookie, err := r.Cookie(a.authOptions.CookieName); err == nil {
		return cookie.Value
	}

	return ""
}

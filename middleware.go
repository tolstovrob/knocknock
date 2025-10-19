package knocknock

import (
	"context"
	"net/http"

	"github.com/tolstovrob/knocknock/sessions"
)

/*
 * Тут содержится код дял http middleware системы авторизации, а также соответствующая работа с контекстом
 */

type contextKey string

const sessionContextKey contextKey = "session"

func (a *Auth) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := extractToken(r)
			if session, err := a.GetSession(r.Context(), token); err == nil && token != "" {
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

func extractToken(r *http.Request) string {
	// Из заголовка HTTP
	if authHeader := r.Header.Get("Authorization"); authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}

	// Из query параметра
	if token := r.URL.Query().Get("token"); token != "" {
		return token
	}

	// Из cookie
	if cookie, err := r.Cookie("session_token"); err == nil {
		return cookie.Value
	}

	return ""
}

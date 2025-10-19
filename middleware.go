package knocknock

/*
 * middleware.go содержит HTTP middleware для проверки аутентификации.
 */

import (
	"context"
	"net/http"
	"strings"
)

type ContextKey string

const SessionContextKey ContextKey = "session"

// Создает HTTP middleware для проверки аутентификации. Функция извлекает токен из запроса и, если сессия
// валидна, добавляет её в контекст запроса.
//
// Пример:
//
//	router := mux.NewRouter()
//	router.Use(auth.Middleware())
//
//	router.HandleFunc("/protected", func(w http.ResponseWriter, r *http.Request) {
//	    session := knocknock.GetSession(r.Context())
//	    if session == nil {
//	        http.Error(w, "Unauthorized", http.StatusUnauthorized)
//	        return
//	    }
//	    // работа с авторизованным пользователем
//	})
func (a *Auth) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := a.extractToken(r)

			if session, err := a.GetSession(r.Context(), token); err == nil {
				ctx := context.WithValue(r.Context(), SessionContextKey, session)
				r = r.WithContext(ctx)
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Возвращает сессию из контекста запроса. Возвращает nil если сессия не найдена в контексте.
//
// Пример:
//
//	func handler(w http.ResponseWriter, r *http.Request) {
//	    session := knocknock.GetSession(r.Context())
//	    if session == nil {
//	        http.Error(w, "Unauthorized", http.StatusUnauthorized)
//	        return
//	    }
//	    userData := session.UserData.(*User)
//	    // работа с данными пользователя
//	}
func GetSession(ctx context.Context) *Session {
	if session, ok := ctx.Value(SessionContextKey).(*Session); ok {
		return session
	}
	return nil
}

// Извлекает токен из HTTP-запроса. Проверяет HTTP-заголовки, query-параметры и cookies.
func (a *Auth) extractToken(r *http.Request) string {
	if authHeader := r.Header.Get(a.AuthOptions.HeaderName); authHeader != "" {
		if strings.HasPrefix(authHeader, "Bearer ") {
			return authHeader[7:]
		}
		return authHeader
	}

	if token := r.URL.Query().Get(a.AuthOptions.QueryParamName); token != "" {
		return token
	}

	if cookie, err := r.Cookie(a.AuthOptions.CookieName); err == nil {
		return cookie.Value
	}

	return ""
}

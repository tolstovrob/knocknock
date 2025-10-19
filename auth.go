package knocknock

/*
 * Точка входа в проект. Здесь содержится всё, что связано непосредственно с авторизацией и аутентификацией. Ключевой
 * является структура Auth, которая управляет хранилищем и состоянием сессий
 */

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"
)

// Структура настроек Auth через функциональные опции
type AuthOptions struct {
	TokenSize      int           // Длина токена в байтах
	DefaultExpiry  time.Duration // Время жизни сессии по умолчанию
	CookieName     string        // Имя cookie для токена
	HeaderName     string        // Имя HTTP-заголовка для токена
	QueryParamName string        // Имя query-параметра для токена
}

type AuthOption func(*AuthOptions)

// Функциональная опция для установки длины токена
func WithTokenSize(length int) AuthOption {
	return func(o *AuthOptions) {
		o.TokenSize = length
	}
}

// Функциональная опция для установки времени жизни по умолчанию
func WithDefaultExpiry(expiry time.Duration) AuthOption {
	return func(o *AuthOptions) {
		o.DefaultExpiry = expiry
	}
}

// Функциональная опция для установки имя cookie токена
func WithCookieName(name string) AuthOption {
	return func(o *AuthOptions) {
		o.CookieName = name
	}
}

// Функциональная опция для установки имени HTTP-заголовка токена
func WithHeaderName(name string) AuthOption {
	return func(o *AuthOptions) {
		o.HeaderName = name
	}
}

// Функциональная опция для установки имени query-параметра токена
func WithQueryParamName(name string) AuthOption {
	return func(o *AuthOptions) {
		o.QueryParamName = name
	}
}

// Создаёт и возвращает конфигурацию Auth по умолчанию
func defaultAuthOptions() *AuthOptions {
	return &AuthOptions{
		TokenSize:      32,
		DefaultExpiry:  24 * time.Hour,
		CookieName:     "session_token",
		HeaderName:     "Authorization",
		QueryParamName: "token",
	}
}

// Структура для управления сессиями аутентификации
type Auth struct {
	store       Store
	authOptions *AuthOptions
}

// Конструктор структуры Auth. Обязательно принимает хранилище, опционально -- набор функциональных опций
//
// Пример:
//
//	store := NewMemoryStore()
//	auth := HandleAuth(store, knocknock.WithDefaultExpiry(2*time.Hour))
func HandleAuth(store Store, authOptions ...AuthOption) *Auth {
	opts := defaultAuthOptions()
	for _, opt := range authOptions {
		opt(opts)
	}
	return &Auth{store, opts}
}

// Конструктор для обновления опций Auth. Принимает набор функциональных опций
func (a *Auth) UpdateAuthOptions(authOptions ...AuthOption) {
	for _, opt := range authOptions {
		opt(a.authOptions)
	}
}

// Создаёт новую сессию для указанных данных. Автоматически генерирует токен сессии и устанавливает время истечения.
// Конфигурируется через опции сессии в session.go
func (a *Auth) CreateSession(ctx context.Context, userData UserData) (*Session, error) {
	token, err := generateToken(a.authOptions.TokenSize)
	if err != nil {
		return nil, err
	}

	session := MakeSession(token, userData, a.authOptions.DefaultExpiry)

	if err := a.store.Save(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}

// Возвращает сессию по токену. Автоматически удаляет сессию если она истекла и возвращает SessionExpiredError
func (a *Auth) GetSession(ctx context.Context, token string) (*Session, error) {
	session, err := a.store.Get(ctx, token)
	if err != nil {
		return nil, err
	}

	if session.IsExpired() {
		_ = a.DeleteSession(ctx, token)
		return nil, SessionExpiredError
	}

	return session, nil
}

// Удаляет сессию по токену
func (a *Auth) DeleteSession(ctx context.Context, token string) error {
	return a.store.Delete(ctx, token)
}

// Генерирует криптографически безопасный случайный токен
func generateToken(size int) (string, error) {
	bytes := make([]byte, size)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

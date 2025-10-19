package knocknock

/*
 * Точка входа в проект. Здесь содержится всё, что связано непосредственно с авторизацией и аутентификацией. Ключевой
 * является структура Auth, которая управляет хранилищем и состоянием сессий
 */

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"
)

// Ошибки аутентефикации
var (
	// Возвращается в случае протухшей сессии
	SessionExpiredError = errors.New("Session expired")
)

// Структура настроек Auth через функциональные опции
type Options struct {
	TokenLength    int           // Длина токена в байтах
	DefaultExpiry  time.Duration // Время жизни сессии по умолчанию
	CookieName     string        // Имя cookie для токена
	HeaderName     string        // Имя HTTP-заголовка для токена
	QueryParamName string        // Имя query-параметра для токена
}

type Option func(*Options)

// Функциональная опция для установки длины токена
func WithTokenLength(length int) Option {
	return func(o *Options) {
		o.TokenLength = length
	}
}

// Функциональная опция для установки времени жизни по умолчанию
func WithDefaultExpiry(expiry time.Duration) Option {
	return func(o *Options) {
		o.DefaultExpiry = expiry
	}
}

// Функциональная опция для установки имя cookie токена
func WithCookieName(name string) Option {
	return func(o *Options) {
		o.CookieName = name
	}
}

// Функциональная опция для установки имени HTTP-заголовка токена
func WithHeaderName(name string) Option {
	return func(o *Options) {
		o.HeaderName = name
	}
}

// Функциональная опция для установки имени query-параметра токена
func WithQueryParamName(name string) Option {
	return func(o *Options) {
		o.QueryParamName = name
	}
}

// Создаёт и возвращает конфигурацию Auth по умолчанию
func defaultOptions() *Options {
	return &Options{
		TokenLength:    32,
		DefaultExpiry:  24 * time.Hour,
		CookieName:     "session_token",
		HeaderName:     "Authorization",
		QueryParamName: "token",
	}
}

// Структура для управления сессиями аутентификации
type Auth struct {
	store   Store
	options *Options
}

// Конструктор структуры Auth. Обязательно принимает хранилище, опционально -- набор функциональных опций из options.go
//
// Пример:
//
//	store := NewMemoryStore()
//	auth := HandleAuth(store, knocknock.WithDefaultExpiry(2*time.Hour))
func HandleAuth(store Store, options ...Option) *Auth {
	opts := defaultOptions()
	for _, opt := range options {
		opt(opts)
	}
	return &Auth{store, opts}
}

// Конструктор для обновления опций Auth. Принимает набор функциональных опций из options.go
func (a *Auth) UpdateOptions(options ...Option) {
	for _, opt := range options {
		opt(a.options)
	}
}

// Создаёт новую сессию для указанных данных. Автоматически генерирует токен сессии и устанавливает время истечения.
// Конфигурируется через опции сессии в sessions/options.go. Потенциально может вернуть ошибку из sessions/go
func (a *Auth) CreateSession(ctx context.Context, userData UserData, options ...SessionOption) (*Session, error) {
	token, err := generateToken(a.options.TokenLength)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	session := &Session{
		Token:     token,
		UserData:  userData,
		CreatedAt: now,
		ExpiresAt: now.Add(a.options.DefaultExpiry),
	}

	for _, opt := range options {
		opt(session)
	}

	if err := a.store.Save(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}

// Возвращает сессию по токену. Автоматически удаляет сессию если она истекла и возвращает SessionExpiredError.
// Потенциально может вернуть ошибку из sessions/go
func (a *Auth) GetSession(ctx context.Context, token string) (*Session, error) {
	session, err := a.store.Get(ctx, token)
	if err != nil {
		return nil, err
	}

	if time.Now().After(session.ExpiresAt) {
		_ = a.DeleteSession(ctx, token)
		return nil, SessionExpiredError
	}

	return session, nil
}

// Удаляет сессию по токену. Потенциально может вернуть ошибку из sessions/go
func (a *Auth) DeleteSession(ctx context.Context, token string) error {
	return a.store.Delete(ctx, token)
}

// Генерирует криптографически безопасный случайный токен.
func generateToken(size int) (string, error) {
	bytes := make([]byte, size)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

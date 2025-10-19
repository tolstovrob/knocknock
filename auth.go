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

	"github.com/tolstovrob/knocknock/sessions"
	"github.com/tolstovrob/knocknock/stores"
)

// Ошибки аутентефикации
var (
	// Возвращается в случае протухшей сессии
	SessionExpiredError = errors.New("Session expired")
)

// Структура для управления сессиями аутентификации
type Auth struct {
	store   stores.Store
	options *Options
}

// Конструктор структуры Auth. Обязательно принимает хранилище, опционально -- набор функциональных опций из options.go
//
// Пример:
//
//	store := stores.NewMemoryStore()
//	auth := New(store, knocknock.WithDefaultExpiry(2*time.Hour))
func New(store stores.Store, options ...Option) *Auth {
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
// Конфигурируется через опции сессии в sessions/options.go. Потенциально может вернуть ошибку из sessions/sessions.go
func (a *Auth) CreateSession(ctx context.Context, userData sessions.UserData, options ...sessions.Option) (*sessions.Session, error) {
	token, err := generateToken()
	if err != nil {
		return nil, err
	}

	now := time.Now()

	session := &sessions.Session{
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
// Потенциально может вернуть ошибку из sessions/sessions.go
func (a *Auth) GetSession(ctx context.Context, token string) (*sessions.Session, error) {
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

// Удаляет сессию по токену. Потенциально может вернуть ошибку из sessions/sessions.go
func (a *Auth) DeleteSession(ctx context.Context, token string) error {
	return a.store.Delete(ctx, token)
}

func generateToken() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

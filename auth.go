package knocknock

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/tolstovrob/knocknock/sessions"
	"github.com/tolstovrob/knocknock/store"
)

/*
 * Было бы здорово инкапсулировать логику работы с Session. Для этого есть структура Auth, которая хранит в себе
 * приватный экземпляр хранилища (определённого в internal/store/store.go), а также реализует публичные аксессоры дляъ
 * взаимодействия сессий и хранилища.
 */

var (
	SessionExpiredError = errors.New("Session expired")
)

type Auth struct {
	store store.Store
}

func New(store store.Store) *Auth {
	return &Auth{store}
}

func (a *Auth) CreateSession(ctx context.Context, userData sessions.UserData, expiresIn time.Duration) (*sessions.Session, error) {
	token, err := generateToken()
	if err != nil {
		return nil, err
	}

	now := time.Now()

	session := &sessions.Session{
		Token:     token,
		UserData:  userData,
		CreatedAt: now,
		ExpiresAt: now.Add(expiresIn),
	}

	if err := a.store.Save(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}

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

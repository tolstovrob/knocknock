package stores

import (
	"context"
	"errors"

	"github.com/tolstovrob/knocknock/sessions"
)

var (
	SessionNotFoundError = errors.New("Session not found")
	SessionExistsError   = errors.New("Session with given token already exists")
)

type Store interface {
	Save(ctx context.Context, session *sessions.Session) error
	Get(ctx context.Context, token string) (*sessions.Session, error)
	Delete(ctx context.Context, token string) error
}

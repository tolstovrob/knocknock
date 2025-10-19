package store

import (
	"context"

	"github.com/tolstovrob/knocknock/sessions"
)

type Store interface {
	Save(ctx context.Context, session *sessions.Session) error
	Get(ctx context.Context, token string) (*sessions.Session, error)
	Delete(ctx context.Context, token string) error
}

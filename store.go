package knocknock

import (
	"context"
)

type Store interface {
	Save(ctx context.Context, session *Session) error
	Get(ctx context.Context, token string) (*Session, error)
	Delete(ctx context.Context, token string) error
}

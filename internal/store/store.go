package store

import (
	"context"

	"github.com/tolstovrob/knocknock/internal/core/sessions"
)

/*
 * Интерфейс для взаимодействия с базой данных. Программист может выбрать любую базу и реализовать свои методы для
 * работы с ней
 */

type Store interface {
	Save(ctx context.Context, session *sessions.Session) error
	Get(ctx context.Context, token string) (*sessions.Session, error)
	Delete(ctx context.Context, token string) error
}

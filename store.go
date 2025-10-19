package knocknock

/*
 * В store.go описан интерфейс для реализации методов работы с любым хранилищем. Библиотека поставляет лишь in-memory
 * хранилище (store_memory.go), оставляя остальное на усмотрение программиста
 */

import (
	"context"
)

// Интерфейс для хранилища сессий. Реализации должны обеспечивать сохранение, получение и удаление сессий.
type Store interface {
	Save(ctx context.Context, session *Session) error
	Get(ctx context.Context, token string) (*Session, error)
	Delete(ctx context.Context, token string) error
}

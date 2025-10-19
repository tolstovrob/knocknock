package knocknock

/*
 * Пример реализации интерфейса Store (store.go). Хранилище для сессий в рантайме. Подойдёт для тестирования, в прод
 * лучше не брать.
 */

import (
	"context"
	"sync"
	"time"
)

// Структура для хранилища, с ограничениями на чтение/запись через мьютексы
type MemoryStore struct {
	mu       sync.RWMutex
	sessions map[string]*Session
}

// Создаёт новое хранилище
func HandleMemoryStore() *MemoryStore {
	return &MemoryStore{
		sessions: make(map[string]*Session),
	}
}

// Реализация Store.Save
func (m *MemoryStore) Save(ctx context.Context, session *Session) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.sessions[session.Token]; exists {
		return SessionExistsError
	}

	m.sessions[session.Token] = session
	return nil
}

// Реализация Store.Get
func (m *MemoryStore) Get(ctx context.Context, token string) (*Session, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	session, exists := m.sessions[token]

	if !exists {
		return nil, SessionNotFoundError
	}

	return session, nil
}

// Реализация Store.Delete
func (m *MemoryStore) Delete(ctx context.Context, token string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.sessions, token)
	return nil
}

// Очищает хранилище от протухших сессий. Примечательно, что функция не реализует Store. Программист может дополнять
// свои хранилища методами не из Store
func (m *MemoryStore) Cleanup() {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	for token, session := range m.sessions {
		if now.After(session.ExpiresAt) {
			delete(m.sessions, token)
		}
	}
}

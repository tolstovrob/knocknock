package impl

import (
	"context"
	"sync"
	"time"

	"github.com/tolstovrob/knocknock/internal/core/auth"
	"github.com/tolstovrob/knocknock/internal/core/sessions"
)

/*
 * Реализация хранилища сессий во временной памяти. Имплементирует интерфейс Store из internal/store/store.go
 */

type MemoryStore struct {
	mu       sync.RWMutex
	sessions map[string]*sessions.Session
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		sessions: make(map[string]*sessions.Session),
	}
}

func (m *MemoryStore) Save(ctx context.Context, session *sessions.Session) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.sessions[session.Token]; exists {
		return auth.SessionExistsError
	}

	m.sessions[session.Token] = session
	return nil
}

func (m *MemoryStore) Get(ctx context.Context, token string) (*sessions.Session, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	session, exists := m.sessions[token]

	if !exists {
		return nil, auth.SessionNotFoundError
	}

	return session, nil
}

func (m *MemoryStore) Delete(ctx context.Context, token string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.sessions, token)
	return nil
}

/*
 * Также добавим опциональную для интерфейса Store и, тем не менее, полезную функцию Cleanup, удаляющую устаревшие
 * сессии из нашего хранилища. Замечу, что метод не возвращает error даже в качестве nil. Это потому, что данный метод
 * не является частью публичного интерфейса Store, а значит мы можем написать любую реализацию
 */

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

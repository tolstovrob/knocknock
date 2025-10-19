package store

import (
	"context"
	"sync"
	"time"

	"github.com/tolstovrob/knocknock/sessions"
)

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
		return SessionExistsError
	}

	m.sessions[session.Token] = session
	return nil
}

func (m *MemoryStore) Get(ctx context.Context, token string) (*sessions.Session, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	session, exists := m.sessions[token]

	if !exists {
		return nil, SessionNotFoundError
	}

	return session, nil
}

func (m *MemoryStore) Delete(ctx context.Context, token string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.sessions, token)
	return nil
}

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

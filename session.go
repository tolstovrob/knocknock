package knocknock

import "time"

type UserData = any

type SessionOption func(*Session)

func WithExpiresAt(expiry time.Time) SessionOption {
	return func(s *Session) {
		s.ExpiresAt = expiry
	}
}

type Session struct {
	Token     string    `json:"token"`
	UserData  UserData  `json:"userData"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}

func MakeSession(token string, userData UserData, options ...SessionOption) *Session {
	ss := &Session{Token: token, UserData: userData}
	for _, opt := range options {
		opt(ss)
	}
	return ss
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

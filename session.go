package knocknock

import "time"

type UserData = any

type Session struct {
	Token     string    `json:"token"`
	UserData  UserData  `json:"userData"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}

func MakeSession(token string, userData UserData, expiresIn time.Duration) *Session {
	now := time.Now()
	ss := &Session{Token: token, UserData: userData, CreatedAt: now, ExpiresAt: now.Add(expiresIn)}
	return ss
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

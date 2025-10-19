package sessions

import "time"

type UserData = any

type Session struct {
	Token     string    `json:"token"`
	UserData  UserData  `json:"userData"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

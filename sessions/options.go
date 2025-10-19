package sessions

import "time"

type Option func(*Session)

func WithExpiresAt(expiry time.Time) func(*Session) {
	return func(s *Session) {
		s.ExpiresAt = expiry
	}
}

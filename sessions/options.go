package sessions

import "time"

type Option func(*Session)

func WithExpiresAt(expiry time.Time) Option {
	return func(s *Session) {
		s.ExpiresAt = expiry
	}
}

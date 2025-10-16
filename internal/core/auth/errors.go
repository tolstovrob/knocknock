package auth

import "errors"

var (
	SessionNotFoundError = errors.New("Session not found")
	SessionExpiredError  = errors.New("Session expired")
)

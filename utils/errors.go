package utils

import "errors"

var (
	SessionNotFoundError = errors.New("Session not found")
	SessionExpiredError  = errors.New("Session expired")
	SessionExistsError   = errors.New("Session with given token already exists")
)

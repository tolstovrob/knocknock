package knocknock

import "errors"

var (
	// Возвращается в случае протухшей сессии
	SessionExpiredError = errors.New("Session expired")
	// Возвращается если сессии не существует
	SessionNotFoundError = errors.New("Session not found")
	// Возвращается если данный токен уже занят
	SessionExistsError = errors.New("Session with given token already exists")
)

package knocknock

/*
 * session.go содержит структуру сессии и сопутствующие методы
 */

import "time"

// UserData представляет пользовательские данные, хранимые в сессии. Может быть любого типа.
type UserData = any // TODO(tolstovrob): think about generics

// Структура сессии
type Session struct {
	Token     string    `json:"token"`     // Токен сессии
	UserData  UserData  `json:"userData"`  // Информация о пользователе
	CreatedAt time.Time `json:"createdAt"` // Время создания
	ExpiresAt time.Time `json:"expiresAt"` // Время истечения
}

// Создает новую сессию с указанным токеном, пользовательскими данными и сроком жизни. Важно: третий аргумент expiresIn
// отвечает именно за время жизни, а значит в структуре Session в поле createdAt будет текущее время, а в expiresAt --
// текущее время + expiresIn.
//
// Пример:
//
//	session := sessions.New("token123", userData, 6*time.Hour)
func MakeSession(token string, userData UserData, expiresIn time.Duration) *Session {
	now := time.Now()
	ss := &Session{Token: token, UserData: userData, CreatedAt: now, ExpiresAt: now.Add(expiresIn)}
	return ss
}

// Проверяет, истекла ли сессия
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

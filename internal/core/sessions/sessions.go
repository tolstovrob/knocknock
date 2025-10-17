package sessions

import "time"

/*
 * Пакет internal/core -- точка входа в проект. Здесь содержится основная информация обо всех структурах и интерфейсах,
 * которые использует библиотека в работе.
 *
 * Итак, у нас есть тип сессии Session, который содержит в себе объект пользовательских данных UserData, а также
 * необходимые для сессии метаданные
 */

type UserData = any

type Session struct {
	Token     string    `json:"token"`
	UserData  UserData  `json:"userData"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}

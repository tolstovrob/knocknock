package sessions

import "time"

/*
 * Пакет internal/core -- точка входа в проект. Здесь содержится основная информация обо всех структурах и интерфейсах,
 * которые использует библиотека в работе.
 *
 * Итак, у нас есть тип сессии Session[UserData], который является дженериком от типа пользовательских данных.
 */

type UserData = any

type Session struct {
	Token     string    `json:"token"`
	UserData  UserData  `json:"userData"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}

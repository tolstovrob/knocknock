package knocknock

/*
 * Здесь содержится всё, что связано с конфигурацией объекта Auth через функциональные опции
 */

import "time"

// Структура настроек Auth через функциональные опции
type Options struct {
	TokenLength    int           // Длина токена в байтах
	DefaultExpiry  time.Duration // Время жизни сессии по умолчанию
	CookieName     string        // Имя cookie для токена
	HeaderName     string        // Имя HTTP-заголовка для токена
	QueryParamName string        // Имя query-параметра для токена
}

type Option func(*Options)

// Функциональная опция для установки длины токена
func WithTokenLength(length int) Option {
	return func(o *Options) {
		o.TokenLength = length
	}
}

// Функциональная опция для установки времени жизни по умолчанию
func WithDefaultExpiry(expiry time.Duration) Option {
	return func(o *Options) {
		o.DefaultExpiry = expiry
	}
}

// Функциональная опция для установки имя cookie токена
func WithCookieName(name string) Option {
	return func(o *Options) {
		o.CookieName = name
	}
}

// Функциональная опция для установки имени HTTP-заголовка токена
func WithHeaderName(name string) Option {
	return func(o *Options) {
		o.HeaderName = name
	}
}

// Функциональная опция для установки имени query-параметра токена
func WithQueryParamName(name string) Option {
	return func(o *Options) {
		o.QueryParamName = name
	}
}

// Создаёт и возвращает конфигурацию Auth по умолчанию
func defaultOptions() *Options {
	return &Options{
		TokenLength:    32,
		DefaultExpiry:  24 * time.Hour,
		CookieName:     "session_token",
		HeaderName:     "Authorization",
		QueryParamName: "token",
	}
}

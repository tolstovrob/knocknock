package knocknock

import "time"

type Options struct {
	TokenLength    int
	DefaultExpiry  time.Duration
	CookieName     string
	HeaderName     string
	QueryParamName string
}

type Option func(*Options)

func WithTokenLength(length int) Option {
	return func(o *Options) {
		o.TokenLength = length
	}
}

func WithDefaultExpiry(expiry time.Duration) Option {
	return func(o *Options) {
		o.DefaultExpiry = expiry
	}
}

func WithCookieName(name string) Option {
	return func(o *Options) {
		o.CookieName = name
	}
}

func WithHeaderName(name string) Option {
	return func(o *Options) {
		o.HeaderName = name
	}
}

func WithQueryParamName(name string) Option {
	return func(o *Options) {
		o.QueryParamName = name
	}
}

func defaultOptions() *Options {
	return &Options{
		TokenLength:    32,
		DefaultExpiry:  24 * time.Hour,
		CookieName:     "session_token",
		HeaderName:     "Authorization",
		QueryParamName: "token",
	}
}

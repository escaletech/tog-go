package sessions

import (
	"context"
)

type Session map[string]bool

func (s Session) IsSet(name string) bool {
	value, _ := s[name]
	return value
}

type SessionOptions struct {
	Force  Session
	Traits []string
}

type ErrorHandler func(context.Context, error)

type ClientOptions struct {
	Addr    string
	Cluster bool
	OnError ErrorHandler
}

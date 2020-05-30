package sessions

import (
	"context"
)

type Session map[string]bool

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

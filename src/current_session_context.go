package main

import "context"

var key = ctxKey("current-session")

type ctxKey string

func WithCurrentSession(ctx context.Context, curr *CurrentSession) context.Context {
	return context.WithValue(ctx, key, curr)
}

func CurrentSessionFromContext(ctx context.Context) *CurrentSession {
	curr, ok := ctx.Value(key).(*CurrentSession)
	if !ok {
		panic("no such current session")
	}
	return curr
}

package main

import (
	"context"
	"errors"
)

const PostgresFmt = "2006-01-02 15:04:05.000000Z"

var Rx = RootRx{}

var (
	ErrUserMustBeUnauth = errors.New("user must be unauthenticated")
	ErrUserMustBeAuth   = errors.New("user must be authenticated")
)

type RootRx struct{}

func (r *RootRx) Ping(ctx context.Context) bool {
	return true
}

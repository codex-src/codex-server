package main

import (
	"context"
	"time"

	graphql "github.com/graph-gophers/graphql-go"
)

type User struct {
	UserID    graphql.ID
	CreatedAt time.Time
	UpdatedAt time.Time
	Username  string

	// TODO
	Notes []*Note
}

func (r *RootRx) TestUsernameTaken(ctx context.Context, args struct{ Username string }) (bool, error) {
	currUser := CurrentSessionFromContext(ctx)
	if !currUser.IsUnauth() {
		return false, ErrUserMustBeUnauth
	}
	var ok bool
	err := DB.QueryRow(`
		select count(*) = 0
		from users
		where username = $1
	`, args.Username).Scan(&ok)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (r *RootRx) Me(ctx context.Context) (*UserRx, error) {
	currUser := CurrentSessionFromContext(ctx)
	if !currUser.IsAuth() {
		return nil, ErrUserMustBeAuth
	}
	user := &User{}
	err := DB.QueryRow(`
		select
			user_id,
			created_at,
			updated_at,
			username
		from users
		where user_id = $1
	`, currUser.UserID).Scan(&user.UserID, &user.CreatedAt, &user.UpdatedAt, &user.Username)
	if err != nil {
		return nil, err
	}
	return &UserRx{user}, nil
}

/*
 * UserRx
 */

type UserRx struct{ user *User }

func (r *UserRx) UserID() graphql.ID {
	return r.user.UserID
}

func (r *UserRx) CreatedAt() string {
	return r.user.CreatedAt.UTC().Format(PostgresFmt)
}

func (r *UserRx) UpdatedAt() string {
	return r.user.UpdatedAt.UTC().Format(PostgresFmt)
}

func (r *UserRx) Username() string {
	return r.user.Username
}

func (r *UserRx) Notes(ctx context.Context, args struct{ Limit, Offset *int32 }) ([]*NoteRx, error) {
	return Rx.Notes(ctx, args)
}

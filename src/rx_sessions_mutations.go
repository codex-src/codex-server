package main

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

func (r *RootRx) CreateSession(ctx context.Context, args struct{ Username, Password string }) (*UserRx, error) {
	currUser := CurrentSessionFromContext(ctx)
	if !currUser.IsUnauth() {
		return nil, ErrUserMustBeUnauth
	}
	tx, err := DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	var dbPasswordHash []byte
	err = tx.QueryRow(`
		select
			user_id,
			password_hash
		from users
		where username = $1
	`, args.Username).Scan(&currUser.UserID, &dbPasswordHash)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword(dbPasswordHash, []byte(args.Password))
	if err != nil {
		return nil, err
	}
	err = tx.QueryRow(`
		insert into sessions (
			user_id,
			expires_at )
		values ( $1, now() + '1 week' )
		returning
			session_id,
			expires_at
	`, currUser.UserID).Scan(&currUser.SessionID, &currUser.SessionExpiry)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	SetCookie(currUser.W, SESSION_KEY, string(currUser.SessionID), currUser.SessionExpiry)
	return Rx.Me(ctx)
}

func (r *RootRx) RevokeSession(ctx context.Context) (bool, error) {
	currUser := CurrentSessionFromContext(ctx)
	if !currUser.IsAuth() {
		return false, ErrUserMustBeAuth
	}
	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}
	defer tx.Rollback()
	_, err = tx.Exec(`
		update sessions
		set revoked_at = now()
		where session_id = $1
	`, currUser.SessionID)
	if err != nil {
		return false, err
	}
	err = tx.Commit()
	if err != nil {
		return false, err
	}
	ClearCookie(currUser.W, SESSION_KEY, string(currUser.SessionID))
	return true, nil
}

package main

import (
	"net/http"
	"time"

	graphql "github.com/graph-gophers/graphql-go"
)

type CurrentSession struct {
	W http.ResponseWriter
	R *http.Request

	UserID        graphql.ID
	SessionID     graphql.ID
	SessionExpiry time.Time
}

func (c *CurrentSession) IsUnauth() bool {
	return c.UserID == "" && c.SessionID == "" && c.SessionExpiry.IsZero()
}

func (c *CurrentSession) IsAuth() bool {
	return c.UserID != "" || c.SessionID != "" || c.SessionExpiry.IsZero()
}

func ExtendCurrentSession(w http.ResponseWriter, r *http.Request) (*CurrentSession, error) {
	currUser := &CurrentSession{W: w, R: r}
	cookie, err := r.Cookie(SESSION_KEY)
	if err == http.ErrNoCookie {
		return currUser, nil
	}
	tx, err := DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	err = tx.QueryRow(`
		update sessions
		set expires_at = now() + '1 week'
		where
			session_id = $1 and
			( select expires_at > now() and revoked_at is null where session_id = $2 )
		returning
			user_id,
			session_id,
			expires_at
	`, cookie.Value, cookie.Value).Scan(&currUser.UserID, &currUser.SessionID, &currUser.SessionExpiry)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return currUser, nil
}

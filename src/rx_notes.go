package main

import (
	"context"
	"fmt"
	"time"

	graphql "github.com/graph-gophers/graphql-go"
)

type Note struct {
	UserID         graphql.ID
	NoteID         graphql.ID
	CreatedAt      time.Time
	UpdatedAt      time.Time
	TitleUTF8Count int32
	Title          string
	DataUTF8Count  int32
	Data           string
}

func (r *RootRx) Notes(ctx context.Context, args struct{ Limit, Offset *int32 }) ([]*NoteRx, error) {
	currUser := CurrentSessionFromContext(ctx)
	if !currUser.IsAuth() {
		return nil, ErrUserMustBeAuth
	}
	var rxs []*NoteRx
	rows, err := DB.Query(`
		select
			user_id,
			note_id,
			created_at,
			updated_at,
			title_utf8_count,
			title,
			data_utf8_count,
			data
		from notes
		where user_id = $1
		order by updated_at desc
		limit coalesce( $2, 25 )
		offset $3
	`, currUser.UserID, args.Limit, args.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		note := &Note{}
		err := rows.Scan(&note.UserID, &note.NoteID, &note.CreatedAt, &note.UpdatedAt, &note.TitleUTF8Count, &note.Title, &note.DataUTF8Count, &note.Data)
		if err != nil {
			return nil, err
		}
		rxs = append(rxs, &NoteRx{note})
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return rxs, nil
}

func (r *RootRx) Note(ctx context.Context, args struct{ NoteID graphql.ID }) (*NoteRx, error) {
	currUser := CurrentSessionFromContext(ctx)
	if !currUser.IsAuth() {
		return nil, ErrUserMustBeAuth
	}
	note := &Note{}
	err := DB.QueryRow(`
		select
			user_id,
			note_id,
			created_at,
			updated_at,
			title_utf8_count,
			title,
			data_utf8_count,
			data
		from notes
		where
			note_id = $1 and
			( select user_id = $2 from notes where note_id = $3 )
	`, args.NoteID, currUser.UserID, args.NoteID).Scan(&note.UserID, &note.NoteID, &note.CreatedAt, &note.UpdatedAt, &note.TitleUTF8Count, &note.Title, &note.DataUTF8Count, &note.Data)
	if err != nil {
		return nil, err
	}
	return &NoteRx{note}, nil
}

/*
 * NoteRx
 */

type NoteRx struct{ note *Note }

func (r *NoteRx) UserID() graphql.ID {
	return r.note.UserID
}

func (r *NoteRx) NoteID() graphql.ID {
	return r.note.NoteID
}

func (r *NoteRx) CreatedAt() string {
	return r.note.CreatedAt.UTC().Format(PostgresFmt)
}

func (r *NoteRx) UpdatedAt() string {
	return r.note.UpdatedAt.UTC().Format(PostgresFmt)
}

func (r *NoteRx) TitleUTF8Count() int32 {
	return r.note.TitleUTF8Count
}

func (r *NoteRx) Title() string {
	return r.note.Title
}

func (r *NoteRx) DataUTF8Count() int32 {
	return r.note.DataUTF8Count
}

func (r *NoteRx) Data280() string {
	return fmt.Sprintf("%0.*s", 280, r.note.Data)
}

func (r *NoteRx) Data() string {
	return r.note.Data
}

package main

import (
	"context"
	"unicode/utf8"

	graphql "github.com/graph-gophers/graphql-go"
)

func (r *RootRx) CreateNote(ctx context.Context, args struct{ Title, Data string }) (*NoteRx, error) {
	currUser := CurrentSessionFromContext(ctx)
	if !currUser.IsAuth() {
		return nil, ErrUserMustBeAuth
	}
	tx, err := DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	var noteID graphql.ID
	err = tx.QueryRow(`
		insert into notes (
			user_id,
			title_utf8_count,
			title,
			data_utf8_count,
			data )
		values ( $1, $2, $3, $4, $5 )
		returning note_id
	`, currUser.UserID, utf8.RuneCountInString(args.Title), args.Title, utf8.RuneCountInString(args.Data), args.Data).Scan(&noteID)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	noteArgs := struct{ NoteID graphql.ID }{noteID}
	return Rx.Note(ctx, noteArgs)
}

type UpdateNoteArgs struct {
	NoteID graphql.ID
	Title  string
	Data   string
}

func (r *RootRx) UpdateNote(ctx context.Context, args UpdateNoteArgs) (*NoteRx, error) {
	currUser := CurrentSessionFromContext(ctx)
	if !currUser.IsAuth() {
		return nil, ErrUserMustBeAuth
	}
	tx, err := DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	_, err = tx.Exec(`
		update notes
		set
			title_utf8_count = $1,
			title = $2,
			data_utf8_count = $3,
			data = $4
		where
			note_id = $5 and
			( select user_id = $6 from notes where note_id = $7 )
	`, utf8.RuneCountInString(args.Title), args.Title, utf8.RuneCountInString(args.Data), args.Data, args.NoteID, currUser.UserID, args.NoteID)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	noteArgs := struct{ NoteID graphql.ID }{args.NoteID}
	return Rx.Note(ctx, noteArgs)
}

type UpdateNoteTitleArgs struct {
	NoteID graphql.ID
	Title  string
}

func (r *RootRx) UpdateNoteTitle(ctx context.Context, args UpdateNoteTitleArgs) (*NoteRx, error) {
	currUser := CurrentSessionFromContext(ctx)
	if !currUser.IsAuth() {
		return nil, ErrUserMustBeAuth
	}
	tx, err := DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	_, err = tx.Exec(`
		update notes
		set
			title_utf8_count = $1,
			title = $2
		where
			note_id = $3 and
			( select user_id = $4 from notes where note_id = $5 )
	`, utf8.RuneCountInString(args.Title), args.Title, args.NoteID, currUser.UserID, args.NoteID)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	noteArgs := struct{ NoteID graphql.ID }{args.NoteID}
	return Rx.Note(ctx, noteArgs)
}

func (r *RootRx) DuplicateNote(ctx context.Context, args struct{ NoteID graphql.ID }) (*NoteRx, error) {
	currUser := CurrentSessionFromContext(ctx)
	if !currUser.IsAuth() {
		return nil, ErrUserMustBeAuth
	}
	tx, err := DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	var noteID graphql.ID
	err = tx.QueryRow(`
		insert into notes (
			user_id,
			title_utf8_count,
			title,
			data_utf8_count,
			data )
		select
			user_id,
			title_utf8_count,
			title,
			data_utf8_count,
			data
		from notes
		where
			note_id = $1 and
			( select user_id = $2 from notes where note_id = $3 )
		returning note_id
	`, args.NoteID, currUser.UserID, args.NoteID).Scan(&noteID)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	noteArgs := struct{ NoteID graphql.ID }{noteID}
	return Rx.Note(ctx, noteArgs)
}

func (r *RootRx) DeleteNote(ctx context.Context, args struct{ NoteID graphql.ID }) (bool, error) {
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
		delete
		from notes
		where
			note_id = $1 and
			( select user_id = $2 from notes where note_id = $3 )
	`, args.NoteID, currUser.UserID, args.NoteID)
	if err != nil {
		return false, err
	}
	err = tx.Commit()
	if err != nil {
		return false, err
	}
	return true, nil
}

package main

import (
	"context"
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

type CreateUserInput struct {
	Username           string
	Password           string
	Passcode           string
	ChargeMonth        bool
	StripeCardID       string
	StripeCardBrand    string
	StripeCardLastFour string
}

func (r *RootRx) CreateUser(ctx context.Context, args struct{ User *CreateUserInput }) (*UserRx, error) {
	currUser := CurrentSessionFromContext(ctx)
	if !currUser.IsUnauth() {
		return nil, ErrUserMustBeUnauth
	}
	tx, err := DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(args.User.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	err = tx.QueryRow(`
		insert into users (
			username,
			password_hash,
			passcode )
		values ( $1, $2, $3 )
		returning user_id
	`, args.User.Username, passwordHash, args.User.Passcode).Scan(&currUser.UserID)
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
	_, err = tx.Exec(`
		insert into billing (
			user_id,
			charge_month,
			stripe_card_id,
			stripe_card_brand,
			stripe_card_last_four )
		values ( $1, $2, $3, $4, $5 )
	`, currUser.UserID, args.User.ChargeMonth, args.User.StripeCardID, args.User.StripeCardBrand, args.User.StripeCardLastFour)
	if err != nil {
		return nil, err
	}
	// FIXME: Stripe.
	_, err = tx.Exec(`
		insert into subscriptions (
			user_id,
			charge_month,
			start_date,
			end_date,
			stripe_charge_id,
			stripe_charge_amount,
			stripe_card_brand,
			stripe_card_last_four )
		values (
			$1,
			$2,
			date_trunc('day', now() at time zone 'utc'),
			date_trunc('day', now() at time zone 'utc') + ( case when $3 = true then '1 month' else '1 year' end )::interval - '1 day'::interval,
			$4,
			$5,
			$6,
			$7 )
	`, currUser.UserID, args.User.ChargeMonth, args.User.ChargeMonth, "ch_xxx", 0, "Visa", "4242") // FIXME
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

type ResetPasswordArgs struct {
	Username    string
	Keychain    string
	NewPassword string
}

func (r *RootRx) ResetPassword(ctx context.Context, args ResetPasswordArgs) (*UserRx, error) {
	currUser := CurrentSessionFromContext(ctx)
	if !currUser.IsUnauth() {
		return nil, ErrUserMustBeUnauth
	}
	tx, err := DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	var userID string
	err = tx.QueryRow(`
		select user_id
		from users
		where username = $1
	`, args.Username).Scan(&userID)
	if err != nil {
		return nil, err
	}
	var matches bool
	err = tx.QueryRow(`
		select users.passcode || '-' || billing.stripe_card_last_four = $1
		from users
		join billing
		on users.user_id = billing.user_id
		where users.user_id = $2
	`, args.Keychain, userID).Scan(&matches)
	if err != nil {
		return nil, err
	}
	// NOTE: Returns `sql.ErrNoRows` for simplicity.
	if !matches {
		return nil, sql.ErrNoRows
	}
	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(args.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	_, err = DB.Exec(`
		update users
		set password_hash = $1
		where user_id = $2
	`, newPasswordHash, userID)
	if err != nil {
		return nil, err
	}
	createSessionArgs := struct{ Username, Password string }{args.Username, args.NewPassword}
	return Rx.CreateSession(ctx, createSessionArgs)
}

// func (r *RootRx) ChangePassword(ctx context.Context, args struct{ NewPassword string }) (bool, error) {
// 	currUser := CurrentSessionFromContext(ctx)
// 	if !currUser.IsAuth() {
// 		return false, ErrUserMustBeAuth
// 	}
// 	tx, err := DB.Begin()
// 	if err != nil {
// 		return false, err
// 	}
// 	defer tx.Rollback()
// 	newHash, err := bcrypt.GenerateFromPassword([]byte(args.NewPassword), bcrypt.DefaultCost)
// 	if err != nil {
// 		return false, err
// 	}
// 	_, err = DB.Exec(`
// 		update users
// 		set password_hash = $1
// 		where user_id = $2
// 	`, newHash, currUser.UserID)
// 	if err != nil {
// 		return false, err
// 	}
// 	err = tx.Commit()
// 	if err != nil {
// 		return false, err
// 	}
// 	return true, nil
// }

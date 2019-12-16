package main

import (
	"context"
	"time"
)

type Date struct {
	Year  int32
	Month int32
	Day   int32
}

func (r *RootRx) NextMonth(ctx context.Context) (*DateRx, error) {
	var t time.Time
	err := DB.QueryRow(`
		select now() at time zone 'utc' + '1 month'
	`).Scan(&t)
	if err != nil {
		return nil, err
	}
	date := &Date{int32(t.Year()), int32(t.Month() - 1), int32(t.Day())}
	return &DateRx{date}, nil
}

func (r *RootRx) NextYear(ctx context.Context) (*DateRx, error) {
	var t time.Time
	err := DB.QueryRow(`
		select now() at time zone 'utc' + '1 year'
	`).Scan(&t)
	if err != nil {
		return nil, err
	}
	date := &Date{int32(t.Year()), int32(t.Month() - 1), int32(t.Day())}
	return &DateRx{date}, nil
}

/*
 * DateRx
 */

type DateRx struct{ date *Date }

func (r *DateRx) Year() int32 {
	return r.date.Year
}

func (r *DateRx) Month() int32 {
	return r.date.Month
}

func (r *DateRx) Day() int32 {
	return r.date.Day
}

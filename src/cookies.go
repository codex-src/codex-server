package main

import (
	"math"
	"net/http"
	"time"
)

const SESSION_KEY = "codex-session-id"

func maxAge(future time.Time) int {
	dur := time.Until(future)
	seconds := int(math.Round(dur.Seconds()))
	return seconds
}

func SetCookie(w http.ResponseWriter, key, value string, expiry time.Time) {
	cookie := &http.Cookie{
		Name:     key,
		Value:    value,
		Path:     "/",
		MaxAge:   maxAge(expiry),
		Secure:   false, //FIXME
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
}

func ClearCookie(w http.ResponseWriter, key, value string) {
	cookie := &http.Cookie{
		Name:     key,
		Value:    value,
		Path:     "/",
		MaxAge:   -1,
		Secure:   false, //FIXME
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
}

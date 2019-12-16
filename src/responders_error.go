package main

// import (
// 	"database/sql"
// 	"errors"
// 	"net/http"
//
// 	"github.com/lib/pq"
// 	"golang.org/x/crypto/bcrypt"
// )
//
// // `computeStatusCode` computes a status code for an error.
// func computeStatusCode(errPtr *error) int {
// 	// Authentication error:
// 	switch *errPtr {
// 	case ErrMustBeUnauth:
// 		fallthrough
// 	case ErrMustBeAuth:
// 		fallthrough
// 	case bcrypt.ErrMismatchedHashAndPassword:
// 		*errPtr = nil
// 		return 401
// 	}
// 	// SQL error:
// 	if errors.Is(*errPtr, sql.ErrNoRows) {
// 		*errPtr = nil
// 		return 404
// 	}
// 	// Postgres error:
// 	pqErr := &pq.Error{}
// 	if errors.As(*errPtr, &pqErr) {
// 		// Integrity constraint violation:
// 		if pqErr.Code.Class() == "23" {
// 			*errPtr = nil
// 			return 402
// 		}
// 	}
// 	// Unexpected error:
// 	return 500
// }
//
// // `RespondError` takes a pointer to an error and returns a
// // responder function based on the computed status code.
// func RespondError(errPtr *error) func(http.ResponseWriter) {
// 	if *errPtr == nil {
// 		return RespondBadRequest
// 	}
// 	statusCode := computeStatusCode(errPtr)
// 	return NewResponder(statusCode)
// }

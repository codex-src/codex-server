package main

// import (
// 	"fmt"
// 	"net/http"
// )
// 
// // stripe.com/docs/api/errors
// //
// // 200: Everything worked as expected.
// // 400: The request was unacceptable, often due to missing a required parameter.
// // 401: No valid API key provided.
// // 402: The parameters were valid but the request failed.
// // 404: The requested resource doesn't exist.
// // 409: The request conflicts with another request.
// // 429: Too many requests hit the API too quickly.
// // 500: Something went wrong on our end.
// 
// const (
// 	StatusCodeOK              = 200
// 	StatusCodeBadRequest      = 400
// 	StatusCodeUnauthorized    = 401
// 	StatusCodeRequestFailed   = 402
// 	StatusCodeNotFound        = 404
// 	StatusCodeConflict        = 409
// 	StatusCodeTooManyRequests = 429
// 	StatusCodeServerError     = 500
// )
// 
// var Statuses = map[int]string{
// 	StatusCodeOK:              "OK",
// 	StatusCodeBadRequest:      "Bad Request",
// 	StatusCodeUnauthorized:    "Unauthorized",
// 	StatusCodeRequestFailed:   "Request Failed",
// 	StatusCodeNotFound:        "Not Found",
// 	StatusCodeConflict:        "Conflict",
// 	StatusCodeTooManyRequests: "Too Many Requests",
// 	StatusCodeServerError:     "Server Error",
// }
// 
// var (
// 	RespondOK              = NewResponder(StatusCodeOK)
// 	RespondBadRequest      = NewResponder(StatusCodeBadRequest)
// 	RespondUnauthorized    = NewResponder(StatusCodeUnauthorized)
// 	RespondRequestFailed   = NewResponder(StatusCodeRequestFailed)
// 	RespondNotFound        = NewResponder(StatusCodeNotFound)
// 	RespondConflict        = NewResponder(StatusCodeConflict)
// 	RespondTooManyRequests = NewResponder(StatusCodeTooManyRequests)
// 	RespondServerError     = NewResponder(StatusCodeServerError)
// )
// 
// func NewResponder(statusCode int) func(http.ResponseWriter) {
// 	respond := func(w http.ResponseWriter) {
// 		if statusCode >= 200 && statusCode <= 299 {
// 			w.WriteHeader(statusCode)
// 			return
// 		}
// 		status, ok := Statuses[statusCode]
// 		if !ok {
// 			err := fmt.Errorf("no such status code %d", statusCode)
// 			panic(err)
// 		}
// 		errStr := fmt.Sprintf("%d %s", statusCode, status)
// 		http.Error(w, errStr, statusCode)
// 	}
// 	return respond
// }

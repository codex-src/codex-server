package main

import "net/http"

func writeCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Content-Type, X-Requested-With")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
}

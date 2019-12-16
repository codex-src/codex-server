package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	graphql "github.com/graph-gophers/graphql-go"
	_ "github.com/lib/pq"
)

// Postgres database:
var DB *sql.DB

// GraphQL schema:
var Schema *graphql.Schema

type Query struct {
	Query     string
	Variables map[string]interface{}
}

func handleGraphQL(w http.ResponseWriter, r *http.Request) {
	// Enable cross-origin resource sharing:
	writeCORSHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	// Extend the current session if authenticated:
	curr, err := ExtendCurrentSession(w, r)
	if err != nil {
		http.Error(w, "500 Server Error", http.StatusInternalServerError)
		check(err, "ExtendCurrentSession")
		return
	}
	// Create a context with the current session as a value:
	ctx := WithCurrentSession(context.Background(), curr)
	// Unmarshal query and variables:
	var query Query
	err = json.NewDecoder(r.Body).Decode(&query)
	if err != nil {
		http.Error(w, "500 Server Error", http.StatusInternalServerError)
		check(err, "json.NewDecoder")
		return
	}
	// Execute query and marshal response and errors:
	res := Schema.Exec(ctx, query.Query, "", query.Variables)
	b, err := json.MarshalIndent(res, "", "\t")
	if err != nil {
		http.Error(w, "500 Server Error", http.StatusInternalServerError)
		check(err, "json.MarshalIndent")
		return
	}
	// Write response:
	fmt.Fprintln(w, string(b))
}

func main() {
	// Connect to the database:
	var err error
	DB, err = sql.Open("postgres", "postgres://zaydek@localhost/codex?sslmode=disable")
	must(err, "sql.Open")
	err = DB.Ping()
	must(err, "DB.Ping")
	defer DB.Close()
	// Parse the schema:
	b, err := ioutil.ReadFile("schema.graphql")
	must(err, "ioutil.ReadFile")
	Schema, err = graphql.ParseSchema(string(b), &RootRx{})
	must(err, "graphql.ParseSchema")
	// Listen and serve:
	http.HandleFunc("/graphql", handleGraphQL)
	// http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
	// 	w.WriteHeader(http.StatusOK)
	// 	return
	// })
	err = http.ListenAndServe(":8000", nil)
	must(err, "http.ListenAndServe")
}

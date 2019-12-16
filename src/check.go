package main

import (
	"fmt"
	"log"
)

// `check` logs on a non-nil error with a description.
func check(err error, desc string) {
	if err == nil {
		return
	}
	wrapped := fmt.Errorf("%s: %w", desc, err)
	log.Println(wrapped)
}

// `must` logs fatally on a non-nil error with a
// description.
func must(err error, desc string) {
	if err == nil {
		return
	}
	wrapped := fmt.Errorf("%s: %w", desc, err)
	log.Fatalln(wrapped)
}

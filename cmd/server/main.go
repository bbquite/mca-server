package main

import (
	"net/http"
)

func main() {
	// http.HandleFunc(`/api`, apiPage)
	// http.HandleFunc(`/`, mainPage)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

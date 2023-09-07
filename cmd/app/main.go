package main

import (
	"log"
	"net/http"

	"github.com/brunoquindeler/go-application/internal"
)

func main() {
	// store := internal.NewInMemoryPlayerStore()
	store, err := internal.NewSQLitePlayerStore()
	if err != nil {
		panic(err)
	}

	playerServer := internal.NewPlayerServer(store)

	log.Fatal(http.ListenAndServe(":5000", playerServer))
}

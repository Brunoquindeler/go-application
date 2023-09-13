package main

import (
	"log"
	"net/http"

	"github.com/brunoquindeler/go-application/internal"
)

func main() {
	// store := internal.NewInMemoryPlayerStore()

	sqliteConn, err := internal.GetSQLiteConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer sqliteConn.Close()

	store, err := internal.NewSQLitePlayerStore(sqliteConn)
	if err != nil {
		log.Fatal(err)
	}

	playerServer := internal.NewPlayerServer(store)

	log.Fatal(http.ListenAndServe(":5000", playerServer))
}

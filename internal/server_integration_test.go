package internal

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	os.Remove(PlayerDB)

	sqliteConn, err := GetSQLiteConnection()
	if err != nil {
		t.Fatal(err)
	}
	defer sqliteConn.Close()

	sqliteStore, err := NewSQLitePlayerStore(sqliteConn)
	if err != nil {
		t.Fatal(err)
	}

	inMemoryStore := NewInMemoryPlayerStore()

	tests := []struct {
		desc  string
		store playerStore
	}{
		{
			desc:  "In Memory Store",
			store: inMemoryStore,
		},
		{
			desc:  "SQLite Store",
			store: sqliteStore,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			server := NewPlayerServer(test.store)
			player := "Pepper"

			server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
			server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
			server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

			t.Run("get score", func(t *testing.T) {
				response := httptest.NewRecorder()
				server.ServeHTTP(response, newGetScoreRequest(player))
				assertStatus(t, response.Code, http.StatusOK)
				assertResponseBody(t, response.Body.String(), "3")
			})

			t.Run("get league", func(t *testing.T) {
				response := httptest.NewRecorder()
				server.ServeHTTP(response, newLeagueRequest())
				assertStatus(t, response.Code, http.StatusOK)

				got := getLeagueFromResponse(t, response.Body)
				want := []Player{
					{"Pepper", 3},
				}

				assertLeague(t, got, want)
			})
		})
	}
}

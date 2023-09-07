package internal

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestRecordingWinsAndRetrievingThemWithInMemory(t *testing.T) {
	store := NewInMemoryPlayerStore()
	server := PlayerServer{Store: store}
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

	response := httptest.NewRecorder()
	server.ServeHTTP(response, newGetScoreRequest(player))

	assertStatus(t, response.Code, http.StatusOK)
	assertResponseBody(t, response.Body.String(), "3")
}

func TestRecordingWinsAndRetrievingThemWithSQLite(t *testing.T) {
	os.Remove(playerDB)

	store, _ := NewSQLitePlayerStore()
	server := PlayerServer{Store: store}
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

	response := httptest.NewRecorder()
	server.ServeHTTP(response, newGetScoreRequest(player))

	assertStatus(t, response.Code, http.StatusOK)
	assertResponseBody(t, response.Body.String(), "3")
}

package internal

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const create = `
CREATE TABLE IF NOT EXISTS players (
    id   INTEGER    PRIMARY KEY AUTOINCREMENT UNIQUE NOT NULL,
    name TEXT (255) NOT NULL UNIQUE,
    wins INTEGER    DEFAULT (0) 
);
`

type SQLitePlayerStore struct {
	db *sql.DB
}

func NewSQLitePlayerStore() (*SQLitePlayerStore, error) {
	db, err := sql.Open("sqlite3", playerDB)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(create); err != nil {
		return nil, err
	}

	return &SQLitePlayerStore{
		db: db,
	}, nil
}

func (s *SQLitePlayerStore) RecordWin(name string) {
	var err error
	if wins := s.GetPlayerScore(name); wins == 0 {
		_, err = s.db.Exec("INSERT INTO players (id, name, wins) VALUES (NULL,?,?);", name, 1)
		if err != nil {
			log.Println(err.Error())
		}
	} else {
		_, err = s.db.Exec("UPDATE players SET wins = ? WHERE name=?;", wins+1, name)
		if err != nil {
			log.Println(err.Error())
		}
	}
}

func (s *SQLitePlayerStore) GetPlayerScore(name string) int {
	row := s.db.QueryRow("SELECT wins FROM players WHERE name=?", name)

	var wins int
	if err := row.Scan(&wins); err != nil && err != sql.ErrNoRows {
		log.Println(err.Error())
	}

	return wins
}

func (s *SQLitePlayerStore) GetLeague() []Player {
	rows, err := s.db.Query("SELECT name, wins FROM players")
	if err != nil {
		log.Println(err.Error())
	}
	defer rows.Close()

	var players []Player
	for rows.Next() {
		player := Player{}
		if err := rows.Scan(&player.Name, &player.Wins); err != nil && err != sql.ErrNoRows {
			log.Println(err.Error())
		}
		players = append(players, player)
	}

	return players
}

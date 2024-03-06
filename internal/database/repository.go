package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Repository struct {
	db       *sql.DB
	PingData *PingData
}

func ConnectDB() (*sql.DB, error) {
	return sql.Open("sqlite3", "API.db")
}

func New(db *sql.DB) *Repository {
	return &Repository{
		db:       db,
		PingData: initPingData(db),
	}
}

func (r *Repository) CreateTables() error {
	if err := r.PingData.createTable(); err != nil {
		return err
	}
	return nil
}

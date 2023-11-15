package pgdb

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type PostgresDB struct {
	DB *sql.DB
}

func New(db *sql.DB) *PostgresDB {
	return &PostgresDB{
		DB: db,
	}
}

func NewDB(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return &sql.DB{}, nil
	}

	if err := db.Ping(); err != nil {
		return &sql.DB{}, err
	}
	return db, nil
}

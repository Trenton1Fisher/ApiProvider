package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func NewPostgreSQLClient(dbURL string) (*sql.DB, error) {
    db, err := sql.Open("postgres", dbURL)
    if err != nil {
        return nil, err
    }

    err = db.Ping()
    if err != nil {
        return nil, err
    }

    return db, err
}
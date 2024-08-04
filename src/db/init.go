package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var db *Queries

// DbInit is so when we share one sql.Open for various
// db calls throughout the program. Also, it means we
// don't have to handle an err on said subsequent db calls
func DbInit() error {
	var d, err = sql.Open("sqlite3", "file:database.db")
	if err != nil {
		return err
	}
	db = New(d)
	return nil
}

func Db() *Queries {
	return db
}

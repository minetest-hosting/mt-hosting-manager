package db

import (
	"database/sql"
	"path"

	_ "github.com/mattn/go-sqlite3"
)

func Init(data_dir string) (*sql.DB, error) {
	var err error
	db, err := sql.Open("sqlite3", path.Join(data_dir, "mt-hosting.sqlite?_timeout=5000&_journal=WAL&_foreign_keys=true"))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

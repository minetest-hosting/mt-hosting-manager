package db

import (
	"database/sql"
	"errors"
	"path"

	_ "modernc.org/sqlite"
)

func Init(data_dir string) (*sql.DB, error) {
	var err error
	db, err := sql.Open("sqlite", path.Join(data_dir, "mt-hosting.sqlite"))
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("pragma journal_mode = wal")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("pragma foreign_keys = ON")
	if err != nil {
		return nil, err
	}

	return db, nil
}

func EnableWAL(db *sql.DB) error {
	result := db.QueryRow("pragma journal_mode;")
	var mode string
	err := result.Scan(&mode)
	if err != nil {
		return err
	}

	if mode != "wal" {
		_, err = db.Exec("pragma journal_mode = wal;")
		if err != nil {
			return errors.New("couldn't switch the db-journal to wal-mode: " + err.Error())
		}
	}

	return nil
}

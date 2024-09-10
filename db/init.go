package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init() (*sql.DB, *gorm.DB, error) {
	connStr := fmt.Sprintf(
		"user=%s password=%s port=%s host=%s dbname=%s sslmode=disable",
		os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"),
		os.Getenv("PGPORT"),
		os.Getenv("PGHOST"),
		os.Getenv("PGDATABASE"))

	logrus.Infof("Connecting to %s", connStr)
	var err error
	DB, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, nil, fmt.Errorf("sql.Open error: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		return nil, nil, fmt.Errorf("DB.Ping error: %v", err)
	}

	g, err := gorm.Open(postgres.Open(connStr))
	if err != nil {
		return nil, nil, fmt.Errorf("gorm.Open error: %v", err)
	}

	return DB, g, nil
}

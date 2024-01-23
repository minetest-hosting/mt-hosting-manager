package db_test

import (
	"database/sql"
	"mt-hosting-manager/db"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func SetupDB(t *testing.T) *sql.DB {
	db_, err := db.Init()
	assert.NoError(t, err)
	assert.NotNil(t, db_)

	err = db.Migrate(db_)
	assert.NoError(t, err)

	return db_
}

func SetupRepos(t *testing.T) *db.Repositories {
	if os.Getenv("PGHOST") == "" {
		t.SkipNow()
	}

	repos := db.NewRepositories(SetupDB(t))
	assert.NoError(t, repos.UserRepo.DeleteAll())
	assert.NoError(t, repos.NodeTypeRepo.DeleteAll())
	assert.NoError(t, repos.ExchangeRateRepo.DeleteAll())
	return repos
}

func TestMigrate(t *testing.T) {
	SetupDB(t)
}

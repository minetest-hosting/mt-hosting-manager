package db_test

import (
	"database/sql"
	"mt-hosting-manager/db"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupDB(t *testing.T) *sql.DB {
	tmpdir, err := os.MkdirTemp(os.TempDir(), "mt-hosting")
	assert.NoError(t, err)

	db_, err := db.Init(tmpdir)
	assert.NoError(t, err)
	assert.NotNil(t, db_)

	err = db.Migrate(db_)
	assert.NoError(t, err)

	return db_
}

func setupRepos(t *testing.T) *db.Repositories {
	return db.NewRepositories(setupDB(t))
}

func TestMigrate(t *testing.T) {
	setupDB(t)
}

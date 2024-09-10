package core_test

import (
	"database/sql"
	"mt-hosting-manager/db"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func SetupDB(t *testing.T) (*sql.DB, *gorm.DB) {
	db_, g, err := db.Init()
	assert.NoError(t, err)
	assert.NotNil(t, db_)
	assert.NotNil(t, g)

	err = db.Migrate(db_)
	assert.NoError(t, err)

	return db_, g
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

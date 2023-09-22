package core_test

import (
	"mt-hosting-manager/db"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func SetupRepos(t *testing.T) *db.Repositories {
	tmpdir, err := os.MkdirTemp(os.TempDir(), "mt-hosting")
	assert.NoError(t, err)

	db_, err := db.Init(tmpdir)
	assert.NoError(t, err)
	assert.NotNil(t, db_)

	err = db.Migrate(db_)
	assert.NoError(t, err)

	return db.NewRepositories(db_)
}

func TestCollect(t *testing.T) {
	//TODO
}

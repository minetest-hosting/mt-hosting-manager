package db_test

import (
	"mt-hosting-manager/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNodeRepository(t *testing.T) {
	repos := setupRepos(t)

	assert.NoError(t, repos.NodeRepo.Insert(&types.Node{
		Provider:   types.ProviderHetzner,
		ServerType: "cx11",
		Name:       "",
	}))

	//TODO: get / update / delete
}

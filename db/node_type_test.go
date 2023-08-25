package db_test

import (
	"mt-hosting-manager/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNodeRepository(t *testing.T) {
	repos := SetupRepos(t)

	assert.NoError(t, repos.NodeTypeRepo.Insert(&types.NodeType{
		Provider:   types.ProviderHetzner,
		ServerType: "cx11",
		Name:       "",
	}))

	//TODO: get / update / delete
}

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

	list, err := repos.NodeTypeRepo.GetAll()
	assert.NoError(t, err)
	assert.True(t, len(list) > 0)

	nt := list[0]
	nt.CpuCount = 666
	assert.NoError(t, repos.NodeTypeRepo.Update(nt))

	assert.NoError(t, repos.NodeTypeRepo.Delete(nt.ID))
}

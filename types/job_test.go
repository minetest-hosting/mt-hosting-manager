package types_test

import (
	"mt-hosting-manager/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestJobDataStruct struct {
	Var int `json:"var"`
}

func TestJobData(t *testing.T) {
	d := &TestJobDataStruct{
		Var: 666,
	}

	// job with data

	j := &types.Job{}
	assert.False(t, j.HasData())
	assert.NoError(t, j.SetData(d))
	assert.True(t, j.HasData())

	d2 := &TestJobDataStruct{}
	assert.NoError(t, j.GetData(d2))
	assert.Equal(t, d.Var, d2.Var)

	// job without data

	j = &types.Job{}
	assert.False(t, j.HasData())

}

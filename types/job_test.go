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

	j := &types.Job{}
	assert.NoError(t, j.SetData(d))

	d2 := &TestJobDataStruct{}
	assert.NoError(t, j.GetData(d2))
	assert.Equal(t, d.Var, d2.Var)
}

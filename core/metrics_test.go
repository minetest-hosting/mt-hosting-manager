package core_test

import (
	"mt-hosting-manager/core"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNodeExporterMetrics(t *testing.T) {
	data, err := os.ReadFile("testdata/metrics.txt")
	assert.NoError(t, err)
	assert.NotNil(t, data)

	metrics, err := core.ParseNodeExporterMetrics(data)
	assert.NoError(t, err)
	assert.NotNil(t, metrics)

	assert.Equal(t, 16, metrics.LoadPercent)
	assert.Equal(t, int64(50531442688), metrics.DiskSize)
	assert.Equal(t, int64(14798311424), metrics.DiskUsed)
	assert.Equal(t, int64(16791117824), metrics.MemorySize)
	assert.Equal(t, int64(12169216000), metrics.MemoryUsed)
}

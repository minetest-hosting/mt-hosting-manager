package mtui_test

import (
	"mt-hosting-manager/api/mtui"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoginDownload(t *testing.T) {
	t.SkipNow() // disabled for automated tests

	url := "https://pandorabox.io/ui"
	//jwt_key := "REDACTED"
	//username := "admin"

	c := mtui.New(url)
	//err := c.Login(username, jwt_key)
	//assert.NoError(t, err)

	stats, err := c.GetStats()
	assert.NoError(t, err)
	assert.NotNil(t, stats)

	/*
		r, err := c.DownloadZip("/")
		assert.NoError(t, err)
		defer r.Close()

		f, err := os.CreateTemp(os.TempDir(), "mtui-download.zip")
		assert.NoError(t, err)
		defer f.Close()

		n, err := io.Copy(f, r)
		assert.NoError(t, err)
		fmt.Printf("copied %d bytes to %s\n", n, f.Name())
	*/
}

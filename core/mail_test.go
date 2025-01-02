package core_test

import (
	"mt-hosting-manager/core"
	"mt-hosting-manager/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendMail(t *testing.T) {

	to := ""

	cfg := types.NewConfig()
	cfg.MailHost = ""
	cfg.MailAddress = ""
	cfg.MailPassword = ""

	if cfg.MailAddress == "" || to == "" {
		t.SkipNow()
	}

	c := core.New(nil, cfg)

	assert.NoError(t, c.SendMail(to, "Test mail", "Test body"))
}

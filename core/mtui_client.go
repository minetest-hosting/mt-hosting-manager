package core

import (
	"fmt"
	"mt-hosting-manager/api/mtui"
	"mt-hosting-manager/types"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
)

var mtui_cache = expirable.NewLRU[string, *mtui.MtuiClient](20, nil, time.Hour*2)

func (c *Core) GetMTUIClient(server *types.MinetestServer) (*mtui.MtuiClient, error) {
	client, found := mtui_cache.Get(server.ID)
	if !found {
		// create a new client and log in
		url := fmt.Sprintf("https://%s.%s/ui", server.DNSName, c.cfg.HostingDomainSuffix)
		client = mtui.New(url)
		err := client.Login(server.Admin, server.JWTKey)
		if err != nil {
			return nil, fmt.Errorf("login error: %v", err)
		}
		mtui_cache.Add(server.ID, client)
	}

	return client, nil
}

package wallee

import (
	"net/http"
)

type WalleeClient struct {
	UserID  string
	Key     string
	SpaceID string
	client  http.Client
}

func New(userid, spaceid, key string) *WalleeClient {
	return &WalleeClient{
		UserID:  userid,
		SpaceID: spaceid,
		Key:     key,
		client:  http.Client{},
	}
}

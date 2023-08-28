package wallee

import (
	"net/http"
	"os"
)

type WalleeClient struct {
	UserID  string
	Key     string
	SpaceID string
	client  http.Client
}

func New() *WalleeClient {
	return &WalleeClient{
		UserID:  os.Getenv("WALLEE_USERID"),
		SpaceID: os.Getenv("WALLEE_SPACEID"),
		Key:     os.Getenv("WALLEE_KEY"),
		client:  http.Client{},
	}
}

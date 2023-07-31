package wallee

import "net/http"

type WalleeClient struct {
	UserID  string
	Key     string
	SpaceID string
	client  http.Client
}

func New(userID, spaceID, key string) *WalleeClient {
	return &WalleeClient{
		UserID:  userID,
		SpaceID: spaceID,
		Key:     key,
		client:  http.Client{},
	}
}

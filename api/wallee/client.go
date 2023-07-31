package wallee

type WalleeClient struct {
	UserID  string
	Key     string
	SpaceID string
}

func New(userID, spaceID, key string) *WalleeClient {
	return &WalleeClient{
		UserID:  userID,
		SpaceID: spaceID,
		Key:     key,
	}
}

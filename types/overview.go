package types

type MinetestServerOverview struct {
	Job *Job `json:"job"`
	*MinetestServer
}

type OverviewData struct {
	*UserNode
	Job     *Job                      `json:"job"`
	Servers []*MinetestServerOverview `json:"servers"`
}

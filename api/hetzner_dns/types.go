package hetzner_dns

type RecordType string

const (
	RecordA     RecordType = "A"
	RecordAAAA  RecordType = "AAAA"
	RecordCNAME RecordType = "CNAME"
)

type Record struct {
	Type     RecordType `json:"type"`
	ID       string     `json:"id"`
	Created  string     `json:"created"`
	Modified string     `json:"modified"`
	ZoneID   string     `json:"zone_id"`
	Name     string     `json:"name"`
	Value    string     `json:"value"`
	TTL      int        `json:"ttl"`
}

type RecordsResponse struct {
	Records []*Record `json:"records"`
}

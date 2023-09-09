package types

func MailQueueProvider() *MailQueue { return &MailQueue{} }

type MailQueueState string

const (
	MailQueueStateCreated     MailQueueState = "CREATED"
	MailQueueStatePending     MailQueueState = "PENDING"
	MailQueueStateDoneSuccess MailQueueState = "DONE_SUCCESS"
	MailQueueStateDoneFailure MailQueueState = "DONE_FAILURE"
)

type MailQueue struct {
	ID        string         `json:"id"`
	State     MailQueueState `json:"state"`
	Timestamp int64          `json:"timestamp"`
	Receiver  string         `json:"receiver"`
	Subject   string         `json:"subject"`
	Content   string         `json:"content"`
}

func (m *MailQueue) Columns(action string) []string {
	return []string{
		"id",
		"state",
		"timestamp",
		"receiver",
		"subject",
		"content",
	}
}

func (m *MailQueue) Table() string {
	return "mail_queue"
}

func (m *MailQueue) Scan(action string, r func(dest ...any) error) error {
	return r(&m.ID, &m.State, &m.Timestamp, &m.Receiver, &m.Subject, &m.Content)
}

func (m *MailQueue) Values(action string) []any {
	return []any{m.ID, m.State, m.Timestamp, m.Receiver, m.Subject, m.Content}
}

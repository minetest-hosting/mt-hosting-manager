package worker

import (
	"fmt"
	"mt-hosting-manager/api/hetzner_dns"

	"github.com/sirupsen/logrus"
)

func (w *Worker) UpdateDNSRecord(records *hetzner_dns.RecordsResponse, t hetzner_dns.RecordType, name, value string) error {
	var existing_record *hetzner_dns.Record
	for _, r := range records.Records {
		if r.Type == t && r.Name == name {
			existing_record = r
			break
		}
	}

	if existing_record == nil {
		// create new record
		logrus.WithFields(logrus.Fields{
			"value": value,
			"name":  name,
			"type":  t,
		}).Info("Creating Record")

		new_record := &hetzner_dns.Record{
			Type:  t,
			Name:  name,
			Value: value,
			TTL:   300,
		}
		err := w.hdc.CreateRecord(new_record)
		if err != nil {
			return fmt.Errorf("create record error type=%s, name=%s, value=%s: %v", t, name, value, err)
		}

	} else if existing_record.Value != value {
		// update record
		existing_record.Value = value
		err := w.hdc.UpdateRecord(existing_record)
		if err != nil {
			return fmt.Errorf("update record error type=%s, name=%s, value=%s: %v", t, name, value, err)
		}
	}

	return nil
}

func (w *Worker) RemoveDNSRecord(records *hetzner_dns.RecordsResponse, t hetzner_dns.RecordType, name string) error {
	var existing_record *hetzner_dns.Record

	for _, r := range records.Records {
		if r.Type == t && r.Name == name {
			existing_record = r
			break
		}
	}

	if existing_record != nil {
		err := w.hdc.DeleteRecord(existing_record.ID)
		if err != nil {
			return fmt.Errorf("remove record error type=%s, name=%s: %v", t, name, err)
		}
	}

	return nil
}

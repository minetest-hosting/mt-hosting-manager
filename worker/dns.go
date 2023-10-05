package worker

import (
	"fmt"
	"mt-hosting-manager/api/hetzner_dns"

	"github.com/sirupsen/logrus"
)

func (w *Worker) CreateDNSRecord(t hetzner_dns.RecordType, name, value string) (*hetzner_dns.Record, error) {
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
	record, err := w.hdc.CreateRecord(new_record)
	if err != nil {
		return nil, fmt.Errorf("create record error type=%s, name=%s, value=%s: %v", t, name, value, err)
	}
	return record.Record, nil
}

func (w *Worker) UpdateDNSRecord(record_id string, value string) error {
	record_response, err := w.hdc.GetRecord(record_id)
	if err != nil {
		return err
	}
	record := record_response.Record
	record.Value = value

	logrus.WithFields(logrus.Fields{
		"id":        record_id,
		"new-value": value,
		"name":      record.Name,
		"type":      record.Type,
	}).Info("Updating Record")

	return w.hdc.UpdateRecord(record)
}

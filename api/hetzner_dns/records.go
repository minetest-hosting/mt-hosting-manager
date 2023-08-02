package hetzner_dns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *HetznerDNSClient) GetRecords() (*RecordsResponse, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://dns.hetzner.com/api/v1/records?zone_id=%s", c.zoneID), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Auth-API-Token", c.Key)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected response-code: %d", resp.StatusCode)
	}

	rr := &RecordsResponse{}
	err = json.NewDecoder(resp.Body).Decode(rr)

	return rr, err
}

func (c *HetznerDNSClient) CreateRecord(rec *Record) error {
	rec.ZoneID = c.zoneID

	data, err := json.Marshal(rec)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, "https://dns.hetzner.com/api/v1/records", bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Set("Auth-API-Token", c.Key)

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected response-code: %d", resp.StatusCode)
	}

	return nil
}

func (c *HetznerDNSClient) UpdateRecord(rec *Record) error {
	data, err := json.Marshal(rec)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("https://dns.hetzner.com/api/v1/records/%s", rec.ID), bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Set("Auth-API-Token", c.Key)

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected response-code: %d", resp.StatusCode)
	}

	return nil
}

func (c *HetznerDNSClient) DeleteRecord(recordID string) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("https://dns.hetzner.com/api/v1/records/%s", recordID), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Auth-API-Token", c.Key)

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected response-code: %d", resp.StatusCode)
	}

	return nil
}

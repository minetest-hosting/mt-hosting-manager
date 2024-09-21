package mtui

import (
	"fmt"
	"net/http"
)

func (a *MtuiClient) SetMaintenanceMode(enable bool) error {
	method := http.MethodPut
	if !enable {
		method = http.MethodDelete
	}

	req, err := a.request(method, "api/maintenance", nil)
	if err != nil {
		return fmt.Errorf("request error: %v", err)
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("http do error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("api-response status: %d", resp.StatusCode)
	}

	return nil
}

package mtui

import (
	"fmt"
	"io"
	"net/http"
)

type MtuiClient struct {
	client http.Client
	url    string
	token  string
}

func New(url string) *MtuiClient {
	return &MtuiClient{
		url:    url,
		client: http.Client{},
	}
}

func (a *MtuiClient) Login(username, jwt_key string) error {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/loginadmin/%s", a.url, username), nil)
	if err != nil {
		return fmt.Errorf("create request failed: %v", err)
	}

	q := req.URL.Query()
	q.Set("key", jwt_key)
	q.Set("disable_redirect", "true")
	req.URL.RawQuery = q.Encode()

	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("http do error: %v", err)
	}
	defer resp.Body.Close()

	resp_bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("io.read error: %v", err)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("api-response status: %d", resp.StatusCode)
	}

	if len(resp_bytes) == 0 {
		return fmt.Errorf("token-length is zero")
	}

	a.token = string(resp_bytes)
	return nil
}

func (a *MtuiClient) DownloadRootZip() (io.ReadCloser, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/filebrowser/zip", a.url), nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %v", err)
	}

	q := req.URL.Query()
	q.Set("dir", "/")
	req.URL.RawQuery = q.Encode()

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.token))

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http do error: %v", err)
	}

	return resp.Body, nil
}

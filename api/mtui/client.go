package mtui

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type MtuiClient struct {
	client  http.Client
	url     string
	cookies []*http.Cookie
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

	if resp.StatusCode != 200 {
		return fmt.Errorf("api-response status: %d", resp.StatusCode)
	}

	a.cookies = resp.Cookies()
	return nil
}

func (a *MtuiClient) DownloadZip(dir string) (io.ReadCloser, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/filebrowser/zip", a.url), nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %v", err)
	}

	q := req.URL.Query()
	q.Set("dir", dir)
	req.URL.RawQuery = q.Encode()

	for _, c := range a.cookies {
		req.AddCookie(c)
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http do error: %v", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("api-response status: %d", resp.StatusCode)
	}

	return resp.Body, nil
}

func (a *MtuiClient) GetDirectorySize(dir string) (int64, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/filebrowser/size", a.url), nil)
	if err != nil {
		return 0, fmt.Errorf("create request failed: %v", err)
	}

	q := req.URL.Query()
	q.Set("dir", dir)
	req.URL.RawQuery = q.Encode()

	for _, c := range a.cookies {
		req.AddCookie(c)
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("http do error: %v", err)
	}

	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("api-response status: %d", resp.StatusCode)
	}

	resp_bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("readall error: %v", err)
	}

	size, err := strconv.ParseInt(string(resp_bytes), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parseint error: %v", err)
	}

	return size, nil
}
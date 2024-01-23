package core

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetMTUIDigest() (string, error) {
	image := "minetest-go/mtui"
	tag := "latest"

	resp, err := http.Get(fmt.Sprintf("https://ghcr.io/token?scope=repository:%s:pull", image))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("unexpected response-code: %d", resp.StatusCode)
	}

	tr := map[string]string{}
	err = json.NewDecoder(resp.Body).Decode(&tr)
	if err != nil {
		return "", err
	}
	token := tr["token"]
	if token == "" {
		return "", fmt.Errorf("no token received")
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://ghcr.io/v2/%s/manifests/%s", image, tag), nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("unexpected response-code: %d", resp.StatusCode)
	}

	digest := resp.Header.Get("docker-content-digest")
	if digest == "" {
		return "", fmt.Errorf("no digest header received")
	}

	return digest, nil
}

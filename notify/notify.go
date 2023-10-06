package notify

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	_ "time/tzdata"

	"github.com/sirupsen/logrus"
)

type NtfyNotification struct {
	Topic    string   `json:"topic"`
	Title    string   `json:"title"`
	Message  string   `json:"message"`
	Tags     []string `json:"tags"`
	Click    *string  `json:"click"`
	Priority int      `json:"priority"`
}

var ntfy_url = os.Getenv("NTFY_URL")
var ntfy_topic = os.Getenv("NTFY_TOPIC")
var ntfy_username = os.Getenv("NTFY_USERNAME")
var ntfy_password = os.Getenv("NTFY_PASSWORD")

func Send(nn *NtfyNotification, cache bool) {
	if ntfy_url == "" || ntfy_topic == "" {
		// nothing to do
		return
	}

	nn.Topic = ntfy_topic

	body, err := json.Marshal(nn)
	if err != nil {
		logrus.WithError(err).Error("failed to marshal notification")
		return
	}

	req, err := http.NewRequest(http.MethodPost, ntfy_url, bytes.NewReader(body))
	if err != nil {
		logrus.WithError(err).Error("failed to create request for notification")
		return
	}

	if !cache {
		req.Header.Add("Cache", "no")
	}

	if ntfy_username != "" {
		req.SetBasicAuth(ntfy_username, ntfy_password)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.WithError(err).Error("failed to post notification")
		return
	}
	defer resp.Body.Close()
}

package infrastructure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/env"
)

func PhoneServiceURL() string {
	return env.Get("PHONE_SERVICE_URL", "http://localhost:9090/notify")
}

type phoneClient struct {
	url        string
	httpClient *http.Client
}

func NewPhoneClient(url string) *phoneClient {
	return &phoneClient{
		url:        url,
		httpClient: &http.Client{},
	}
}

type notificationPayload struct {
	Token        string       `json:"token"`
	Notification notification `json:"notification"`
}

type notification struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func (c *phoneClient) Notify(ctx context.Context, deviceToken string) error {
	payload := notificationPayload{
		Token: deviceToken,
		Notification: notification{
			Title: "Daily check-in",
			Body:  "How are you today?",
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("phone service returned status %d", resp.StatusCode)
	}

	return nil
}

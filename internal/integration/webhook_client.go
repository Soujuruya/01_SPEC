package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type WebhookClient struct {
	url    string
	client *http.Client
}

func NewWebhookClient(url string, timeout time.Duration) *WebhookClient {
	return &WebhookClient{
		url: url,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *WebhookClient) Send(ctx context.Context, payload WebhookPayload) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	if resp.StatusCode >= 300 {
		return http.ErrHandlerTimeout // временно
	}

	return nil
}

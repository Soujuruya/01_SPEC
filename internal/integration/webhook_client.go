package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Soujuruya/01_SPEC/internal/pkg/logger"
)

type WebhookClient struct {
	url    string
	client *http.Client
	lg     *logger.Logger
}

func NewWebhookClient(url string, timeout time.Duration, lg *logger.Logger) *WebhookClient {
	return &WebhookClient{
		url: url,
		client: &http.Client{
			Timeout: timeout,
		},
		lg: lg,
	}
}

func (c *WebhookClient) Send(ctx context.Context, payload WebhookPayload, lg *logger.Logger) error {
	body, err := json.Marshal(payload)
	if err != nil {
		lg.Error("WebhookClient: failed to marshal payload", "error", err, "payload", payload)
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewReader(body))
	if err != nil {
		lg.Error("WebhookClient: failed to create request", "error", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		lg.Error("WebhookClient: request failed", "error", err, "url", c.url, "payload", payload)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		lg.Warn("WebhookClient: non-2xx response", "status", resp.StatusCode, "url", c.url, "payload", payload)
		return http.ErrHandlerTimeout
	}

	lg.Info("WebhookClient: webhook sent successfully", "url", c.url, "payload", payload)
	return nil
}

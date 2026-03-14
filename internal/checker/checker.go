package checker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type healthResponse struct {
	Status string `json:"status"`
	DB     string `json:"db"`
}

type Checker struct {
	url    string
	client *http.Client
}

func New(url string, timeout time.Duration) *Checker {
	return &Checker{
		url:    url,
		client: &http.Client{Timeout: timeout},
	}
}

func (c *Checker) Check() error {
	resp, err := c.client.Get(c.url)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var health healthResponse
	if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if health.Status != "ok" {
		return fmt.Errorf("unhealthy: status=%s, db=%s", health.Status, health.DB)
	}

	return nil
}

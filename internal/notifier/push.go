package notifier

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Pusher struct {
	url    string
	client *http.Client
}

func NewPusher(token string) *Pusher {
	return &Pusher{
		url:    fmt.Sprintf("https://torchi.app/api/v1/push/%s", token),
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (p *Pusher) Send(message string) error {
	resp, err := p.client.Post(p.url, "text/plain", strings.NewReader(message))
	if err != nil {
		return fmt.Errorf("push request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("push failed: status %d", resp.StatusCode)
	}

	return nil
}

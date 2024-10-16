package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shariqali-dev/toll-calculator/internal/types"
)

type Client struct {
	Endpoint string
}

func NewClient(endpoint string) *Client {
	return &Client{
		Endpoint: endpoint,
	}
}

func (c *Client) AggregateInvoice(distance types.Distance) error {
	b, err := json.Marshal(distance)
	if err != nil {
		return fmt.Errorf("error marshalling distance: %v", distance)
	}
	req, err := http.NewRequest(http.MethodPost, c.Endpoint, bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("error creating request from bytes: %v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(" service responded with non 200 status code %d", resp.StatusCode)
	}

	return nil
}

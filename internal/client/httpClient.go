package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shariqali-dev/toll-calculator/internal/types"
)

type HTTPClient struct {
	Endpoint string
}

func NewHTTPClient(endpoint string) *HTTPClient {
	return &HTTPClient{
		Endpoint: endpoint,
	}
}

func (c *HTTPClient) AggregateInvoice(distance types.Distance) error {
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

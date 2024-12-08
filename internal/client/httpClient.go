package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shariqali-dev/toll-calculator/internal/types"
	"github.com/sirupsen/logrus"
)

type HTTPClient struct {
	Endpoint string
}

func NewHTTPClient(endpoint string) *HTTPClient {
	return &HTTPClient{
		Endpoint: endpoint,
	}
}

func (c *HTTPClient) Aggregate(ctx context.Context, aggReq *types.AggregateRequest) error {
	b, err := json.Marshal(aggReq)
	if err != nil {
		return fmt.Errorf(aggReq.String())
	}
	req, err := http.NewRequest(http.MethodPost, c.Endpoint+"/aggregate", bytes.NewReader(b))
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

func (c *HTTPClient) GetInvoice(ctx context.Context, id int) (*types.Invoice, error) {
	invRequest := types.GetInvoiceRequest{
		OBUID: int32(id),
	}
	b, err := json.Marshal(invRequest)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("%s/%s?obuID=%d", c.Endpoint, "invoice", invRequest.OBUID)
	logrus.Infof("requesting get invoice -> %s", endpoint)
	req, err := http.NewRequest(http.MethodGet, endpoint, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("the service responded with non 200 status code %d", resp.StatusCode)
	}

	var invoice types.Invoice
	if err := json.NewDecoder(resp.Body).Decode(&invoice); err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return &invoice, nil
}

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/manish-neemnarayan/toll-calculator/types"
)

type HttpClient struct {
	Endpoint string
}

func NewHttpClient(endpoint string) Client {
	return &HttpClient{
		Endpoint: endpoint,
	}
}

func (c *HttpClient) GetInvoice(ctx context.Context, id int) (*types.Invoice, error) {
	invReq := &types.GetInvoiceRequest{
		ObuID: int32(id),
	}

	b, err := json.Marshal(&invReq)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("%s/%s?obuId=%d", c.Endpoint, "invoice", invReq.ObuID)

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("invoice err: %v\n endpoint: %v\n", err, c.Endpoint)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("the aggregate service returned Not Ok error: %d", resp.StatusCode)
	}

	var inv *types.Invoice
	if err := json.NewDecoder(resp.Body).Decode(&inv); err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return inv, nil
}

func (c *HttpClient) Aggregate(ctx context.Context, distance *types.AggregateRequest) error {
	b, err := json.Marshal(distance)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", c.Endpoint+"/aggregate", bytes.NewReader(b))
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("the aggregate service returned Not Ok error: %d", resp.StatusCode)
	}
	resp.Body.Close()
	return nil
}

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

func (c *HttpClient) Aggregate(ctx context.Context, distance *types.AggregateRequest) error {
	b, err := json.Marshal(distance)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", c.Endpoint, bytes.NewReader(b))
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

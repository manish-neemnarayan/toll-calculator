package client

import (
	"context"
	"fmt"

	"github.com/manish-neemnarayan/toll-calculator/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	Endpoint string
	types.AggregatorClient
}

func NewGRPCClient(endpoint string) (Client, error) {
	fmt.Println(endpoint)
	conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("error coming from new-grpc-client")
		return nil, err
	}

	c := types.NewAggregatorClient(conn)
	return &GRPCClient{
		Endpoint:         endpoint,
		AggregatorClient: c,
	}, nil
}

func (c *GRPCClient) Aggregate(ctx context.Context, req *types.AggregateRequest) error {
	fmt.Println("came here hey man----")
	_, err := c.AggregatorClient.Aggregate(ctx, req)
	return err
}

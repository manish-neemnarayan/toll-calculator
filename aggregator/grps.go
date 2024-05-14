package main

import (
	"context"

	"github.com/manish-neemnarayan/toll-calculator/types"
)

type GRPCServer struct {
	types.UnimplementedAggregatorServer
	svc Aggregator
}

func NewGRPCServer(svc Aggregator) *GRPCServer {
	return &GRPCServer{
		svc: svc,
	}
}

func (s *GRPCServer) Aggregate(ctx context.Context, req *types.AggregateRequest) (*types.None, error) {
	distance := types.Distance{
		OBUID: int(req.OBUID),
		Value: req.Value,
		Unix:  req.Unix,
	}
	return &types.None{}, s.svc.AggregateDistance(distance)
}

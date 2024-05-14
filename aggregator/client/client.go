package client

import (
	"context"

	"github.com/manish-neemnarayan/toll-calculator/types"
)

type Client interface {
	Aggregate(context.Context, *types.AggregateRequest) error
}

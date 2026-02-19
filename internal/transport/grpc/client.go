package grpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"x5_test/internal/api/proto/gen"
)

type FulfillmentClient struct {
	client gen.FulfillmentServiceClient
	conn   *grpc.ClientConn
}

func NewFulfillmentClient(addr string) (*FulfillmentClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to grpc server: %w", err)
	}

	return &FulfillmentClient{
		client: gen.NewFulfillmentServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *FulfillmentClient) Close() error {
	return c.conn.Close()
}

func (c *FulfillmentClient) ProcessOrder(ctx context.Context, orderID string) error {
	req := &gen.ProcessOrderRequest{
		OrderId: orderID,
	}

	_, err := c.client.ProcessOrder(ctx, req)
	if err != nil {
		return fmt.Errorf("grpc call failed: %w", err)
	}

	return nil
}

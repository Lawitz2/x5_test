package grpc

import (
	"context"
	"fmt"
	"x5_test/internal/api/proto/gen"
	"x5_test/internal/domain"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

func (c *FulfillmentClient) ProcessOrder(ctx context.Context, orderID string, items []domain.Item) error {
	protoItems := make([]*gen.OrderItem, len(items))
	for i, item := range items {
		protoItems[i] = &gen.OrderItem{
			Sku: item.SKU,
			Qty: int32(item.Quantity),
		}
	}

	req := &gen.ProcessOrderRequest{
		OrderId: orderID,
		Items:   protoItems,
	}

	_, err := c.client.ProcessOrder(ctx, req)
	if err != nil {
		return fmt.Errorf("grpc call failed: %w", err)
	}

	return nil
}

package grpc

import (
	"context"
	"x5_test/internal/api/proto/gen"
	"x5_test/internal/service"

	"google.golang.org/grpc"
)

type FulfillmentServer struct {
	gen.UnimplementedFulfillmentServiceServer
	service *service.FulfillmentService
}

func NewFulfillmentServer(grpcServer *grpc.Server, svc *service.FulfillmentService) {
	srv := &FulfillmentServer{
		service: svc,
	}
	gen.RegisterFulfillmentServiceServer(grpcServer, srv)
}

func (s *FulfillmentServer) ProcessOrder(ctx context.Context, req *gen.ProcessOrderRequest) (*gen.ProcessOrderResponse, error) {
	err := s.service.ProcessOrder(ctx, req.GetOrderId())
	if err != nil {
		return nil, err
	}

	return &gen.ProcessOrderResponse{
		Status:  "SUCCESS",
		Message: "Order processed successfully",
	}, nil
}

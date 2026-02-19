package grpc

import (
	"context"
	gen2 "x5_test/internal/api/proto/gen"
	"x5_test/internal/service"

	"google.golang.org/grpc"
)

type FulfillmentServer struct {
	gen2.UnimplementedFulfillmentServiceServer
	service *service.FulfillmentService
}

func NewFulfillmentServer(grpcServer *grpc.Server, svc *service.FulfillmentService) {
	srv := &FulfillmentServer{
		service: svc,
	}
	gen2.RegisterFulfillmentServiceServer(grpcServer, srv)
}

func (s *FulfillmentServer) ProcessOrder(ctx context.Context, req *gen2.ProcessOrderRequest) (*gen2.ProcessOrderResponse, error) {
	err := s.service.ProcessOrder(ctx, req.GetOrderId())
	if err != nil {
		return nil, err
	}

	return &gen2.ProcessOrderResponse{
		Status:  "SUCCESS",
		Message: "Order processed successfully",
	}, nil
}

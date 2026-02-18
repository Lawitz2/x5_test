package grpc

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		return nil, status.Errorf(codes.Internal, "processing failed: %v", err)
	}

	return &gen2.ProcessOrderResponse{
		Status:  "SUCCESS",
		Message: "Order processed successfully",
	}, nil
}

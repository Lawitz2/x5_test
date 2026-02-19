// cmd/app/main.go

package main

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"x5_test/config"
	"x5_test/internal/domain"

	"x5_test/internal/api"
	"x5_test/internal/repository/postgres"
	"x5_test/internal/service"
	transportgrpc "x5_test/internal/transport/grpc"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	pool, err := pgxpool.New(context.Background(), cfg.DBConfig.ConnString)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	err = goose.Up(stdlib.OpenDBFromPool(pool), cfg.MigrationsDir)
	if err != nil {
		log.Fatal(err)
	}

	repo := postgres.NewRepository(pool)

	orderSvc := service.NewOrderService(repo)
	fulfillmentSvc := service.NewFulfillmentService(repo)

	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("failed to listen for gRPC: %v", err)
	}

	grpcServer := grpc.NewServer()
	transportgrpc.NewFulfillmentServer(grpcServer, fulfillmentSvc)

	go func() {
		log.Printf("gRPC server listening at %v", lis.Addr())
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	e := echo.New()
	h := api.NewHandler(orderSvc, cfg.PageLimit)
	h.Register(e)

	go func() {
		log.Printf("HTTP server starting on :%s", cfg.HTTPPort)
		if err := e.Start(":" + cfg.HTTPPort); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("failed to start HTTP server: %v", err)
		}
	}()

	grpcClient, err := transportgrpc.NewFulfillmentClient("localhost:" + cfg.GRPCPort)
	if err != nil {
		log.Fatalf("failed to create grpc client: %v", err)
	}

	ctx, cancelProcessor := context.WithCancel(context.Background())
	go func() {
		orderProcessor(ctx, repo, grpcClient, cfg.PageLimit)
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	cancelProcessor()

	log.Println("Shutting down servers...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	grpcServer.GracefulStop()

	if err := grpcClient.Close(); err != nil {
		log.Printf("gRPC client close error: %v", err)
	}
	pool.Close()

	log.Println("Server exited properly")
}

// Поллит базу каждые 5 сек на наличие новых заказов, обрабатывает их если таковые есть.
func orderProcessor(
	ctx context.Context,
	repo *postgres.Repository,
	client *transportgrpc.FulfillmentClient,
	limit int) {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			orders, err := repo.ListOrders(ctx, "", domain.StatusNew, limit)
			if err != nil {
				log.Printf("failed to list orders: %v", err)
				continue
			}

			// По-хорошему обработка заказов должна выполняться в отдельных го-рутинах
			// (воркер пул), из-за простоты задания решил сделать в этом же треде.
			for _, order := range orders {
				ctxTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
				err = client.ProcessOrder(ctxTimeout, order.ID.String())
				cancel()
				if err != nil {
					log.Printf("failed to process order: %v", err)
				}
			}
		case <-ctx.Done():
			log.Printf("order processor is stopping")
			return
		}
	}
}

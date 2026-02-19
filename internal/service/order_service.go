package service

import (
	"context"
	"fmt"
	"time"

	"x5_test/internal/domain"

	"github.com/google/uuid"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *domain.Order) error
	GetOrder(ctx context.Context, id uuid.UUID) (*domain.Order, error)
	ListOrders(ctx context.Context, customerID string, status domain.OrderStatus, limit int) ([]domain.Order, error)
}

type FulfillmentClient interface {
	ProcessOrder(ctx context.Context, orderID string, items []domain.Item) error
}

type OrderService struct {
	repo OrderRepository
}

func NewOrderService(repo OrderRepository) *OrderService {
	return &OrderService{
		repo: repo,
	}
}

// CreateOrder создает новый заказ.
func (s *OrderService) CreateOrder(ctx context.Context, customerID string, items []domain.Item) (*domain.Order, error) {
	order := &domain.Order{
		ID:         uuid.New(),
		CustomerID: customerID,
		Items:      items,
		Status:     domain.StatusNew,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.repo.CreateOrder(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to create order in repo: %w", err)
	}

	return order, nil
}

// GetOrder возвращает заказ по его идентификатору.
func (s *OrderService) GetOrder(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	return s.repo.GetOrder(ctx, id)
}

// ListOrders возвращает список заказов с возможностью фильтрации.
func (s *OrderService) ListOrders(
	ctx context.Context,
	customerID string,
	status domain.OrderStatus,
	limit int,
) ([]domain.Order, error) {
	return s.repo.ListOrders(ctx, customerID, status, limit)
}

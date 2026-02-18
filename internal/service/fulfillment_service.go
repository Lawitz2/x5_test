package service

import (
	"context"
	"log"
	"math/rand"
	"time"

	"x5_test/internal/domain"

	"github.com/google/uuid"
)

type FulfillmentRepository interface {
	UpdateOrderStatus(ctx context.Context, id uuid.UUID, status domain.OrderStatus) error
}

type FulfillmentService struct {
	repo FulfillmentRepository
}

func NewFulfillmentService(repo FulfillmentRepository) *FulfillmentService {
	return &FulfillmentService{repo: repo}
}

func (s *FulfillmentService) ProcessOrder(ctx context.Context, orderID string) error {
	id, err := uuid.Parse(orderID)
	if err != nil {
		return err
	}

	if err := s.repo.UpdateOrderStatus(ctx, id, domain.StatusProcessing); err != nil {
		return err
	}

	// Имитируем работу (200ms)
	time.Sleep(time.Millisecond * 200)

	// Определяем финальный статус (70% успех, 30% ошибка)
	finalStatus := domain.StatusFulfilled
	if rand.Float32() < 0.3 {
		finalStatus = domain.StatusFailed
	}

	if err := s.repo.UpdateOrderStatus(ctx, id, finalStatus); err != nil {
		log.Printf("failed to update final status for order %s: %v", orderID, err)
		return err
	}

	log.Printf("Order %s processed with status %s", orderID, finalStatus)
	return nil
}

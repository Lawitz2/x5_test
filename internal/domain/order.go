package domain

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	StatusNew        OrderStatus = "NEW"
	StatusProcessing OrderStatus = "PROCESSING"
	StatusFulfilled  OrderStatus = "FULFILLED"
	StatusFailed     OrderStatus = "FAILED"
)

type Item struct {
	SKU      string `json:"sku"`
	Quantity int    `json:"qty"`
}

type Order struct {
	ID         uuid.UUID   `json:"id"`
	CustomerID string      `json:"customer_id"`
	Items      []Item      `json:"items"`
	Status     OrderStatus `json:"status"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

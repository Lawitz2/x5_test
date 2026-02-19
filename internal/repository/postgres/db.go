package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"x5_test/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) CreateOrder(ctx context.Context, order *domain.Order) error {
	itemsJSON, err := json.Marshal(order.Items)
	if err != nil {
		return fmt.Errorf("failed to marshal items: %w", err)
	}

	query := `
		INSERT INTO orders (id, customer_id, items, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err = r.pool.Exec(ctx, query,
		order.ID,
		order.CustomerID,
		itemsJSON,
		order.Status,
		order.CreatedAt,
		order.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
	}

	return nil
}

func (r *Repository) GetOrder(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	query := `
		SELECT id, customer_id, items, status, created_at, updated_at
		FROM orders
		WHERE id = $1
	`
	var order domain.Order
	var itemsJSON []byte

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&order.ID,
		&order.CustomerID,
		&itemsJSON,
		&order.Status,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	if err := json.Unmarshal(itemsJSON, &order.Items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal items: %w", err)
	}

	return &order, nil
}

func (r *Repository) ListOrders(
	ctx context.Context,
	customerID string,
	status domain.OrderStatus,
	limit int,
) ([]domain.Order, error) {
	query := `
		SELECT id, customer_id, items, status, created_at, updated_at
		FROM orders
		WHERE ($1 = '' OR customer_id = $1)
		  AND ($2 = '' OR status::text = $2)
		ORDER BY created_at DESC
		LIMIT $3
	`
	rows, err := r.pool.Query(ctx, query, customerID, string(status), limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var order domain.Order
		var itemsJSON []byte
		err := rows.Scan(
			&order.ID,
			&order.CustomerID,
			&itemsJSON,
			&order.Status,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}

		if err := json.Unmarshal(itemsJSON, &order.Items); err != nil {
			return nil, fmt.Errorf("failed to unmarshal items: %w", err)
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return orders, nil
}

func (r *Repository) UpdateOrderStatus(ctx context.Context, id uuid.UUID, status domain.OrderStatus) error {
	query := `
		UPDATE orders
		SET status = $1, updated_at = NOW()
		WHERE id = $2
	`
	result, err := r.pool.Exec(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("order with id %s not found for status update", id)
	}

	return nil
}

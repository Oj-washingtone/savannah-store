package repocitory

import (
	"context"
	"time"

	"github.com/Oj-washingtone/savannah-store/internal/database"
	"github.com/Oj-washingtone/savannah-store/internal/model"
	"github.com/google/uuid"
)

type OrdersRepository interface {
	Create(ctx context.Context, order *model.Orders) (*model.Orders, error)
	GetById(ctx context.Context, id uuid.UUID) (*model.Orders, error)
	GetByUser(ctx context.Context, userId uuid.UUID) ([]*model.Orders, error)
	UpdateStatus(ctx context.Context, orderId uuid.UUID, status model.OrderStatus) error
	UpdatePaidStatus(ctx context.Context, orderId uuid.UUID, paid bool) error
	GetAll(ctx context.Context) ([]*model.Orders, error)
}

type ordersRepository struct {
	db *database.DB
}

func NewOrdersRepository() OrdersRepository {
	return &ordersRepository{db: database.GetDB()}
}

func (r *ordersRepository) Create(ctx context.Context, order *model.Orders) (*model.Orders, error) {
	query := `
		INSERT INTO orders (id, user_id, total, paid)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, status, total, paid, created_at, updated_at
	`

	err := r.db.Pool.QueryRow(ctx, query,
		order.ID,
		order.UserID,
		order.Total,
		order.Paid,
	).Scan(
		&order.ID,
		&order.UserID,
		&order.Status,
		&order.Total,
		&order.Paid,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return order, nil
}

func (r *ordersRepository) GetById(ctx context.Context, id uuid.UUID) (*model.Orders, error) {
	query := `SELECT id, user_id, status, total, paid, created_at, updated_at
			  FROM orders WHERE id = $1`

	row := r.db.Pool.QueryRow(ctx, query, id)

	var order model.Orders
	err := row.Scan(
		&order.ID,
		&order.UserID,
		&order.Status,
		&order.Total,
		&order.Paid,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *ordersRepository) GetByUser(ctx context.Context, userId uuid.UUID) ([]*model.Orders, error) {
	query := `SELECT id, user_id, status, total, paid, created_at, updated_at
			  FROM orders WHERE user_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.Pool.Query(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*model.Orders
	for rows.Next() {
		order := &model.Orders{}
		if err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.Status,
			&order.Total,
			&order.Paid,
			&order.CreatedAt,
			&order.UpdatedAt,
		); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (r *ordersRepository) UpdateStatus(ctx context.Context, orderId uuid.UUID, status model.OrderStatus) error {
	query := `UPDATE orders SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Pool.Exec(ctx, query, status, time.Now(), orderId)
	return err
}

func (r *ordersRepository) UpdatePaidStatus(ctx context.Context, orderId uuid.UUID, paid bool) error {
	query := `UPDATE orders SET paid = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Pool.Exec(ctx, query, paid, time.Now(), orderId)
	return err
}

func (r *ordersRepository) GetAll(ctx context.Context) ([]*model.Orders, error) {
	query := `SELECT id, user_id, status, total, paid, created_at, updated_at
			  FROM orders ORDER BY created_at DESC`

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*model.Orders
	for rows.Next() {
		order := &model.Orders{}
		if err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.Status,
			&order.Total,
			&order.Paid,
			&order.CreatedAt,
			&order.UpdatedAt,
		); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

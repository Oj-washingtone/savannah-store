package repocitory

import (
	"context"
	"strconv"
	"strings"

	"github.com/Oj-washingtone/savannah-store/internal/database"
	"github.com/Oj-washingtone/savannah-store/internal/model"
)

type OrderItemsRepository interface {
	Create(ctx context.Context, item *model.OrderItems) error

	CreateBulk(ctx context.Context, items []*model.OrderItems) error
}

type orderItemsRepository struct {
	db *database.DB
}

func NewOrderItemsRepository() OrderItemsRepository {
	return &orderItemsRepository{db: database.GetDB()}
}

func (r *orderItemsRepository) Create(ctx context.Context, item *model.OrderItems) error {
	query := `
		INSERT INTO order_items (id, order_id, product_id, quantity, price)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		item.ID,
		item.OrderID,
		item.ProductID,
		item.Quantity,
		item.Price,
	)

	return err
}

func (r *orderItemsRepository) CreateBulk(ctx context.Context, items []*model.OrderItems) error {
	if len(items) == 0 {
		return nil
	}

	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `INSERT INTO order_items (id, order_id, product_id, quantity, price) VALUES `
	args := []interface{}{}
	for i, item := range items {
		start := i*5 + 1
		query += `($` + strconv.Itoa(start) + `,$` + strconv.Itoa(start+1) + `,$` + strconv.Itoa(start+2) + `,$` + strconv.Itoa(start+3) + `,$` + strconv.Itoa(start+4) + `),`
		args = append(args, item.ID, item.OrderID, item.ProductID, item.Quantity, item.Price)
	}
	query = strings.TrimRight(query, ",")

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

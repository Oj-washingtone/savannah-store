package repocitory

import (
	"context"

	"github.com/Oj-washingtone/savannah-store/internal/database"
	"github.com/Oj-washingtone/savannah-store/internal/model"
	"github.com/google/uuid"
)

type CartItemsRepository interface {
	AddItem(ctx context.Context, item *model.CartItem) error
	RemoveItem(ctx context.Context, itemId uuid.UUID) error
	GetItems(ctx context.Context, cartId uuid.UUID) ([]*model.CartItem, error)
	UpdateQuantity(ctx context.Context, itemId uuid.UUID, quantity int) error
	Exists(ctx context.Context, cartId uuid.UUID, productId string) (bool, error)
	ClearCart(ctx context.Context, cartId uuid.UUID) error
}

type cartItemsRepository struct {
	db *database.DB
}

func NewCartItemsRepository() CartItemsRepository {
	return &cartItemsRepository{db: database.GetDB()}
}

func (r *cartItemsRepository) AddItem(ctx context.Context, item *model.CartItem) error {

	query := `INSERT INTO cart_items (id, cart_id, product_id, quantity, price) VALUES ($1, $2, $3, $4, $5) RETURNING created_at, updated_at`

	return r.db.Pool.QueryRow(ctx, query,
		item.ID,
		item.CartId,
		item.ProductId,
		item.Quantity,
		item.Price,
	).Scan(&item.CreatedAt, &item.UpdatedAt)

}

func (r *cartItemsRepository) RemoveItem(ctx context.Context, itemId uuid.UUID) error {
	query := `DELETE FROM cart_items WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, itemId)
	return err
}

func (r *cartItemsRepository) GetItems(ctx context.Context, cartId uuid.UUID) ([]*model.CartItem, error) {
	query := `SELECT id, cart_id, product_id, quantity, price, created_at, updated_at FROM cart_items WHERE cart_id = $1`
	rows, err := r.db.Pool.Query(ctx, query, cartId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*model.CartItem
	for rows.Next() {
		item := &model.CartItem{}
		if err := rows.Scan(&item.ID, &item.CartId, &item.ProductId, &item.Quantity, &item.Price, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *cartItemsRepository) UpdateQuantity(ctx context.Context, itemId uuid.UUID, quantity int) error {
	query := `UPDATE cart_items SET quantity = $1 WHERE id = $2`
	_, err := r.db.Pool.Exec(ctx, query, quantity, itemId)
	return err
}

func (r *cartItemsRepository) Exists(ctx context.Context, cartId uuid.UUID, productId string) (bool, error) {
	query := `SELECT EXISTS (
        SELECT 1 FROM cart_items WHERE cart_id = $1 AND product_id = $2
    )`
	var exists bool
	err := r.db.Pool.QueryRow(ctx, query, cartId, productId).Scan(&exists)
	return exists, err
}

func (r *cartItemsRepository) ClearCart(ctx context.Context, cartId uuid.UUID) error {
	query := `DELETE FROM cart_items WHERE cart_id = $1`
	_, err := r.db.Pool.Exec(ctx, query, cartId)
	return err
}

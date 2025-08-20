package repocitory

import (
	"context"

	"github.com/Oj-washingtone/savannah-store/internal/database"
	"github.com/Oj-washingtone/savannah-store/internal/model"
	"github.com/google/uuid"
)

type ShoppingCartRepository interface {
	CreateCart(ctx context.Context, userId uuid.UUID) (*model.Cart, error)
	GetShoppingCart(ctx context.Context, userId uuid.UUID) (*model.Cart, error)
}

type shoppingCartRepository struct {
	db *database.DB
}

func NewShoppingCartRepository() ShoppingCartRepository {
	return &shoppingCartRepository{db: database.GetDB()}
}

func (r *shoppingCartRepository) CreateCart(ctx context.Context, userId uuid.UUID) (*model.Cart, error) {
	cartId := uuid.New()
	query := `INSERT INTO carts (id, user_id) VALUES ($1, $2) RETURNING id, user_id, created_at, updated_at`

	var cart model.Cart
	err := r.db.Pool.QueryRow(ctx, query, cartId, userId).Scan(
		&cart.ID,
		&cart.UserId,
		&cart.CreatedAt,
		&cart.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &cart, nil
}

func (r *shoppingCartRepository) GetShoppingCart(ctx context.Context, userId uuid.UUID) (*model.Cart, error) {
	query := `SELECT id, user_id, created_at, updated_at FROM carts WHERE user_id = $1 AND deleted_at IS NULL`

	var cart model.Cart
	err := r.db.Pool.QueryRow(ctx, query, userId).Scan(
		&cart.ID,
		&cart.UserId,
		&cart.CreatedAt,
		&cart.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &cart, nil
}

package repocitory

import (
	"context"

	"github.com/Oj-washingtone/savannah-store/internal/database"
	"github.com/Oj-washingtone/savannah-store/internal/model"
	"github.com/google/uuid"
)

type ProductRepository interface {
	Create(ctx context.Context, product *model.Product) error
	GetById(ctx context.Context, id uuid.UUID) (*model.Product, error)
	List(ctx context.Context, limit, offset int) ([]model.Product, error)
	Update(ctx context.Context, product *model.Product) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type productRepository struct {
	db *database.DB
}

func NewProductRepository() ProductRepository {
	return &productRepository{db: database.GetDB()}
}

func (r *productRepository) Create(ctx context.Context, product *model.Product) error {
	query := `
		INSERT INTO products (id, category_id, name, description, price, stock)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at
	`

	return r.db.Pool.QueryRow(ctx, query,
		product.ID,
		product.CategoryID,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
	).Scan(&product.CreatedAt, &product.UpdatedAt)
}

func (r *productRepository) GetById(ctx context.Context, id uuid.UUID) (*model.Product, error) {
	query := `
		SELECT id, category_id, name, description, price, stock, created_at, updated_at
		FROM products
		WHERE id = $1 AND deleted_at IS NULL
	`

	row := r.db.Pool.QueryRow(ctx, query, id)

	var product model.Product
	err := row.Scan(
		&product.ID,
		&product.CategoryID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Stock,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *productRepository) List(ctx context.Context, limit, offset int) ([]model.Product, error) {
	query := `
		SELECT id, category_id, name, description, price, stock, created_at, updated_at
		FROM products
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(
			&p.ID,
			&p.CategoryID,
			&p.Name,
			&p.Description,
			&p.Price,
			&p.Stock,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func (r *productRepository) Update(ctx context.Context, product *model.Product) error {
	query := `
		UPDATE products
		SET category_id = $1, name = $2, description = $3, price = $4, stock = $5, updated_at = now()
		WHERE id = $6
		RETURNING updated_at
	`

	return r.db.Pool.QueryRow(ctx, query,
		product.CategoryID,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
		product.ID,
	).Scan(&product.UpdatedAt)
}

func (r *productRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE products
		SET deleted_at = now()
		WHERE id = $1 AND deleted_at IS NULL
	`
	_, err := r.db.Pool.Exec(ctx, query, id)
	return err
}

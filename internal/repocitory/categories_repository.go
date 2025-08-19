package repocitory

import (
	"context"

	"github.com/Oj-washingtone/savannah-store/internal/database"
	"github.com/Oj-washingtone/savannah-store/internal/model"
	"github.com/google/uuid"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *model.ProductCategory) error
	GetById(ctx context.Context, id uuid.UUID) (*model.ProductCategory, error)
	List(ctx context.Context, limit, offset int) ([]model.ProductCategory, error)
	Update(ctx context.Context, category *model.ProductCategory) error
}

type categoryRepository struct {
	db *database.DB
}

func NewCategoryRepository() CategoryRepository {
	return &categoryRepository{db: database.GetDB()}
}

func (r *categoryRepository) Create(ctx context.Context, category *model.ProductCategory) error {
	query := `
		INSERT INTO categories (id, name, parent_id)
		VALUES ($1, $2, $3)
		RETURNING created_at, updated_at
	`

	return r.db.Pool.QueryRow(ctx, query,
		category.ID, category.Name, category.ParentId,
	).Scan(&category.CreatedAt, &category.UpdatedAt)
}

func (r *categoryRepository) GetById(ctx context.Context, id uuid.UUID) (*model.ProductCategory, error) {
	query := `
		SELECT id, name, parent_id, created_at, updated_at
		FROM categories
		WHERE id = $1
	`

	row := r.db.Pool.QueryRow(ctx, query, id)

	var category model.ProductCategory
	err := row.Scan(
		&category.ID,
		&category.Name,
		&category.ParentId,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *categoryRepository) List(ctx context.Context, limit, offset int) ([]model.ProductCategory, error) {
	query := `
		SELECT id, name, parent_id, created_at, updated_at
		FROM categories
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []model.ProductCategory
	for rows.Next() {
		var c model.ProductCategory
		if err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.ParentId,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	return categories, nil
}

func (r *categoryRepository) Update(ctx context.Context, category *model.ProductCategory) error {
	query := `
		UPDATE categories
		SET name = $1, parent_id = $2, updated_at = now()
		WHERE id = $3
		RETURNING updated_at
	`

	return r.db.Pool.QueryRow(ctx, query,
		category.Name,
		category.ParentId,
		category.ID,
	).Scan(&category.UpdatedAt)
}

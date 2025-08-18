package repocitory

import (
	"context"

	"github.com/Oj-washingtone/savannah-store/internal/database"
	"github.com/Oj-washingtone/savannah-store/internal/model"
	"github.com/google/uuid"
)

type userRepository struct {
	db *database.DB
}

func NewUserRepository() *userRepository {
	return &userRepository{db: database.GetDB()}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (id, name, email, role, auth0_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at, updated_at
	`

	return r.db.Pool.QueryRow(ctx, query,
		user.ID, user.Name, user.Email, user.Role, user.Auth0Id,
	).Scan(&user.CreatedAt, &user.UpdatedAt)
}

func (r *userRepository) GetById(context context.Context, id uuid.UUID) (*model.User, error) {
	row := r.db.Pool.QueryRow(context,
		`SELECT id, name, email, role FROM users WHERE id=$1`,
		id,
	)

	user := &model.User{}

	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Role)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	row := r.db.Pool.QueryRow(ctx,
		`SELECT id, name, email, role, created_at, updated_at FROM users WHERE email=$1`,
		email,
	)

	user := &model.User{}
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}

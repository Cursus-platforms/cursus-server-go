package role

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	"github.com/Cursus-platforms/cursus-server-go/internal/model"
)

type Repository interface {
	GetByName(ctx context.Context, name string) (model.Role, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.Role, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetByName(ctx context.Context, name string) (model.Role, error) {
	var role model.Role
	query := `SELECT id, name, description, created_at, updated_at FROM roles WHERE name = $1`

	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt,
	)

	if err != nil {
		return model.Role{}, err
	}

	return role, nil
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (model.Role, error) {
	var role model.Role
	query := `SELECT id, name, description, created_at, updated_at FROM roles WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt,
	)

	if err != nil {
		return model.Role{}, err
	}

	return role, nil
}

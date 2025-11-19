package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/Cursus-platforms/cursus-server-go/internal/model"
)

type Repository interface {
	Create(ctx context.Context, user model.User) (model.User, error)
	GetByEmail(ctx context.Context, email string) (model.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.User, error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, user model.User) (model.User, error) {
	query := `INSERT INTO users (fullname, email, password, role_id)
			VALUES ($1, $2, $3, $4)
			RETURN id, created_at, updated_at`

	var createdUser model.User

	err := r.db.QueryRowxContext(ctx, query, user.FullName, user.Email, user.Password, user.RoleID).StructScan(&createdUser)

	if err != nil {
		return model.User{}, err
	}

	return createdUser, nil
}

func (r *repository) GetByEmail(ctx context.Context, email string) (model.User, error) {
	var user model.User
	query := `SELECT id, fullname, email, password, role_id, created_at, updated_at 
				FROM users WHERE email=$1`

	err := r.db.GetContext(ctx, &user, query, email)

	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (model.User, error) {
	var user model.User
	query := `SELECT id, fullname, email, role_id, created_at, updated_at, FROM users WHERE id = $1`

	err := r.db.GetContext(ctx, &user, query, id)

	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

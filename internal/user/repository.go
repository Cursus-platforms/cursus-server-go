package user

import (
	"context"
	"database/sql"

	"github.com/Cursus-platforms/cursus-server-go/internal/model"
)

type Repository interface {
	Create(ctx context.Context, user model.User)
	GetByEmail(ctx context.Context, email string)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &re
}

package model

import (
	"time"

	"github.com/google/uuid"
)

const (
	ROLE_ADMIN      = "ADMIN"
	ROLE_INSTRUCTOR = "INSTRUCTOR"
	ROLE_STUDENT    = "STUDENT"
)

type Role struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description *string   `json:"description,omitempty" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

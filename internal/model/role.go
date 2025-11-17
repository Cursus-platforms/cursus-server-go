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
	ID uuid.UUID
	Name string
	Description *string
	CreatedAt time.Time
	UpdatedAt time.Time
}


package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID              uuid.UUID
	Fullname        *string
	Email           *string
	Password        *string
	Image           *string
	Subscribe       *int
	IsActive        *bool
	Status          *string
	Description     *string
	Facebook        *string
	Twitter         *string
	Linkedin        *string
	Youtube         *string
	RejectionReason *string
	CategoryID      *uuid.UUID
	SubCategoryID   *uuid.UUID
	RoleID          *uuid.UUID
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

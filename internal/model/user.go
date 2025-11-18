package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	Fullname        *string    `json:"fullname,omitempty" db:"fullname"`
	Email           *string    `json:"email,omitempty" db:"email"`
	Password        *string    `json:"-" db:"password"`
	Image           *string    `json:"image,omitempty" db:"image"`
	Subscribe       *int       `json:"subscribe,omitempty" db:"subscribe"`
	IsActive        *bool      `json:"is_active,omitempty" db:"is_active"`
	Status          *string    `json:"status,omitempty" db:"status"`
	Description     *string    `json:"description,omitempty" db:"description"`
	Facebook        *string    `json:"facebook,omitempty" db:"facebook"`
	Twitter         *string    `json:"twitter,omitempty" db:"twitter"`
	Linkedin        *string    `json:"linkedin,omitempty" db:"linkedin"`
	Youtube         *string    `json:"youtube,omitempty" db:"youtube"`
	RejectionReason *string    `json:"rejection_reason,omitempty" db:"rejection_reason"`
	CategoryID      *uuid.UUID `json:"category_id,omitempty" db:"category_id"`
	SubCategoryID   *uuid.UUID `json:"sub_category_id,omitempty" db:"sub_category_id"`
	RoleID          *uuid.UUID `json:"role_id,omitempty" db:"role_id"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}

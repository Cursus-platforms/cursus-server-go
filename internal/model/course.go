package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Course struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	UserID        *uuid.UUID `json:"user_id,omitempty" db:"user_id"`
	SubCategoryID *uuid.UUID `json:"sub_category_id,omitempty" db:"sub_category_id"`

	Title string  `json:"title" db:"title"`
	Price float64 `json:"price" db:"price"`

	Description      *string `json:"description,omitempty" db:"description"`
	ShortDescription *string `json:"short_description,omitempty" db:"short_description"`
	Requirements     *string `json:"requirements,omitempty" db:"requirements"`
	StudentLearn     *string `json:"student_learn,omitempty" db:"student_learn"`
	RejectionReason  *string `json:"rejection_reason,omitempty" db:"rejection_reason"`

	Image      *string `json:"image,omitempty" db:"image"`
	IntroVideo *string `json:"intro_video,omitempty" db:"intro_video"`
	Slug       *string `json:"slug,omitempty" db:"slug"`

	Status     string     `json:"status" db:"status"`
	ApprovedAt *time.Time `json:"approved_at,omitempty" db:"approved_at"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`

	TotalSold     int     `json:"total_sold,omitempty" db:"total_sold"`
	TotalView     int     `json:"total_view,omitempty" db:"total_view"`
	TotalRating   int     `json:"total_rating,omitempty" db:"total_rating"`
	TotalEnrolled int     `json:"total_enrolled,omitempty" db:"total_enrolled"`
	Likes         int     `json:"likes,omitempty" db:"likes"`
	Dislikes      int     `json:"dislikes,omitempty" db:"dislikes"`
	Discount      float64 `json:"discount,omitempty" db:"discount"`

	IsDeleted     bool `json:"is_deleted,omitempty" db:"is_deleted"`
	IsBestseller  bool `json:"is_bestseller,omitempty" db:"is_bestseller"`
	RequireLogIn  bool `json:"require_log_in,omitempty" db:"require_log_in"`
	RequireEnroll bool `json:"require_enroll,omitempty" db:"require_enroll"`

	Levels   json.RawMessage `json:"levels,omitempty" db:"levels"`
	Captions json.RawMessage `json:"captions,omitempty" db:"captions"`
	Language json.RawMessage `json:"language,omitempty" db:"language"`
}

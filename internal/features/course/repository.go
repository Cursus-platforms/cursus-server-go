package course

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/Cursus-platforms/cursus-server-go/internal/model"
)

type Repository interface {
	Create(ctx context.Context, course model.Course) (model.Course, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.Course, error)
	Update(ctx context.Context, course model.Course) (model.Course, error)
	Delete(ctx context.Context, id uuid.UUID) error
	CheckSlugExists(ctx context.Context, slug string) (bool, error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, course model.Course) (model.Course, error) {
	query := `
        INSERT INTO courses (
            user_id, title, description, short_description, requirements, student_learn, 
            price, slug, sub_category_id, status
        ) 
        VALUES (
            :user_id, :title, :description, :short_description, :requirements, :student_learn, 
            :price, :slug, :sub_category_id, :status
        )`

	rows, err := r.db.NamedQueryContext(ctx, query, course)
	if err != nil {
		return model.Course{}, fmt.Errorf("repository: failed to create course: %w", err)
	}
	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("WARNING: Error closing rows: %v", err)
		}
	}(rows)

	if rows.Next() {
		err := rows.StructScan(&course)
		return course, err
	}

	return model.Course{}, errors.New("repository: failed to retrieve created course data")
}

// GetByID GetByID
func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (model.Course, error) {
	var course model.Course
	query := `SELECT * FROM courses WHERE id = $1 AND is_deleted = FALSE LIMIT 1`

	err := r.db.GetContext(ctx, &course, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Course{}, sql.ErrNoRows
		}
		return model.Course{}, fmt.Errorf("repository: failed to get course by ID: %w", err)
	}
	return course, nil
}

// Update Update
func (r *repository) Update(ctx context.Context, course model.Course) (model.Course, error) {
	query := `
        UPDATE courses SET
            title = :title,
            description = :description,
            short_description = :short_description,
            requirements = :requirements,
            student_learn = :student_learn,
            price = :price,
            slug = :slug,
            sub_category_id = :sub_category_id,
            status = :status,
            image = :image,
            intro_video = :intro_video,
            levels = :levels,
            captions = :captions,
            language = :language,
            approved_at = :approved_at,
            updated_at = NOW()
        WHERE id = :id AND is_deleted = FALSE`

	rows, err := r.db.NamedQueryContext(ctx, query, course)
	if err != nil {
		return model.Course{}, fmt.Errorf("repository: failed to update course: %w", err)
	}
	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("WARNING: Error closing rows: %v", err)
		}
	}(rows)

	if rows.Next() {
		err := rows.StructScan(&course)
		return course, err
	}

	return model.Course{}, sql.ErrNoRows
}

// Delete DELETE (Soft Delete)
func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE courses SET is_deleted = TRUE, updated_at = NOW() WHERE id = $1 AND is_deleted = FALSE`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("repository: failed to soft delete course: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("repository: failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *repository) CheckSlugExists(ctx context.Context, slug string) (bool, error) {
	var exists bool

	query := `SELECT EXISTS(SELECT 1 FROM courses WHERE slug = $1 AND is_deleted = FALSE)`

	err := r.db.GetContext(ctx, &exists, query, slug)

	if err != nil {
		return false, fmt.Errorf("repository: failed to check slug existence in DB: %w", err)
	}

	return exists, nil
}

package course

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Cursus-platforms/cursus-server-go/internal/lib/utils"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/Cursus-platforms/cursus-server-go/internal/lib/slug"
	"github.com/Cursus-platforms/cursus-server-go/internal/model"
)

const CacheDuration = 30 * time.Minute

type Service interface {
	CreateCourse(ctx context.Context, course model.Course) (model.Course, error)
	GetCourseByID(ctx context.Context, id uuid.UUID) (model.Course, error)
	// GetCourseBySlug(ctx context.Context, slug string) (model.Course, error)
}

type service struct {
	repo Repository
	rdb  *redis.Client
}

func NewService(repo Repository, rdb *redis.Client) Service {
	return &service{repo: repo, rdb: rdb}
}

func (s *service) CreateCourse(ctx context.Context, course model.Course) (model.Course, error) {
	baseSlug := slug.Generate(course.Title)
	finalSlug := baseSlug

	for i := 0; i < 5; i++ {
		exists, err := s.repo.CheckSlugExists(ctx, finalSlug)
		if err != nil {
			return model.Course{}, fmt.Errorf("service: error checking slug uniqueness: %w", err)
		}

		if !exists {
			break
		}

		suffix := utils.RandomString(4)
		finalSlug = fmt.Sprintf("%s-%s", baseSlug, suffix)

		if i == 4 {
			return model.Course{}, errors.New("could not generate unique slug after 5 attempts")
		}
	}

	course.Slug = &finalSlug
	course.Status = "draft"

	createdCourse, err := s.repo.Create(ctx, course)
	if err != nil {
		return model.Course{}, fmt.Errorf("service: database error during course creation: %w", err)
	}

	return createdCourse, nil
}

func (s *service) GetCourseByID(ctx context.Context, id uuid.UUID) (model.Course, error) {
	cacheKey := fmt.Sprintf("course:%s", id.String())

	val, err := s.rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		var course model.Course
		if err := json.Unmarshal([]byte(val), &course); err == nil {
			log.Printf("CACHE HIT: Course %s", id.String())
			return course, nil
		}
	} else if !errors.Is(err, redis.Nil) {
		log.Printf("Warning: Redis error on GET course: %v", err)
	}

	course, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return model.Course{}, err
	}

	if course.Status == "approved" {
		jsonData, err := json.Marshal(course)
		if err == nil {
			s.rdb.Set(ctx, cacheKey, jsonData, CacheDuration)
			log.Printf("CACHE REFRESHED: Course %s", id.String())
		}
	}

	return course, nil
}

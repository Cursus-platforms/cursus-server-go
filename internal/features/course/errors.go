package course

import "errors"

var (
	ErrPriceRequired = errors.New("price is required and must be positive")
	// ErrSlugGenerationFailed = errors.New("could not generate unique slug after 5 attempts")
	// ErrCourseNotFound = errors.New("course not found")
)

package validator

import (
	"github.com/go-playground/validator/v10"
)

var GlobalValidator *validator.Validate

func init() {
	GlobalValidator = validator.New()
}

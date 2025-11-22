package validator

import (
	"errors"
	"fmt"

	"github.com/Cursus-platforms/cursus-server-go/internal/lib/response"
	"github.com/go-playground/validator/v10"
)

func FormatValidationErrors(err error) []response.ValidationErrorDetail {
	var details []response.ValidationErrorDetail
	var fieldErrors validator.ValidationErrors
	if errors.As(err, &fieldErrors) {
		for _, e := range fieldErrors {
			details = append(details, response.ValidationErrorDetail{
				Field: e.Field(),
				Error: fmt.Sprintf("Field '%s' failed on the '%s' tag", e.Field(), e.Tag()),
			})
		}
	}
	return details
}

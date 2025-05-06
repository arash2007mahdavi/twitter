package validations

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Filed string `json:"field"`
	Tag   string `json:"tag"`
	Param string `json:"param"`
}

func GetValidationError(err error) *[]ValidationError {
	validationerrors := []ValidationError{}
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		for _, err := range err.(validator.ValidationErrors) {
			element := ValidationError{}
			element.Filed = err.Field()
			element.Param = err.Param()
			element.Tag = err.Tag()
			validationerrors = append(validationerrors, element)
		}
		return &validationerrors
	}
	return nil
}
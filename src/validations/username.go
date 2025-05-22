package validations

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func ValidateUsername(fld validator.FieldLevel) bool {
	val, ok := fld.Field().Interface().(string)
	if !ok {
		return false
	}
	res, err := regexp.MatchString(`^[a-zA-Z0-9_]{7,35}$`, val)
	if err != nil {
		panic("error in match string")
	}
	return res
}
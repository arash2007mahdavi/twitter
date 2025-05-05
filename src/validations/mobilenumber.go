package validations

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func ValidateMobileNumber(fld validator.FieldLevel) bool {
	val, ok := fld.Field().Interface().(string)
	if !ok {
		return false
	}
	res, err := regexp.MatchString(`^09[0-9]{9}$`, val)
	if err != nil {
		panic("error in match string")
	}
	return res
}
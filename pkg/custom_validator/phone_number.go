package custom_validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func ValidatePhoneNumber(fl validator.FieldLevel) bool {
	if match, err := regexp.MatchString(`^\+62*`, fl.Field().String()); err != nil || !match {
		return false
	}

	return true
}

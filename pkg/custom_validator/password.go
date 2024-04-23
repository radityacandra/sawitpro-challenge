package custom_validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func ValidatePassword(fl validator.FieldLevel) bool {
	// upperCase validation
	regexUppercase := regexp.MustCompile("^.*[A-Z].*$")
	// number validation
	regexNumber := regexp.MustCompile("^.*[0-9].*$")
	// symbol validation
	regexSymbol := regexp.MustCompile("^.*[-!$%^&*()_+|~=`{}\\[\\]:\";'<>?,.\\/].*$")

	if regexUppercase.MatchString(fl.Field().String()) &&
		regexNumber.MatchString(fl.Field().String()) && regexSymbol.MatchString(fl.Field().String()) {
		return true
	}

	return false
}

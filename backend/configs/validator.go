package configs

import (
	"regexp"
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	initValidatorOnce sync.Once
	validate          *validator.Validate
)

func passwordValidationFunc(fl validator.FieldLevel) bool {
	hasDigit := regexp.MustCompile(`[0-9]`)
	hasLower := regexp.MustCompile(`[a-z]`)
	hasUpper := regexp.MustCompile(`[A-Z]`)
	hasSpecial := regexp.MustCompile(`[\W]`)
	noSpace := regexp.MustCompile(`\s`)

	fieldValue := fl.Field().String()

	// Check all conditions
	return hasDigit.MatchString(fieldValue) &&
		hasLower.MatchString(fieldValue) &&
		hasUpper.MatchString(fieldValue) &&
		hasSpecial.MatchString(fieldValue) &&
		!noSpace.MatchString(fieldValue) &&
		len(fieldValue) >= 8
}

func GetValidator() *validator.Validate {
	initValidatorOnce.Do(func() {
		validate = validator.New()

		// Register custom validator
		validate.RegisterValidation("password", func(fl validator.FieldLevel) bool {
			return false
		})
		validate.RegisterValidation("password", passwordValidationFunc)
	})

	return validate
}

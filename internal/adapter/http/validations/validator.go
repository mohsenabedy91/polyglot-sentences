package validations

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"log"
	"regexp"
	"unicode"
)

func RegisterValidator(cfg config.Config) error {
	if val, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := val.RegisterValidation("regex_alpha", RegexAlpha, true); err != nil {
			return err
		}
		if err := val.RegisterValidation("password_complexity", PasswordComplexity, true); err != nil {
			return err
		}
		if err := val.RegisterValidation("token_length", func(fl validator.FieldLevel) bool {
			return TokenLength(fl, cfg.OTP.Digits)
		}, true); err != nil {
			return err
		}
	}

	return nil
}

// RegexAlpha This Go code defines a function to validates if a given string only contains alphabetic characters and spaces.
// It uses a regular expression to check for Unicode letters and marks.
// If the validation fails, it logs the error and returns false.
func RegexAlpha(field validator.FieldLevel) bool {
	value, ok := field.Field().Interface().(string)
	if !ok {
		return false
	}
	res, err := regexp.MatchString(`^[\p{L}\p{M} ]+$`, value)
	if err != nil {
		log.Print(err.Error())
	}
	return res
}

// PasswordComplexity This Go code defines a function to validate the complexity of a password based on certain criteria such as
// length, presence of uppercase letters, lowercase letters, digits, and special characters.
func PasswordComplexity(field validator.FieldLevel) bool {
	var (
		hasUpper   bool
		hasLower   bool
		hasDigit   bool
		hasSpecial bool
	)

	if value, ok := field.Field().Interface().(string); !ok {
		return false
	} else {
		for _, char := range value {
			switch {
			case unicode.IsUpper(char):
				hasUpper = true
			case unicode.IsLower(char):
				hasLower = true
			case unicode.IsDigit(char):
				hasDigit = true
			case unicode.IsPunct(char) || unicode.IsSpace(char):
				hasSpecial = true
			}
		}
	}

	return hasUpper && hasLower && hasDigit && hasSpecial
}

func TokenLength(field validator.FieldLevel, length int8) bool {
	if value, ok := field.Field().Interface().(string); !ok {
		return false
	} else {
		return len(value) == int(length)
	}
}

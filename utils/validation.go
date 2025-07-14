package utils

import (
	"regexp"

	"github.com/pquerna/otp/totp"
	"github.com/xeipuuv/gojsonschema"
)

// Valida fuerza de contraseña
func ValidPassword(pw string) bool {
	if len(pw) < 12 {
		return false
	}

	hasNumber := false
	hasSymbol := false

	for _, c := range pw {
		switch {
		case c >= '0' && c <= '9':
			hasNumber = true
		case (c >= 33 && c <= 47) || (c >= 58 && c <= 64) || (c >= 91 && c <= 96) || (c >= 123 && c <= 126):
			hasSymbol = true
		}
	}

	return hasNumber && hasSymbol
}

// Valida un objeto JSON contra un esquema
func ValidateJSONSchema(data interface{}, schemaJSON string) (bool, []string) {
	schemaLoader := gojsonschema.NewStringLoader(schemaJSON)
	documentLoader := gojsonschema.NewGoLoader(data)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return false, []string{"error al validar: " + err.Error()}
	}

	if result.Valid() {
		return true, nil
	}

	errors := []string{}
	for _, e := range result.Errors() {
		errors = append(errors, e.String())
	}
	return false, errors
}

func ValidateTOTP(secret, code string) bool {
	return totp.Validate(code, secret)
}
func IsValidEmail(email string) bool {
	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return regex.MatchString(email)
}

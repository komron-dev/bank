package val

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	isValidFullName = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString
)

func ValidateString(value string, minLength int, maxLength int) error {
	n := len(value)
	if n < minLength || n > maxLength {
		return fmt.Errorf("must contain %d-%d characters", minLength, maxLength)
	}
	return nil
}

func ValidateUsername(value string) error {
	if err := ValidateString(value, 3, 50); err != nil {
		return err
	}

	if !isValidUsername(value) {
		return fmt.Errorf("username can only contain lowercase letters, digits, or underscore")
	}

	return nil
}

func ValidateFullName(value string) error {
	if err := ValidateString(value, 3, 50); err != nil {
		return err
	}

	if !isValidFullName(value) {
		return fmt.Errorf("full name can only contain letters or space")
	}

	return nil
}

func ValidatePassword(value string) error {
	return ValidateString(value, 6, 50)
}

func ValidateEmail(value string) error {
	if err := ValidateString(value, 3, 50); err != nil {
		return err
	}

	_, err := mail.ParseAddress(value)
	if err != nil {
		return fmt.Errorf("not a valid email address")
	}

	return nil
}

func ValidateEmailId(value int64) error {
	if value <= 0 {
		return fmt.Errorf("must be a positive integer")
	}
	return nil
}

func ValidateSecretCode(value string) error {
	return ValidateString(value, 32, 128)
}

package util

import (
	"regexp"

	constant "github.com/nozzlium/belimang/internal/constants"
)

func ValidateEmailAddress(
	email string,
) error {
	emailRegex, err := regexp.Compile(
		`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`,
	)
	if err != nil {
		return err
	}

	if !emailRegex.MatchString(email) {
		return constant.ErrBadInput
	}

	return nil
}

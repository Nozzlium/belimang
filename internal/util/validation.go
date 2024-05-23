package util

import (
	"regexp"

	"github.com/nozzlium/belimang/internal/constant"
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

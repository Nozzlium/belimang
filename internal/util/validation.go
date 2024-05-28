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

func ValidateURL(url string) error {
	regex := `^[(http(s)?):\/\/(www\.)?a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)$`
	urlRegex, err := regexp.Compile(
		regex,
	)
	if err != nil {
		return err
	}

	if !urlRegex.MatchString(url) {
		return constant.ErrBadInput
	}

	return nil
}

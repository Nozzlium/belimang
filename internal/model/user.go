package model

import (
	"github.com/google/uuid"
	"github.com/nozzlium/belimang/internal/constant"
	"github.com/nozzlium/belimang/internal/util"
)

type User struct {
	ID       uuid.UUID
	Username string
	Email    string
	Password string
}

type UserRegisterBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (body UserRegisterBody) IsValid() (User, error) {
	var user User
	if unameLen := len(body.Username); unameLen < 5 ||
		unameLen > 30 {
		return user, constant.ErrInvalidBody
	}
	user.Username = body.Username

	if passLen := len(body.Password); passLen < 5 ||
		passLen > 30 {
		return user, constant.ErrBadInput
	}
	user.Password = body.Password

	if err := util.ValidateEmailAddress(body.Email); err != nil {
		return user, err
	}
	user.Email = body.Email

	return user, nil
}

type UserLoginBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (body *UserLoginBody) IsValid() (User, error) {
	var user User
	if unameLen := len(body.Username); unameLen < 5 ||
		unameLen > 30 {
		return user, constant.ErrInvalidBody
	}
	user.Username = body.Username

	if passLen := len(body.Password); passLen < 5 ||
		passLen > 30 {
		return user, constant.ErrBadInput
	}
	user.Password = body.Password

	return user, nil
}

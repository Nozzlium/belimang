package service

import (
	"context"

	"github.com/nozzlium/belimang/internal/model"
	"github.com/nozzlium/belimang/internal/repository"
)

type UserService struct {
	userRepository *repository.UserRepository
}

func NewUserService(
	userRepository *repository.UserRepository,
) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) RegisterAdmin(
	ctx context.Context,
	user model.User,
) (string, error) {
}

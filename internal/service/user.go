package service

import (
	"context"
	"encoding/base64"
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/nozzlium/belimang/internal/constant"
	"github.com/nozzlium/belimang/internal/model"
	"github.com/nozzlium/belimang/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepository *repository.UserRepository
	secret         string
	salt           int
}

func NewUserService(
	userRepository *repository.UserRepository,
	secret string,
	salt int,
) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) RegisterAdmin(
	ctx context.Context,
	user model.User,
) (string, error) {
	foundUsername, err := s.userRepository.FindUserUsername(
		ctx,
		user.Username,
	)
	if err != nil {
		if !errors.Is(
			err,
			constant.ErrNotFound,
		) {
			return "", err
		}
	}
	log.Println(
		foundUsername,
		user.Username,
	)

	if foundUsername == user.Username {
		return "", constant.ErrConflict
	}

	userId, err := uuid.NewV7()
	if err != nil {
		return "", err
	}

	generatedHashBytes, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password),
		s.salt,
	)
	if err != nil {
		return "", err
	}

	user.ID = userId
	user.Password = string(
		generatedHashBytes,
	)
	savedUser, err := s.userRepository.CreateAdmin(
		ctx,
		user,
	)
	if err != nil {
		return "", err
	}

	token, err := generateJwtToken(
		s.secret,
		savedUser,
	)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *UserService) LoginAdmin(
	ctx context.Context,
	user model.User,
) (string, error) {
	savedUser, err := s.userRepository.FindAdminByUsername(
		ctx,
		user,
	)
	if err != nil {
		if errors.Is(
			err,
			constant.ErrNotFound,
		) {
			return "", constant.ErrBadInput
		}
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(savedUser.Password),
		[]byte(user.Password),
	); err != nil {
		return "", constant.ErrBadInput
	}

	token, err := generateJwtToken(
		s.secret,
		savedUser,
	)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *UserService) RegisterUser(
	ctx context.Context,
	user model.User,
) (string, error) {
	foundUsername, err := s.userRepository.FindAdminUsername(
		ctx,
		user.Username,
	)
	if err != nil {
		if !errors.Is(
			err,
			constant.ErrNotFound,
		) {
			return "", err
		}
	}

	if foundUsername == user.Username {
		return "", constant.ErrConflict
	}

	userId, err := uuid.NewV7()
	if err != nil {
		return "", err
	}

	generatedHashBytes, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password),
		s.salt,
	)
	if err != nil {
		return "", err
	}

	user.ID = userId
	user.Password = string(
		generatedHashBytes,
	)
	savedUser, err := s.userRepository.CreateUser(
		ctx,
		user,
	)
	if err != nil {
		return "", err
	}

	token, err := generateJwtToken(
		s.secret,
		savedUser,
	)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *UserService) LoginUser(
	ctx context.Context,
	user model.User,
) (string, error) {
	savedUser, err := s.userRepository.FindUserByUsername(
		ctx,
		user,
	)
	if err != nil {
		if errors.Is(
			err,
			constant.ErrNotFound,
		) {
			return "", constant.ErrBadInput
		}
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(savedUser.Password),
		[]byte(user.Password),
	); err != nil {
		return "", constant.ErrBadInput
	}

	token, err := generateJwtToken(
		s.secret,
		savedUser,
	)
	if err != nil {
		return "", err
	}

	return token, nil
}

func generateJwtToken(
	secret string,
	user model.User,
) (string, error) {
	token := jwt.New(
		jwt.SigningMethodHS256,
	)

	claims := token.Claims.(jwt.MapClaims)
	userID := base64.RawStdEncoding.EncodeToString(
		[]byte(user.ID.String()),
	)
	email := base64.RawStdEncoding.EncodeToString(
		[]byte(user.Email),
	)
	username := base64.RawStdEncoding.EncodeToString(
		[]byte(user.Username),
	)
	claims["si"] = userID
	claims["mm"] = email
	claims["rn"] = username
	claims["exp"] = time.Now().
		Add(time.Hour * 72).
		Unix()

	t, err := token.SignedString(
		[]byte(secret),
	)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return t, nil
}

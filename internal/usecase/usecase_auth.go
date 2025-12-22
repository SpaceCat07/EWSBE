package usecase

import (
	"EWSBE/internal/entity"
	"EWSBE/internal/repository"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	userRepo repository.UserRepository
}

func NewAuthUsecase(userRepo repository.UserRepository) *AuthUsecase {
	return &AuthUsecase{userRepo: userRepo}
}

func (uc *AuthUsecase) Register(username, password string) error {
	if strings.TrimSpace(username) == "" || strings.TrimSpace(password) == "" {
		return errors.New("username and password cannot be empty")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &entity.User{
		Username: username,
		Password: string(hashedPassword),
	}

	return uc.userRepo.CreateUser(user)
}

func (uc *AuthUsecase) Login(username, password string) (*entity.User, error) {
	user, err := uc.userRepo.GetUserByUsername(username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid username or password")
	}

	return user, nil
}

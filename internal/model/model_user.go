package model

import (
	"EWSBE/internal/entity"
	"EWSBE/internal/repository"

	"gorm.io/gorm"
)

type userModel struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) repository.UserRepository {
	return &userModel{db: db}
}

func (r *userModel) CreateUser(user *entity.User) error {
	return r.db.Create(user).Error
}

func (r *userModel) GetUserByUsername(username string) (*entity.User, error) {
	var user entity.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

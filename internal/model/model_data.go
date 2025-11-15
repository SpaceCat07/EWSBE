package model

import (
	"EWSBE/internal/entity"
	"EWSBE/internal/repository"

	"gorm.io/gorm"
)

type dataModel struct {
	db *gorm.DB
}

func NewDataRepo(db *gorm.DB) repository.DataRepository {
	return &dataModel{db: db}
}

func (r *dataModel) CreateData(u *entity.SensorData) error {
	return r.db.Create(u).Error
}

func (r *dataModel) GetAllData() ([]entity.SensorData, error) {
	var u []entity.SensorData
	if err := r.db.Find(&u).Error; err != nil {
		return nil, err
	}

	return u, nil
}

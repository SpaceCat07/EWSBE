package repository

import "EWSBE/internal/entity"

type DataRepository interface {
	CreateData(u *entity.SensorData) error
	GetAllData() ([]entity.SensorData, error)
}

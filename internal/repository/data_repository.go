package repository

import (
	"EWSBE/internal/entity"
	"time"
)

type DataRepository interface {
	CreateData(u *entity.SensorData) error
	GetAllData() ([]entity.SensorData, error)
	GetLatestData() (*entity.SensorData, error)
	GetDataByTimeRange(start, end time.Time) ([]entity.SensorData, error)
	GetDataByLimit(limit int) ([]entity.SensorData, error)
	GetAggregatedData(interval string, start, end time.Time) ([]entity.AggregatedData, error)
	GetDataInsights() (*entity.DataInsights, error)
}

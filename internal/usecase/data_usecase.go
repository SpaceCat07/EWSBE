package usecase

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

type DataUsecase struct {
	repo DataRepository
}

func NewDataUsecase(r DataRepository) *DataUsecase {
	return &DataUsecase{repo: r}
}

func (uc *DataUsecase) Create(u *entity.SensorData) error {
	return uc.repo.CreateData(u)
}

func (uc *DataUsecase) GetAllData() ([]entity.SensorData, error) {
	return uc.repo.GetAllData()
}

func (uc *DataUsecase) GetLatestData() (*entity.SensorData, error) {
	return uc.repo.GetLatestData()
}

func (uc *DataUsecase) GetDataByTimeRange(start, end time.Time) ([]entity.SensorData, error) {
	return uc.repo.GetDataByTimeRange(start, end)
}

func (uc *DataUsecase) GetDataByLimit(limit int) ([]entity.SensorData, error) {
	return uc.repo.GetDataByLimit(limit)
}

func (uc *DataUsecase) GetAggregatedData(interval string, start, end time.Time) ([]entity.AggregatedData, error) {
	return uc.repo.GetAggregatedData(interval, start, end)
}

func (uc *DataUsecase) GetDataInsights() (*entity.DataInsights, error) {
	return uc.repo.GetDataInsights()
}

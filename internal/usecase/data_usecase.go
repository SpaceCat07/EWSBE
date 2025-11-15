package usecase

import "EWSBE/internal/entity"

type DataRepository interface {
	CreateData(u *entity.SensorData) error
	GetAllData() ([]entity.SensorData, error)
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

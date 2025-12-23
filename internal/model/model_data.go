package model

import (
	"EWSBE/internal/entity"
	"EWSBE/internal/repository"
	"fmt"
	"time"

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
	if err := r.db.Order("timestamp desc").Find(&u).Error; err != nil {
		return nil, err
	}

	return u, nil
}

func (r *dataModel) GetLatestData() (*entity.SensorData, error) {
	var data entity.SensorData
	if err := r.db.Order("timestamp desc").First(&data).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &data, nil
}

func (r *dataModel) GetDataByTimeRange(start, end time.Time) ([]entity.SensorData, error) {
	var data []entity.SensorData
	if err := r.db.Where("timestamp >= ? AND timestamp <= ?", start, end).
		Order("timestamp desc").
		Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (r *dataModel) GetDataByLimit(limit int) ([]entity.SensorData, error) {
	var data []entity.SensorData
	if err := r.db.Order("timestamp desc").Limit(limit).Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// ======================================================== //
// =================== For Insight Page =================== //
// ======================================================== //

func (r *dataModel) GetAggregatedData(interval string, start, end time.Time) ([]entity.AggregatedData, error) {
	var results []entity.AggregatedData

	// map interval to PostgreSQL date_trunc parameter
	var period string
	switch interval {
	case "hourly":
		period = "hour"
	case "daily":
		period = "day"
	case "weekly":
		period = "week"
	case "monthly":
		period = "month"
	default:
		period = "hour"
	}

	querySelect := fmt.Sprintf("date_trunc('%s', timestamp) as timestamp, AVG(temperature) as temperature, AVG(humidity) as humidity, AVG(pressure) as pressure, AVG(wind_speed) as wind_speed, AVG(rainfall) as rainfall, AVG(co2) as co2, AVG(altitude) as altitude, AVG(wind_direction) as wind_direction", period)
	queryGroup := fmt.Sprintf("date_trunc('%s', timestamp)", period)

	err := r.db.Model(&entity.SensorData{}).
		Select(querySelect).
		Where("timestamp >= ? AND timestamp <= ?", start, end).
		Group(queryGroup).
		Order("timestamp desc").
		Scan(&results).Error

	return results, err
}

func (r *dataModel) GetDataInsights() (*entity.DataInsights, error) {
	var insights entity.DataInsights
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	startOfPrevMonth := startOfMonth.AddDate(0, -1, 0)
	endOfPrevMonth := startOfMonth.Add(-time.Nanosecond)

	// Current Month Stats
	type Stats struct {
		MinTemp float64
		MaxTemp float64
		AvgTemp float64
		MinHum  float64
		MaxHum  float64
		AvgHum  float64
	}
	var currentStats Stats

	err := r.db.Model(&entity.SensorData{}).
		Select("MIN(temperature) as min_temp, MAX(temperature) as max_temp, AVG(temperature) as avg_temp, MIN(humidity) as min_hum, MAX(humidity) as max_hum, AVG(humidity) as avg_hum").
		Where("timestamp >= ?", startOfMonth).
		Scan(&currentStats).Error

	if err != nil {
		return nil, err
	}

	insights.MinTemp = currentStats.MinTemp
	insights.MaxTemp = currentStats.MaxTemp
	insights.AvgTemp = currentStats.AvgTemp
	insights.MinHum = currentStats.MinHum
	insights.MaxHum = currentStats.MaxHum
	insights.AvgHum = currentStats.AvgHum

	// previous month avg temp for comparison
	var prevMonthAvgTemp float64
	err = r.db.Model(&entity.SensorData{}).
		Select("COALESCE(AVG(temperature), 0)").
		Where("timestamp >= ? AND timestamp <= ?", startOfPrevMonth, endOfPrevMonth).
		Scan(&prevMonthAvgTemp).Error
	
	if err == nil {
		insights.PrevMonthDiff = insights.AvgTemp - prevMonthAvgTemp
	}

	type PeakHourStat struct {
		Hour       int
		AvgTemp    float64
	}
	var peakHour PeakHourStat
	
	err = r.db.Model(&entity.SensorData{}).
		Select("EXTRACT(HOUR FROM timestamp) as hour, AVG(temperature) as avg_temp").
		Where("timestamp >= ?", startOfMonth).
		Group("EXTRACT(HOUR FROM timestamp)").
		Order("avg_temp DESC").
		Limit(1).
		Scan(&peakHour).Error

	if err == nil {
		insights.PeakHour = peakHour.Hour
		insights.PeakHourAvg = peakHour.AvgTemp
	}

	return &insights, nil
}

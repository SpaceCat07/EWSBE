package entity

import (
	"time"
)

// for database storage of sensor readings
type SensorData struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	Timestamp     time.Time `json:"timestamp" gorm:"index"` // when sensor reading was taken
	Temperature   float64   `json:"temperature"`            // °C (from suhu)
	Humidity      float64   `json:"humidity"`               // % (from lembap)
	Pressure      float64   `json:"pressure"`               // hPa (from tekanan)
	Altitude      float64   `json:"altitude"`               // meters (from ketinggian)
	Co2           float64   `json:"co2"`                    // ppm (from co2)
	Distance      float64   `json:"distance"`               // cm (from jarak)
	WindSpeed     float64   `json:"windSpeed"`              // m/s (from angin)
	WindDirection float64   `json:"windDirection"`          // degrees 0-360 (from arahAngin)
	Rainfall      float64   `json:"rainfall"`               // mm (from rain)
	Voltage       float64   `json:"voltage"`                // V (from voltSensor)
	BusVoltage    float64   `json:"busVoltage"`             // V (from busVoltage)
	Current       float64   `json:"current"`                // mA (from current_mA)
}

// grouped/aggregated data for insights
type AggregatedData struct {
	Timestamp     time.Time `json:"timestamp"`
	Temperature   float64   `json:"temperature"`
	Humidity      float64   `json:"humidity"`
	Pressure      float64   `json:"pressure"`
	WindSpeed     float64   `json:"windSpeed"`
	Rainfall      float64   `json:"rainfall"`
	Co2           float64   `json:"co2"`
	Altitude      float64   `json:"altitude"`
	WindDirection float64   `json:"windDirection"`
}

// for analytical insights
type DataInsights struct {
	MinTemp        float64 `json:"minTemp"`
	MaxTemp        float64 `json:"maxTemp"`
	AvgTemp        float64 `json:"avgTemp"`
	MinHum         float64 `json:"minHum"`
	MaxHum         float64 `json:"maxHum"`
	AvgHum         float64 `json:"avgHum"`
	PrevMonthDiff  float64 `json:"prevMonthDiff"` // difference in average temp vs previous month
	PeakHour       int     `json:"peakHour"`      // hour of day with highest average temp
	PeakHourAvg    float64 `json:"peakHourAvg"`   // avg temp at peak hour
}

// from sensor MQTT payload
type MQTTSensorPayload struct {
	Waktu      int64   `json:"waktu"`       // unix timestamp milliseconds
	Suhu       float64 `json:"suhu"`        // temperature
	Lembap     float64 `json:"lembap"`      // humidity
	Tekanan    float64 `json:"tekanan"`     // pressure
	Ketinggian float64 `json:"ketinggian"`  // altitude
	Co2        float64 `json:"co2"`         // co2
	Jarak      float64 `json:"jarak"`       // distance
	Angin      float64 `json:"angin"`       // wind speed
	ArahAngin  float64 `json:"arahAngin"`   // wind direction
	BusVoltage float64 `json:"busVoltage"`  // bus voltage
	CurrentMA  float64 `json:"current_mA"`  // current
	VoltSensor float64 `json:"voltSensor"`  // voltage sensor
	Rain       float64 `json:"rain"`        // rainfall
}

func (m *MQTTSensorPayload) ToSensorData() *SensorData {
	timestamp := time.Unix(0, m.Waktu*int64(time.Millisecond))

	return &SensorData{
		Timestamp:     timestamp,
		Temperature:   m.Suhu,
		Humidity:      m.Lembap,
		Pressure:      m.Tekanan,
		Altitude:      m.Ketinggian,
		Co2:           m.Co2,
		Distance:      m.Jarak,
		WindSpeed:     m.Angin,
		WindDirection: m.ArahAngin,
		Rainfall:      m.Rain,
		Voltage:       m.VoltSensor,
		BusVoltage:    m.BusVoltage,
		Current:       m.CurrentMA,
	}
}

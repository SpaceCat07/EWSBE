package entity

import (
	"time"
)

type SensorData struct {
	ID          uint      // through Model
	CreatedAt   time.Time // through Model
	UpdatedAt   time.Time // through Model
	Temperature float64
	Humidity    float64
	Pressure    float64
	WaterLevel  float64
	Co2         float64
	WindSpeed   float64
	Rain        float64
	WindVane    int
}

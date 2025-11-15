package db

import (
	"EWSBE/internal/config"
	"errors"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(cfg config.Config) (*gorm.DB, error) {
	switch cfg.DBDriver {
    case "postgres":
        return gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
    default:
        return nil, errors.New("unsupported db driver: " + cfg.DBDriver)
    }
}
package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBDriver string
	DSN		 string
	Port	 string
}

func LoadConfig() Config {
	if err := godotenv.Load(); err != nil{
		log.Fatalf("Error loading .env file: %s", err)
	}
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	c := Config{
		DBDriver: os.Getenv("DBDRIVER"),
		DSN: fmt.Sprintf(dsn),
		Port: port,
	}

	return c
}
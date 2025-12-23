package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBDriver string
	DSN      string
	Port     string
}

func LoadConfig() Config {
	if err := godotenv.Load("../.env"); err != nil {
		log.Printf("Error loading .env file: %s", err)
	}
	
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	serverPort := os.Getenv("PORT")

	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "5432"
	}
	if serverPort == "" {
		serverPort = "8080"
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	c := Config{
		DBDriver: os.Getenv("DBDRIVER"),
		DSN:      dsn,
		Port:     serverPort,
	}

	return c
}

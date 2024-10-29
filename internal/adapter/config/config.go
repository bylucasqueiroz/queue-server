package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	ConString string
}

func NewConfig() *Config {
	return &Config{
		ConString: loadConString(),
	}
}

func loadConString() string {
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB_NAME")

	port, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if err != nil {
		panic("error to load postgres port")
	}

	conString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	return conString
}

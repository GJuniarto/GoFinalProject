package config

import "fmt"

type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
}

func NewConfig() *Config {
	return &Config{
		DBHost:     "localhost",
		DBPort:     5432,
		DBUser:     "postgres",
		DBPassword: "Qwerty90@!",
		DBName:     "GoFinal",
	}
}

func (congfig *Config) GetPostgresConfig() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		congfig.DBHost, congfig.DBPort, congfig.DBUser, congfig.DBPassword, congfig.DBName)
}

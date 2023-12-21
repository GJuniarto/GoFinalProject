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
		DBHost:     "db.ihcxbizvpajujndbtmhq.supabase.co",
		DBPort:     5432,
		DBUser:     "postgres",
		DBPassword: "Y5R4EuEHKB5j7dnh",
		DBName:     "postgres",
	}
}

//user=postgres password=[YOUR-PASSWORD] host=db.ihcxbizvpajujndbtmhq.supabase.co port=5432 dbname=postgres
func (congfig *Config) GetPostgresConfig() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		congfig.DBHost, congfig.DBPort, congfig.DBUser, congfig.DBPassword, congfig.DBName)
}

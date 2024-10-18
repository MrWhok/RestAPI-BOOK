package storage

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host     string
	Port     string
	Password string
	User     string
	DBName   string
	SSLMode  string //what is sslmode? it will encrypt the data
}

func NewConnection(config *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf( //what is dsn in here? data source name, it will contain the information to connect to the database
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{}) //what this code do? open connection to the databse using the dsn and gorm config
	if err != nil {
		return nil, err
	}

	return db, nil
}

package database

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func Connect() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=toor dbname=c3c2 sslmode=disable"
	db, err := gorm.Open("postgres", dsn)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %v", err)
	}

	return db, nil
}

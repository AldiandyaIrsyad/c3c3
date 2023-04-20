package main

import (
	"log"

	"github.com/aldiandyaIrsyad/c3c2/database"
	"github.com/aldiandyaIrsyad/c3c2/models"
	"github.com/jinzhu/gorm"
)

func main() {
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.AutoMigrate(&models.User{}, &models.Product{})

	seedUsers(db)
}

func seedUsers(db *gorm.DB) {
	users := []models.User{
		{Username: "user1", Password: "password1", Role: models.Creator},
		{Username: "user2", Password: "password2", Role: models.Creator},
		{Username: "admin", Password: "adminpassword", Role: models.Admin},
	}

	for _, user := range users {
		if err := db.Create(&user).Error; err != nil {
			log.Printf("failed to seed user '%s': %v", user.Username, err)
		}
	}
}

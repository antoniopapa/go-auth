package database

import (
	"github.com/antoniopapa/go-auth/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	connection, err := gorm.Open(mysql.Open("root:@/rest-test"), &gorm.Config{})

	if err != nil {
		panic("Could not connect to the database!")
	}

	DB = connection

	connection.AutoMigrate(&models.User{}, &models.PasswordReset{})
}

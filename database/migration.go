package db

import (
	"github.com/showbaba/go-auth-service/models"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(
		&models.User{},
		&models.Token{},
	)
}
package repository

import (
	"errors"

	"github.com/showbaba/go-auth-service/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	database *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		database: db,
	}
}

func (a *UserRepository) Create(data *models.User) (*models.User, error) {
	if err := a.database.Create(data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (a *UserRepository) Update(id uint, updates models.User) error {
	return a.database.Model(&models.User{}).Where("id = ?", id).Updates(updates).Error
}

func (a *UserRepository) Fetch(q models.User) (*models.User, bool, error) {
	var user models.User
	if err := a.database.Where(q).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return &user, true, nil
}

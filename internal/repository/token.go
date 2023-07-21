package repository

import (
	"errors"

	"github.com/showbaba/go-auth-service/models"
	"gorm.io/gorm"
)

type TokenRepository struct {
	database *gorm.DB
}

func NewTokenRepository(db *gorm.DB) *TokenRepository {
	return &TokenRepository{
		database: db,
	}
}

func (t *TokenRepository) Create(data *models.Token) (*models.Token, error) {
	if err := t.database.Create(data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (t *TokenRepository) Fetch(q models.Token) (*models.Token, bool, error) {
	var token models.Token
	if err := t.database.Where(q).First(&token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return &token, true, nil
}

func (t *TokenRepository) Delete(condition *models.Token) error {
	result := t.database.Where(&condition).Delete(models.Token{})
	return result.Error
}
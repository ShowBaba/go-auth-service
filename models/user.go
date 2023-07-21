package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID             uint `gorm:"primaryKey"`
	Email          string
	FirstName      string
	LastName       string
	Password       string
	PhoneNumber    string
	Username       string
	ProfilePicture string
	IsVerified     *bool     `gorm:"default:false"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (u *User) Insert(db *gorm.DB) (uint, error) {
	if err := db.Create(u).Error; err != nil {
		return 0, err
	}
	return u.ID, nil
}

func (u *User) GetUser(db *gorm.DB, q User) (*User, error) {
	var user User
	err := db.Where(q).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (u *User) GetByID(idb *gorm.DB, d int) (*User, error) {
	return nil, nil
}

type UserWithoutPassword struct {
	ID             uint `gorm:"primaryKey"`
	Email          string
	FirstName      string
	LastName       string
	PhoneNumber    string
	Username       string
	ProfilePicture string
	IsVerified     *bool     `gorm:"default:false"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}


func RemovePasswordFromUser(u *User) *UserWithoutPassword {
	userWithoutPassword := UserWithoutPassword{
		ID:             u.ID,
		Email:          u.Email,
		FirstName:      u.FirstName,
		LastName:       u.LastName,
		PhoneNumber:    u.PhoneNumber,
		Username:       u.Username,
		ProfilePicture: u.ProfilePicture,
		IsVerified:     u.IsVerified,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
	}

	return &userWithoutPassword
}

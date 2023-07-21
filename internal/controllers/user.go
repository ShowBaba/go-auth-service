package controllers

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/showbaba/go-auth-service/internal/repository"
	"github.com/showbaba/go-auth-service/models"
	"gorm.io/gorm"
)

type UserController struct {
	database       *gorm.DB
	userRepository *repository.UserRepository
}

func NewUserController(db *gorm.DB, userRepository *repository.UserRepository) *UserController {
	return &UserController{
		database:       db,
		userRepository: userRepository,
	}
}

func (u *UserController) FecthProfile(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")
	email := c.GetRespHeader("email")

	user, exist, err := u.userRepository.Fetch(models.User{Email: email})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": fmt.Sprint(err)})
	}

	if !exist {
		if err != nil {
			return c.Status(http.StatusNotFound).JSON(
				&fiber.Map{"message": "user with email not found"})
		}
	}
	userWithoutPassword := models.RemovePasswordFromUser(user)
	c.Status(200)
	return c.JSON(&fiber.Map{
		"success": true,
		"message": "fetch profile successful",
		"user":    userWithoutPassword,
	})
}

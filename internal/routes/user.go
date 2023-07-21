package routes

import (
	fiber "github.com/gofiber/fiber/v2"
	"github.com/showbaba/go-auth-service/internal/controllers"
	"github.com/showbaba/go-auth-service/internal/middleware"
	"github.com/showbaba/go-auth-service/internal/repository"
	"gorm.io/gorm"
)

func RegisterUserRoutes(router fiber.Router, database *gorm.DB) {
	userRepository := repository.NewUserRepository(database)
	userController := controllers.NewUserController(database, userRepository)
	authMiddleware := middleware.NewAuthMiddleware()

	userRouter := router.Group("user")
	userRouter.Get("/profile", authMiddleware.ValidateAuthHeaderToken, userController.FecthProfile)
}

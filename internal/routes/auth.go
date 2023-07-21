package routes

import (
	"github.com/go-redis/redis/v8"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/rabbitmq/amqp091-go"
	"github.com/showbaba/go-auth-service/internal/controllers"
	"github.com/showbaba/go-auth-service/internal/repository"
	"github.com/showbaba/go-auth-service/internal/validators"
	"gorm.io/gorm"
)

func RegisterAuthRoutes(router fiber.Router, database *gorm.DB, redis *redis.Client, qC *amqp091.Connection) {
	userRepository := repository.NewUserRepository(database)
	authController := controllers.NewAuthController(database, redis, qC, userRepository)

	userRouter := router.Group("auth")
	userRouter.Post("signup", validators.ValidateSignup, authController.Signup)
	userRouter.Post("verify", validators.ValidateVerifyEmail, authController.VerifyEmail)
	userRouter.Post("login", validators.ValidateLogin, authController.Login)
}

package routes

import (
	fiber "github.com/gofiber/fiber/v2"
	"github.com/rabbitmq/amqp091-go"
	"github.com/showbaba/go-auth-service/internal/controllers"
	"github.com/showbaba/go-auth-service/internal/repository"
	"github.com/showbaba/go-auth-service/internal/validators"
	"gorm.io/gorm"
)

func RegisterAuthRoutes(router fiber.Router, database *gorm.DB, qC *amqp091.Connection) {
	userRepository := repository.NewUserRepository(database)
	tokenRepository := repository.NewTokenRepository(database)
	authController := controllers.NewAuthController(database, qC, tokenRepository, userRepository)

	userRouter := router.Group("auth")
	userRouter.Post("signup", validators.ValidateSignup, authController.Signup)
	userRouter.Post("verify", validators.ValidateVerifyEmail, authController.VerifyEmail)
	userRouter.Post("login", validators.ValidateLogin, authController.Login)
	userRouter.Post("otp-resend", validators.ValidateResendOTP, authController.ResendOTP)
}

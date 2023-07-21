package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rabbitmq/amqp091-go"
	"github.com/showbaba/go-auth-service/internal/routes"
	"gorm.io/gorm"
)

func welcome(c *fiber.Ctx) error {
	return c.SendString("Welcome to gp-auth-service backend API")
}

func Routes(app *fiber.App, database *gorm.DB, qC *amqp091.Connection) {
	apiURL := "/"
	router := app.Group(apiURL)

	app.Get(apiURL, welcome)
	routes.RegisterAuthRoutes(router, database, qC)
	routes.RegisterUserRoutes(router, database)
}

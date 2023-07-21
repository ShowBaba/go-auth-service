package main

import (
	"log"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/lib/pq"
	"github.com/rabbitmq/amqp091-go"
	db "github.com/showbaba/go-auth-service/database"
	"github.com/showbaba/go-auth-service/notification"
	"github.com/showbaba/go-auth-service/router"
	"github.com/showbaba/go-auth-service/utils"
)

func main() {
	qConn, err := amqp091.Dial(utils.GetConfig().RabbitmqServerURL)
	if err != nil {
		panic(err)
	}
	defer qConn.Close()
	dbClient, conn, err := db.ConnectToPgDB(
		utils.GetConfig().DbHost,
		utils.GetConfig().DbUser,
		utils.GetConfig().DbPassword,
		utils.GetConfig().DbName,
		utils.GetConfig().DbPort,
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()


	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		notification.InitNotificationQueue(qConn)
	}()

	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())
	router.Routes(app, dbClient, qConn)
	db.Migrate(dbClient)

	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404)
	})
	port := utils.GetConfig().Port
	log.Printf("starting server on port: %s", port)
	app.Listen(port)
	wg.Wait()
}

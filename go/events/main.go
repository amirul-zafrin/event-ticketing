package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/amirul-zafrin/event-ticketing/events.git/database"
	"github.com/amirul-zafrin/event-ticketing/events.git/routes"
	"github.com/amirul-zafrin/event-ticketing/events.git/services"
	"github.com/amirul-zafrin/event-ticketing/events.git/utilities"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func ping(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Ping!"})
}

func setupRoutes(app *fiber.App) {
	app.Get("/api", ping)
	app.Post("/api/event", routes.CreateEvent)
	app.Put("/api/event/:id", routes.UpdateEvent)
	app.Delete("/api/event/:id", routes.DeleteEvent)
	app.Delete("/api/event/permanent/:id", routes.PermanentDeleteEvent)
	app.Get("/api/event", routes.GetEvents)
	app.Get("/api/event/:id", routes.GetEvent)
	app.Post("/api/event/:id/seats", routes.SetSeatCategory)
}

func main() {
	database.ConnectDB()
	app := fiber.New()
	app.Use(logger.New())

	go func() {
		rmq, err := services.NewRabbitMQ("amqp://guest:guest@127.0.0.1:5672/")
		if err != nil {
			log.Printf("Error connecting to RabbitMQ: %s", err)
		}
		rmq.ConsumeUpdate("seat_requests", utilities.UpdateSeat)
	}()

	go func() {
		setupRoutes(app)
		log.Fatal(app.Listen(":3001"))
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	// Graceful shutdown
	app.Shutdown()
}

package main

import (
	"log"

	"github.com/amirul-zafrin/event-ticketing/events.git/database"
	"github.com/amirul-zafrin/event-ticketing/events.git/routes"
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
}

func main() {
	database.ConnectDB()
	app := fiber.New()
	app.Use(logger.New())

	setupRoutes(app)
	log.Fatal(app.Listen(":3001"))

}

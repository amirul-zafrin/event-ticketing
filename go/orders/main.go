package main

import (
	"log"

	"github.com/amirul-zafrin/event-ticketing/orders.git/database"
	"github.com/amirul-zafrin/event-ticketing/orders.git/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func ping(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Ping!"})
}

func setupRoutes(app *fiber.App) {
	app.Get("/api", ping)
	app.Post("/api/order", routes.CreateOrder, routes.TotalPriceService, routes.UpdateEventSeatRMQ)
	app.Put("/api/order/:id", routes.UpdateOrder, routes.UpdateEventSeatRMQ)
}
func main() {
	database.ConnectDB()

	app := fiber.New()
	app.Use(logger.New())

	setupRoutes(app)

	log.Fatal(app.Listen(":3000"))
}

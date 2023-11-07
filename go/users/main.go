package main

import (
	"log"

	"github.com/amirul-zafrin/event-ticketing/users.git/database"
	"github.com/amirul-zafrin/event-ticketing/users.git/middlewares"
	"github.com/amirul-zafrin/event-ticketing/users.git/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func welcome(c *fiber.Ctx) error {
	return c.SendString("Ping")
}

func setupRoutes(app *fiber.App) {
	app.Get("/api", welcome)
	app.Post("/api/users", routes.CreateUser)
	app.Get("/api/users", middlewares.DeserializeUser, routes.GetUsers)
	app.Get("/api/users/:id", middlewares.DeserializeUser, routes.GetUser)
	app.Put("/api/users/:id", middlewares.DeserializeUser, routes.UpdateUser)
	app.Delete("/api/users/:id", middlewares.DeserializeUser, routes.DeleteUser)
	app.Post("/api/login", routes.UserLogin)
	app.Get("/api/logout", routes.UserLogout)
	app.Get("/api/whoami", middlewares.DeserializeUser, routes.WhoAmI)
}

func main() {
	database.ConnectDB()
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowCredentials: true, // Use this when using HTTPOnly Cookie
	}))

	setupRoutes(app)

	log.Fatal(app.Listen(":3000"))
}

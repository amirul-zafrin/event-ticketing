package api

import (
	"log"

	db "github.com/amirul-zafrin/event/db/sqlc"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	store  *db.Queries
	router *fiber.App
}

func SetupRoutes(app *fiber.App) {
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Ping!"})
	})
}

func NewServer(store *db.Queries) *Server {
	server := &Server{store: store}
	router := fiber.New()

	SetupRoutes(router)

	server.router = router
	return server
}

func (server *Server) Start(address string) {
	log.Fatal(server.router.Listen(address))
}

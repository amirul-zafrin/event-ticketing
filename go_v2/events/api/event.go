package api

import (
	"database/sql"

	db "github.com/amirul-zafrin/event/db/sqlc"
	"github.com/gofiber/fiber/v2"
)

func (server *Server) CreateEvent(c *fiber.Ctx) error {
	params := db.CreateEventParams{}

	if err := c.BodyParser(params); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	event, err := server.store.CreateEvent(c.Context(), params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": event,
	})

}

func (server *Server) GetEvent(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	event, err := server.store.GetEvent(c.Context(), int64(id))
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(err)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": event,
	})
}

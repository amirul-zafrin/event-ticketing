package routes

import (
	"time"

	"github.com/amirul-zafrin/event-ticketing/events.git/database"
	"github.com/amirul-zafrin/event-ticketing/events.git/models"
	"github.com/amirul-zafrin/event-ticketing/events.git/utilities"
	"github.com/gofiber/fiber/v2"
)

func CreateEvent(c *fiber.Ctx) error {
	var event models.Events

	if err := c.BodyParser(&event); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "failed", "error": err.Error()})
	}
	if err := utilities.ValidatorStruct(event); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	var existEvent models.Events
	if err := utilities.FindEventByName(event.Name, &existEvent); err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failed", "message": err.Error()})
	}
	if err := database.Database.Db.Create(&event).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": event})
}

func UpdateEvent(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "failed",
			"message": err.Error(),
		})
	}

	var event models.Events

	if err := c.BodyParser(&event); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failed", "message": err.Error()})
	}

	var existEvent models.Events
	if err := utilities.FindEvent(id, &existEvent); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "failed",
			"message": err.Error(),
		})
	}

	if err := database.Database.Db.Updates(&event); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "failed",
			"message": err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": event,
	})
}

func DeleteEvent(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "failed",
			"message": "Please make sure id is in int",
		})
	}
	var event models.Events
	if err := utilities.FindEvent(id, &event); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "failed",
			"message": err,
		})
	}

	if err := database.Database.Db.Model(&event).Update("deleted_at", time.Now()); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "failed",
			"message": err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Event deleted",
	})
}

func PermanentDeleteEvent(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "failed",
			"message": "Please make sure id is in int",
		})
	}
	var event models.Events
	if err := utilities.FindEvent(id, &event); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "failed",
			"message": err,
		})
	}

	if err := database.Database.Db.Delete(&event); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "failed",
			"message": err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Event deleted permanently",
	})
}

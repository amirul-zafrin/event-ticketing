package routes

import (
	"github.com/amirul-zafrin/event-ticketing/users.git/database"
	"github.com/amirul-zafrin/event-ticketing/users.git/models"
	"github.com/amirul-zafrin/event-ticketing/users.git/utilities"
	"github.com/gofiber/fiber/v2"
)

func HardDeleteUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	var user models.User
	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	if err := utilities.Authorization(c); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "failed", "message": err.Error()})
	}

	if err := utilities.FindUser(id, &user); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err := database.Database.Db.Delete(&user).Error; err != nil {
		return c.Status(404).JSON(err.Error())
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "permanently remove user!"})

}

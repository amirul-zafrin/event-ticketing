package routes

import (
	"time"

	"github.com/amirul-zafrin/event-ticketing/users.git/models"
	"github.com/amirul-zafrin/event-ticketing/users.git/utilities"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type LoginCredential struct {
	Email    string
	Password string
}

func UserLogin(c *fiber.Ctx) error {

	var loginCredential LoginCredential
	if err := c.BodyParser(&loginCredential); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}

	var user models.User
	if err := utilities.FindUserByEmail(loginCredential.Email, &user); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginCredential.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid credentials")
	}
	token, err := utilities.GenerateJWT(&user)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		HTTPOnly: true,
	}

	c.Cookie(&cookie)
	responseUser := CreateResponseUser(user)
	return c.Status(200).JSON(responseUser)
}

func UserLogout(c *fiber.Ctx) error {
	expired := time.Now().Add(-time.Hour * 24)

	cookie := fiber.Cookie{
		Name:     "",
		Value:    "",
		Expires:  expired,
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}

func WhoAmI(c *fiber.Ctx) error {
	user := c.Locals("user")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"user": user}})
}

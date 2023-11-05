package routes

import (
	"github.com/amirul-zafrin/event-ticketing/users.git/database"
	"github.com/amirul-zafrin/event-ticketing/users.git/errors"
	"github.com/amirul-zafrin/event-ticketing/users.git/models"
	"github.com/amirul-zafrin/event-ticketing/users.git/utilities"
	"github.com/amirul-zafrin/event-ticketing/users.git/validators"
	"github.com/gofiber/fiber/v2"
)

type User struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
}

type UserResponse struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func CreateResponseUser(userModel models.User) UserResponse {
	return UserResponse{ID: userModel.ID, FirstName: userModel.FirstName,
		LastName: userModel.LastName, Email: userModel.Email}
}

func CreateUser(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}
	if err := validators.ValidateStruct(&user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}
	if err := utilities.FindUserByEmail(user.Email, &user); err == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(errors.ErrEmailAlreadyUsed.Error())
	} else {
		user.Password, _ = utilities.HashPassword(user.Password)
		if err := database.Database.Db.Create(&user).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}
	}
	responseUser := CreateResponseUser(user)
	return c.Status(200).JSON(responseUser)
}

func GetUsers(c *fiber.Ctx) error {
	users := []models.User{}

	database.Database.Db.Find(&users)
	responseUsers := []UserResponse{}
	for _, user := range users {
		responseUser := CreateResponseUser(user)
		responseUsers = append(responseUsers, responseUser)
	}

	return c.Status(200).JSON(responseUsers)
}

func GetUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var user models.User

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	if err := utilities.FindUser(id, &user); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	responseUser := CreateResponseUser(user)

	return c.Status(200).JSON(responseUser)
}

func UpdateUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	var user models.User

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	if err := utilities.FindUser(id, &user); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	type UpdateUser struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}

	var updateData UpdateUser

	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(500).JSON(err.Error())
	}

	user.FirstName = updateData.FirstName
	user.LastName = updateData.LastName
	// TODO: Dont update if the value nil
	user.Email = updateData.Email
	user.Password = updateData.Password

	database.Database.Db.Save(&user)

	responseUser := CreateResponseUser(user)
	return c.Status(200).JSON(responseUser)
}

func DeleteUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	var user models.User

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	if err := utilities.FindUser(id, &user); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err := database.Database.Db.Delete(&user).Error; err != nil {
		return c.Status(404).JSON(err.Error())
	}

	return c.Status(200).SendString("Delete Success!")
}

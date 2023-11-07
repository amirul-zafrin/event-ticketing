package utilities

import (
	"errors"
	"time"

	"github.com/amirul-zafrin/event-ticketing/users.git/database"
	"github.com/amirul-zafrin/event-ticketing/users.git/models"
	"github.com/gofiber/fiber/v2"
	jtoken "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func GetBasedResponseObject() map[string]interface{} {
	response := make(map[string]interface{})
	response["status"] = "fail"
	response["message"] = "something went wrong"
	return response
}

func FindUser(id int, user *models.User) error {
	database.Database.Db.Find(&user, "id = ?", id)
	if user.ID == 0 {
		return errors.New("user does not exists")
	}
	return nil
}

func FindUserByEmail(email string, user *models.User) error {
	database.Database.Db.Find(&user, "email = ?", email)
	if user.ID == 0 {
		return errors.New("user does not exists")
	}
	return nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func GenerateJWT(user *models.User) (string, error) {
	day := time.Hour * 24
	claims := jtoken.MapClaims{
		"ID":    user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(day * 1).Unix(),
	}
	const SecretKey = "secret"
	token := jtoken.NewWithClaims(jtoken.SigningMethodHS256, claims)
	return token.SignedString([]byte(SecretKey))
}

func Authorization(c *fiber.Ctx, users ...*models.User) error {
	thisUser := c.Locals("user").(models.User)
	if thisUser.Admin {
		return nil
	}

	var found bool
	for _, user := range users {
		if user.ID == thisUser.ID {
			found = true
			break
		}
	}

	if !found {
		return errors.New("user don't have access permission")
	}

	return nil
}

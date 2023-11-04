package utilities

import (
	"errors"

	"github.com/amirul-zafrin/event-ticketing/users.git/database"
	"github.com/amirul-zafrin/event-ticketing/users.git/models"
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

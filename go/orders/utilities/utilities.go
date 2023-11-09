package utilities

import (
	"errors"

	"github.com/amirul-zafrin/event-ticketing/orders.git/database"
	"github.com/amirul-zafrin/event-ticketing/orders.git/models"
)

func FindOrder(id int, order *models.Orders) error {
	database.Database.Db.Find(&order, "id = ? AND deleted_at is NULL", id)
	if order.ID == 0 {
		return errors.New("no order found")
	}
	return nil
}

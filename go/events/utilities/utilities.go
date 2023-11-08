package utilities

import (
	"errors"

	"github.com/amirul-zafrin/event-ticketing/events.git/database"
	"github.com/amirul-zafrin/event-ticketing/events.git/models"
)

func FindEvent(id int, event *models.Events) error {
	database.Database.Db.Find(&event, "id = ? AND deleted_at is NULL", id)
	if event.ID == 0 {
		return errors.New("no event found")
	}
	return nil
}

func FindEventByName(name string, event *models.Events) error {
	database.Database.Db.Find(&event, "name ILIKE ? AND deleted_at is NULL", name)
	if event.ID == 0 {
		return errors.New("no event found")
	}
	return nil
}

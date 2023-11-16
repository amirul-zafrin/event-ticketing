package utilities

import (
	"errors"
	"strconv"

	"github.com/amirul-zafrin/event-ticketing/events.git/constants"
	"github.com/amirul-zafrin/event-ticketing/events.git/database"
	"github.com/amirul-zafrin/event-ticketing/events.git/models"
	"gorm.io/datatypes"
)

func FindEvent(id int, event *models.Events) error {
	database.Database.Db.Find(&event, "id = ? AND deleted_at IS NULL", id)
	if event.ID == 0 {
		return errors.New("no event found")
	}
	return nil
}

func FindEventByName(name string, event *models.Events) error {
	database.Database.Db.Find(&event, "name ILIKE ? deleted_at is NULL", name)
	if event.ID == 0 {
		return errors.New("no event found")
	}
	return nil
}

func GenerateSeats(numStart int, numEnd int, class string) map[string]interface{} {
	seats := make(map[string]interface{})
	for i := numStart; i <= numEnd; i++ {
		seats[strconv.Itoa(i)] = map[string]string{
			"class":  class,
			"status": constants.Available,
		}
	}
	return seats
}

func MergeMap(n *datatypes.JSONMap, m *map[string]interface{}) {
	if *n == nil {
		*n = make(map[string]interface{})
	}
	for k, v := range *m {
		(*n)[k] = v
	}
}

func GetPricesByEventClass(id uint, class string) (float64, error) {
	prices := models.Prices{}
	database.Database.Db.Find(&prices, "event_id = ? AND class = ?", id, class)
	if prices.ID == 0 {
		return 0, errors.New("event ID or Class not in this table")
	}
	return prices.Price, nil
}

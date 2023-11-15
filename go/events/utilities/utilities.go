package utilities

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"

	"github.com/amirul-zafrin/event-ticketing/events.git/constants"
	"github.com/amirul-zafrin/event-ticketing/events.git/database"
	"github.com/amirul-zafrin/event-ticketing/events.git/models"
	"github.com/amirul-zafrin/event-ticketing/events.git/services"
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

func MergeMap(n, m map[string]interface{}) {
	for k, v := range m {
		n[k] = v
	}
}

func UpdateSeat(opt *services.LockingRequest) error {
	event := models.Events{}
	if err := FindEvent(opt.EventID, &event); err != nil {
		log.Printf("Failed to update seat: %s", err)
		return err
	}
	if event.Seats == nil {
		return errors.New("no seat category were set")
	}
	var resultMap map[string][]int
	jsonbytes, err := opt.Details.MarshalJSON()
	if err != nil {
		log.Printf("Error when marshalling")
		return err
	}
	err = json.Unmarshal(jsonbytes, &resultMap)
	if err != nil {
		log.Printf("Error when unmarshalling")
		return err
	}
	status := "locked"
	if opt.IsPaid {
		status = "occupied"
	}
	for _, element := range resultMap {
		for _, val := range element {
			if seatInfo, isMap := event.Seats[strconv.Itoa(val)].(map[string]interface{}); isMap {
				seatInfo["status"] = status
			}
		}
	}
	database.Database.Db.Save(&event)
	return nil
}

package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Events struct {
	gorm.Model
	Name     string            `json:"name"`
	Date     *time.Time        `json:"date"`
	Location string            `json:"location"`
	Capacity uint              `json:"capacity"`
	Seats    datatypes.JSONMap `json:"seats" example:"{1:{\"status\": \"available\", \"Class\": \"classA\" }\"}"`
	Prices   []Prices          `json:"prices" gorm:"foreignKey:EventID"`
}

type Prices struct {
	gorm.Model
	Class   string  `json:"class"`
	Price   float64 `json:"price"`
	EventID uint
}

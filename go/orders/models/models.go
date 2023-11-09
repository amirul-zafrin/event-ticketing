package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Orders struct {
	gorm.Model
	EventID     uint           `json:"event_id"`
	UserID      uint           `json:"user_id"`
	Details     datatypes.JSON `json:"details" example:"{\"classA\": [1,2,3], \"classB\": [1]}"`
	TotalAmount float32        `json:"total_amount"`
	IsPaid      bool           `json:"is_paid"`
}

// models/tombola.go
package models

import (
	"gorm.io/gorm"
)

type Ticket struct {
	gorm.Model
	UserID     uint `json:"user_id"`
	KermesseID uint `json:"kermesse_id"`
}

type Lot struct {
	gorm.Model
	Name       string `json:"name"`
	KermesseID uint   `json:"kermesse_id"`
	WinnerID   uint   `json:"winner_id"`
}

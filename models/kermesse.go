// models/kermesse.go
package models

import "gorm.io/gorm"

// Modèle Kermesse avec une relation one-to-many avec les stands
type Kermesse struct {
	gorm.Model
	Name   string  `json:"name"`
	Stands []Stand `json:"stands" gorm:"foreignKey:KermesseID"`
}

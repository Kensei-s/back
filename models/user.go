// models/user.go
package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name"`                        // Facultatif
	Email    string `json:"email" binding:"required"`    // Binding obligatoire
	Password string `json:"password" binding:"required"` // Binding obligatoire
	Role     string `json:"role"`                        // Facultatif, par exemple : "organisateur"
	Tokens   int    `json:"tokens"`                      // Nombre de jetons détenus par l'utilisateur
}

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

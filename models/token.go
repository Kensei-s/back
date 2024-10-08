// models/token.go
package models

import "gorm.io/gorm"

type TokenTransaction struct {
	gorm.Model
	UserID      uint   `json:"user_id"`     // ID de l'utilisateur effectuant la transaction
	Amount      int    `json:"amount"`      // Quantité de jetons transférés
	Description string `json:"description"` // Description de la transaction (achat, utilisation)
}

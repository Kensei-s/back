// models/stand.go
package models

import "gorm.io/gorm"

// Modèle Stand avec une clé étrangère pour la kermesse
type Stand struct {
	gorm.Model
	Name       string `json:"name"`
	Type       string `json:"type"`        // Type: nourriture, boisson, activité
	Stock      int    `json:"stock"`       // Quantité de stock pour les stands de nourriture/boisson
	Price      int    `json:"price"`       // Prix en jetons pour une unité
	KermesseID uint   `json:"kermesse_id"` // Clé étrangère vers la kermesse
}

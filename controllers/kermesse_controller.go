// controllers/kermesse_controller.go
package controllers

import (
	"net/http"

	"github.com/Kensei-s/back/config"
	"github.com/Kensei-s/back/models"

	"github.com/gin-gonic/gin"
)

// Créer une nouvelle kermesse (seulement pour les organisateurs)
func CreateKermesse(c *gin.Context) {
	var kermesse models.Kermesse
	if err := c.ShouldBindJSON(&kermesse); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&kermesse).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, kermesse)
}

// Récupérer toutes les kermesses
func GetKermesses(c *gin.Context) {
	var kermesses []models.Kermesse
	if err := config.DB.Preload("Stands").Find(&kermesses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, kermesses)
}

// Mettre à jour une kermesse (seulement pour les organisateurs)
func UpdateKermesse(c *gin.Context) {
	var kermesse models.Kermesse

	// Vérifier si la kermesse existe
	if err := config.DB.Where("id = ?", c.Param("id")).First(&kermesse).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kermesse non trouvée"})
		return
	}

	// Lier les données reçues à la kermesse
	if err := c.ShouldBindJSON(&kermesse); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Mettre à jour la kermesse
	config.DB.Save(&kermesse)

	c.JSON(http.StatusOK, kermesse)
}

// Supprimer une kermesse (seulement pour les organisateurs)
func DeleteKermesse(c *gin.Context) {
	var kermesse models.Kermesse

	// Vérifier si la kermesse existe
	if err := config.DB.Where("id = ?", c.Param("id")).First(&kermesse).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kermesse non trouvée"})
		return
	}

	// Supprimer la kermesse
	config.DB.Delete(&kermesse)

	c.JSON(http.StatusOK, gin.H{"message": "Kermesse supprimée avec succès"})
}

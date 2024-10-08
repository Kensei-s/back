// controllers/stand_controller.go
package controllers

import (
	"net/http"

	"github.com/Kensei-s/back/config"
	"github.com/Kensei-s/back/models"

	"github.com/gin-gonic/gin"
)

// Créer un nouveau stand pour une kermesse (seulement pour les teneurs de stand)
func CreateStand(c *gin.Context) {
	var stand models.Stand
	if err := c.ShouldBindJSON(&stand); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Créer le stand
	if err := config.DB.Create(&stand).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stand)
}

// Récupérer tous les stands d'une kermesse
func GetStands(c *gin.Context) {
	var stands []models.Stand

	// Récupérer tous les stands liés à une kermesse spécifique
	if err := config.DB.Where("kermesse_id = ?", c.Param("kermesse_id")).Find(&stands).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stands)
}

// Mettre à jour un stand (seulement pour les teneurs de stand)
func UpdateStand(c *gin.Context) {
	var stand models.Stand

	// Vérifier si le stand existe
	if err := config.DB.Where("id = ?", c.Param("id")).First(&stand).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stand non trouvé"})
		return
	}

	// Lier les données reçues au stand
	if err := c.ShouldBindJSON(&stand); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Mettre à jour le stand
	config.DB.Save(&stand)

	c.JSON(http.StatusOK, stand)
}

// Supprimer un stand (seulement pour les teneurs de stand)
func DeleteStand(c *gin.Context) {
	var stand models.Stand

	// Vérifier si le stand existe
	if err := config.DB.Where("id = ?", c.Param("id")).First(&stand).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stand non trouvé"})
		return
	}

	// Supprimer le stand
	config.DB.Delete(&stand)

	c.JSON(http.StatusOK, gin.H{"message": "Stand supprimé avec succès"})
}

// Utiliser des jetons pour interagir avec un stand (achat de produits)
func UseTokens(c *gin.Context) {
	var input struct {
		StandID uint `json:"stand_id"` // L'ID du stand avec lequel interagir
		Amount  int  `json:"amount"`   // Nombre d'unités à acheter (en jetons)
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Récupérer l'utilisateur connecté
	userID := c.MustGet("user_id").(uint)
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Utilisateur non trouvé"})
		return
	}

	// Récupérer le stand
	var stand models.Stand
	if err := config.DB.First(&stand, input.StandID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Stand non trouvé"})
		return
	}

	// Calcul du coût total de la transaction (prix unitaire * nombre d'unités)
	totalCost := stand.Price * input.Amount

	// Vérifier que l'utilisateur a assez de jetons
	if user.Tokens < totalCost {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Jetons insuffisants"})
		return
	}

	// Vérifier le stock du stand
	if stand.Stock < input.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Stock insuffisant"})
		return
	}

	// Déduire les jetons de l'utilisateur
	user.Tokens -= totalCost

	// Réduire le stock du stand
	stand.Stock -= input.Amount

	// Sauvegarder les modifications dans la base de données
	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour des jetons"})
		return
	}

	if err := config.DB.Save(&stand).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour du stock"})
		return
	}

	// Réponse avec les informations de la transaction réussie
	c.JSON(http.StatusOK, gin.H{
		"message":          "Transaction réussie",
		"remaining_tokens": user.Tokens,
		"remaining_stock":  stand.Stock,
	})
}

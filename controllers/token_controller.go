// controllers/token_controller.go
package controllers

import (
	"net/http"

	"github.com/Kensei-s/back/config"
	"github.com/Kensei-s/back/models"

	"github.com/gin-gonic/gin"
)

// Acheter des jetons (simulation de paiement)
func BuyTokens(c *gin.Context) {
	var input struct {
		Amount int `json:"amount"` // Nombre de jetons à acheter
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

	// Simuler l'ajout des jetons au compte de l'utilisateur
	user.Tokens += input.Amount
	config.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{"message": "Jetons ajoutés avec succès", "new_balance": user.Tokens})
}

// Distribuer des jetons à un élève (parent)
func DistributeTokens(c *gin.Context) {
	var input struct {
		ChildID uint `json:"child_id"`
		Amount  int  `json:"amount"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Enregistrer la transaction de distribution de jetons pour l'enfant
	transaction := models.TokenTransaction{
		UserID:      input.ChildID, // Enregistrer la transaction au nom de l'enfant
		Amount:      input.Amount,
		Description: "Distribution de jetons par le parent",
	}
	config.DB.Create(&transaction)

	c.JSON(http.StatusOK, gin.H{"message": "Jetons distribués avec succès à l'enfant", "child_id": input.ChildID})
}

// // Utiliser des jetons pour un stand (élève)
// func UseTokens(c *gin.Context) {
// 	var input struct {
// 		StandID uint `json:"stand_id"`
// 		Amount  int  `json:"amount"`
// 	}
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	studentID, _ := c.Get("user_id") // ID de l'élève connecté

// 	// Enregistrer la transaction d'utilisation de jetons
// 	transaction := models.TokenTransaction{
// 		UserID:      studentID.(uint),
// 		Amount:      -input.Amount, // Les jetons sont utilisés (réduction)
// 		Description: "Utilisation de jetons pour un stand",
// 	}
// 	config.DB.Create(&transaction)

// 	c.JSON(http.StatusOK, gin.H{"message": "Jetons utilisés avec succès pour le stand", "stand_id": input.StandID})
// }

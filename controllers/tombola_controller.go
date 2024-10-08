// controllers/tombola_controller.go
package controllers

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/Kensei-s/back/config"
	"github.com/Kensei-s/back/models"

	"github.com/gin-gonic/gin"
)

// Acheter des tickets de tombola avec des jetons
func BuyTicket(c *gin.Context) {
	var input struct {
		KermesseID uint `json:"kermesse_id"`
		Quantity   int  `json:"quantity"` // Nombre de tickets à acheter
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

	// Prix du ticket (par exemple, 5 jetons par ticket)
	const ticketPrice = 5
	totalCost := ticketPrice * input.Quantity

	// Vérifier que l'utilisateur a assez de jetons
	if user.Tokens < totalCost {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Jetons insuffisants"})
		return
	}

	// Déduire les jetons et créer les tickets
	user.Tokens -= totalCost
	for i := 0; i < input.Quantity; i++ {
		ticket := models.Ticket{
			UserID:     userID,
			KermesseID: input.KermesseID,
		}
		config.DB.Create(&ticket)
	}
	config.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{"message": "Tickets achetés avec succès", "remaining_tokens": user.Tokens})
}

// Créer un lot pour la tombola (organisateur)
func CreateLot(c *gin.Context) {
	var lot models.Lot
	if err := c.ShouldBindJSON(&lot); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Enregistrer le lot
	if err := config.DB.Create(&lot).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création du lot"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Lot créé avec succès"})
}

// Tirage au sort pour la tombola (organisateur)
func DrawWinner(c *gin.Context) {
	var input struct {
		KermesseID uint `json:"kermesse_id"`
		LotID      uint `json:"lot_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Vérifier que l'utilisateur est un organisateur
	role, exists := c.Get("role")
	if !exists || role != "organisateur" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission refusée"})
		return
	}

	// Récupérer tous les tickets associés à la kermesse
	var tickets []models.Ticket
	if err := config.DB.Where("kermesse_id = ?", input.KermesseID).Find(&tickets).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des tickets"})
		return
	}

	// Vérifier qu'il y a des tickets disponibles pour le tirage
	if len(tickets) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Aucun ticket disponible pour le tirage"})
		return
	}

	// Tirage aléatoire parmi les tickets disponibles
	rand.Seed(time.Now().UnixNano())
	winnerTicket := tickets[rand.Intn(len(tickets))]

	// Associer le gagnant au lot
	var lot models.Lot
	if err := config.DB.First(&lot, input.LotID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lot non trouvé"})
		return
	}

	// Vérifier si le lot a déjà un gagnant
	if lot.WinnerID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ce lot a déjà un gagnant"})
		return
	}

	lot.WinnerID = winnerTicket.UserID
	if err := config.DB.Save(&lot).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de l'enregistrement du gagnant"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Tirage effectué avec succès",
		"winner_id": lot.WinnerID,
		"user_id":   winnerTicket.UserID,
	})
}

// Récupérer tous les lots associés à une kermesse
func GetLotsByKermesse(c *gin.Context) {
	kermesseID := c.Param("kermesse_id")

	var lots []models.Lot
	if err := config.DB.Where("kermesse_id = ?", kermesseID).Find(&lots).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des lots"})
		return
	}

	// Si aucun lot n'est trouvé
	if len(lots) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Aucun lot trouvé pour cette kermesse"})
		return
	}

	c.JSON(http.StatusOK, lots)
}

// controllers/webhook_controller.go
package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/Kensei-s/back/config"
	"github.com/Kensei-s/back/models"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/webhook"
)

// Webhook Stripe pour traiter les paiements réussis
func StripeWebhook(c *gin.Context) {
	const MaxBodyBytes = int64(65536)
	req := c.Request
	req.Body = http.MaxBytesReader(c.Writer, req.Body, MaxBodyBytes)
	payload, err := ioutil.ReadAll(req.Body)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Erreur de lecture du body"})
		return
	}

	// Vérification de la signature du webhook Stripe
	endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	sigHeader := req.Header.Get("Stripe-Signature")
	event, err := webhook.ConstructEvent(payload, sigHeader, endpointSecret)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Signature webhook incorrecte"})
		return
	}

	// Traiter l'événement Stripe
	if event.Type == "checkout.session.completed" {
		var session stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Erreur de deserialization"})
			return
		}

		// Récupérer l'ID de l'utilisateur associé à la session
		userID, err := strconv.ParseUint(session.ClientReferenceID, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ClientReferenceID invalide"})
			return
		}

		// Ajouter des jetons à l'utilisateur après paiement réussi
		var user models.User
		if err := config.DB.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Utilisateur non trouvé"})
			return
		}

		// Ajouter les jetons (par exemple 50 jetons) après paiement
		amountInEuros := session.AmountTotal / 100 // Conversion en euros
		user.Tokens += int(amountInEuros)          // Ajout de jetons correspondant au montant payé
		config.DB.Save(&user)

		fmt.Printf("Paiement réussi, jetons ajoutés à l'utilisateur %s\n", user.Email)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Webhook reçu avec succès"})
}

// controllers/payment_controller.go
package controllers

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/checkout/session"
)

// Initier un paiement Stripe pour l'achat de jetons
func CreateCheckoutSession(c *gin.Context) {
	var input struct {
		Amount int `json:"amount"` // Nombre de jetons à acheter
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Récupérer l'utilisateur connecté
	userID := c.MustGet("user_id").(uint)

	// Initialiser Stripe avec la clé secrète
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	// Créer une session de paiement Stripe avec référence à l'utilisateur
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}), // Paiement par carte
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("eur"), // Devise
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("Achat de jetons"),
					},
					UnitAmount: stripe.Int64(int64(input.Amount) * 100), // Prix en centimes
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode:              stripe.String("payment"),
		SuccessURL:        stripe.String("http://localhost:3000/success"),        // URL de succès
		CancelURL:         stripe.String("http://localhost:3000/cancel"),         // URL d'annulation
		ClientReferenceID: stripe.String(strconv.FormatUint(uint64(userID), 10)), // ID de l'utilisateur en référence
	}

	s, err := session.New(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Retourner l'URL de paiement Stripe
	c.JSON(http.StatusOK, gin.H{
		"checkout_url": s.URL,
	})
}

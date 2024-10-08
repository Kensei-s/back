package routes

import (
	"github.com/Kensei-s/back/controllers"
	"github.com/Kensei-s/back/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Routes d'authentification
	r.POST("/login", controllers.Login)
	r.POST("/register", controllers.Register)

	// Route pour la déconnexion
	r.POST("/logout", middlewares.AuthMiddleware(), controllers.Logout)

	// Route pour récupérer les informations de l'utilisateur connecté
	r.GET("/users/me", middlewares.AuthMiddleware(), controllers.GetCurrentUser) // Ajout de la route pour récupérer l'utilisateur connecté

	// Routes protégées par authentification et rôle
	r.GET("/users", middlewares.AuthMiddleware(), controllers.GetUsers)

	// Routes pour les kermesses (CRUD complet, réservé aux organisateurs)
	r.POST("/kermesses", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("organisateur"), controllers.CreateKermesse)
	r.GET("/kermesses", middlewares.AuthMiddleware(), controllers.GetKermesses)
	r.PUT("/kermesses/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("organisateur"), controllers.UpdateKermesse)
	r.DELETE("/kermesses/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("organisateur"), controllers.DeleteKermesse)

	// Routes pour les stands (CRUD complet, réservé aux teneurs de stand)
	r.POST("/stands", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("teneur_de_stand"), controllers.CreateStand)
	r.GET("/stands/:kermesse_id", middlewares.AuthMiddleware(), controllers.GetStands)
	r.PUT("/stands/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("teneur_de_stand"), controllers.UpdateStand)
	r.DELETE("/stands/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("teneur_de_stand"), controllers.DeleteStand)

	// Routes pour la gestion des rôles par les administrateurs
	r.PUT("/admin/change-role/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("administrateur"), controllers.AdminChangeRole)

	// Routes pour les jetons
	r.POST("/buy-tokens", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("parent"), controllers.BuyTokens)
	r.POST("/distribute-tokens", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("parent"), controllers.DistributeTokens)
	r.POST("/use-tokens", middlewares.AuthMiddleware(), middlewares.RolesMiddleware("parent", "eleve"), controllers.UseTokens)

	// Les parents et élèves peuvent acheter des tickets de la tombola
	r.POST("/buy-ticket", middlewares.AuthMiddleware(), controllers.BuyTicket)

	// Seuls les organisateurs peuvent gérer la tombola
	r.POST("/create-lot", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("organisateur"), controllers.CreateLot)
	r.POST("/draw-winner", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("organisateur"), controllers.DrawWinner)

	// Route pour récupérer les lots d'une kermesse
	r.GET("/kermesses/:kermesse_id/lots", controllers.GetLotsByKermesse)

	// Route pour Stripe Checkout et webhook
	r.POST("/create-checkout-session", middlewares.AuthMiddleware(), controllers.CreateCheckoutSession)
	r.POST("/webhook", controllers.StripeWebhook)
}

// main.go
package main

import (
	"log"
	"os"

	"github.com/Kensei-s/back/config"
	"github.com/Kensei-s/back/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Charger les variables d'environnement depuis le fichier .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Erreur lors du chargement du fichier .env : %v", err)
	}
	// Initialisation de la base de données
	config.ConnectDatabase()

	// Initialisation du routeur
	r := gin.Default()

	// Charger les routes
	routes.SetupRoutes(r)

	// Démarrer le serveur
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(r.Run(":" + port))
}

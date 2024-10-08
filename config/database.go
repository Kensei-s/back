// config/database.go
package config

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/Kensei-s/back/models" // Replace "your-package-path" with the actual package path where the models are defined
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "host=localhost user=postgres password=yourpassword dbname=kermesse port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Erreur lors de la connexion à la base de données: ", err)
	}

	// Migrer les modèles
	// config/database.go
	db.AutoMigrate(&models.User{}, &models.Kermesse{}, &models.Stand{}, &models.Ticket{}, &models.Lot{}, &models.BlacklistedToken{})

	DB = db
}

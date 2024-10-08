// controllers/auth_controller.go
package controllers

import (
	"net/http"
	"strings" // Add the import statement for the strings package
	"time"

	"github.com/Kensei-s/back/config"
	"github.com/Kensei-s/back/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Liste des rôles acceptés pour les utilisateurs normaux
var validRoles = []string{"teneur_de_stand", "parent", "eleve"}

// Liste des rôles restreints qui ne peuvent être attribués que par un administrateur
var restrictedRoles = []string{"organisateur", "administrateur"}

// Register - Enregistrement d'un nouvel utilisateur avec gestion des rôles restreints
func Register(c *gin.Context) {
	var user models.User

	// Lier les données reçues au modèle User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Vérifier si le rôle spécifié est restreint
	if isRestrictedRole(user.Role) {
		// On vérifie que l'utilisateur connecté a le droit d'attribuer ce rôle
		loggedInUserRole, exists := c.Get("role")
		if !exists || loggedInUserRole != "administrateur" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Vous n'avez pas la permission d'attribuer ce rôle"})
			return
		}
	} else {
		// Si le rôle n'est pas restreint, on vérifie qu'il est valide
		if !isValidRole(user.Role) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Rôle invalide. Rôles acceptés : teneur_de_stand, parent, eleve"})
			return
		}
	}

	// Hacher le mot de passe avant de le stocker
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du hachage du mot de passe"})
		return
	}
	user.Password = string(hashedPassword)

	// Enregistrer l'utilisateur dans la base de données
	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de l'enregistrement de l'utilisateur"})
		return
	}

	// Retourner un message de succès
	c.JSON(http.StatusOK, gin.H{"message": "Utilisateur enregistré avec succès", "role": user.Role})
}

// Fonction pour vérifier si le rôle fait partie des rôles restreints
func isRestrictedRole(role string) bool {
	for _, restrictedRole := range restrictedRoles {
		if role == restrictedRole {
			return true
		}
	}
	return false
}

// Fonction pour vérifier si le rôle est valide parmi les rôles non restreints
func isValidRole(role string) bool {
	for _, validRole := range validRoles {
		if role == validRole {
			return true
		}
	}
	return false
}

// Login - Connexion utilisateur
func Login(c *gin.Context) {
	var input models.LoginInput
	var user models.User

	// Lier les données reçues au modèle LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Les données fournies sont incorrectes"})
		return
	}

	// Vérifier que l'utilisateur existe dans la base de données
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email ou mot de passe incorrect"})
		return
	}

	// Comparer le mot de passe hashé
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email ou mot de passe incorrect"})
		return
	}

	// Générer un token JWT si les informations sont valides
	token, err := generateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de générer le token"})
		return
	}

	// Retourner le token JWT
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Logout - Invalider un token JWT en l'ajoutant à la liste noire
func Logout(c *gin.Context) {
	// Récupérer le token depuis le header Authorization
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Aucun token fourni"})
		return
	}

	// Vérifier que le token commence bien par "Bearer " et le retirer
	if !strings.HasPrefix(tokenString, "Bearer ") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format de token invalide"})
		return
	}
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// Ajouter le token à la liste noire
	blacklistedToken := models.BlacklistedToken{Token: tokenString}
	if err := config.DB.Create(&blacklistedToken).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible d'invalider le token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Déconnexion réussie, le token a été invalidé"})
}

// GetUsers - Récupère tous les utilisateurs
func GetUsers(c *gin.Context) {
	var users []models.User

	// Récupérer tous les utilisateurs dans la base de données
	if err := config.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des utilisateurs"})
		return
	}

	// Retourner les utilisateurs sous forme de JSON
	c.JSON(http.StatusOK, users)
}

// GetCurrentUser - Récupère les détails de l'utilisateur connecté
func GetCurrentUser(c *gin.Context) {
	var user models.User

	// Récupérer l'ID de l'utilisateur à partir du token JWT
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non authentifié"})
		return
	}

	// Rechercher l'utilisateur dans la base de données en fonction de l'ID
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Utilisateur non trouvé"})
		return
	}

	// Retourner les détails de l'utilisateur
	c.JSON(http.StatusOK, user)
}

// Fonction pour générer un token JWT
func generateToken(user models.User) (string, error) {
	// Créer un token JWT avec les informations utilisateur et rôle
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,                             // Inclure le rôle de l'utilisateur dans le token
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // Expire dans 72 heures
	})

	// Signer le token avec une clé secrète
	return token.SignedString([]byte("secret"))
}

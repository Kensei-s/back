// middlewares/auth_middleware.go

package middlewares

import (
	"net/http"
	"strings"

	"github.com/Kensei-s/back/config" // Import de la configuration de la base de données
	"github.com/Kensei-s/back/models" // Import des modèles
	"github.com/dgrijalva/jwt-go"     // Import du package JWT pour la gestion des tokens
	"github.com/gin-gonic/gin"        // Import de Gin
)

// AuthMiddleware - Vérifie la validité du token et si le token est blacklisté
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Récupération du token depuis l'en-tête Authorization
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header manquant"})
			c.Abort()
			return
		}

		// Vérifier que le token commence par "Bearer " puis le retirer
		if !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format de token invalide"})
			c.Abort()
			return
		}
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Vérifier si le token est blacklisté
		if isTokenBlacklisted(tokenString) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token invalide ou expiré"})
			c.Abort()
			return
		}

		// Validation du token JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil // Clé secrète utilisée pour signer les tokens
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token invalide"})
			c.Abort()
			return
		}

		// Extraction des claims du token
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("user_id", uint(claims["user_id"].(float64))) // Récupération de l'ID utilisateur
			c.Set("role", claims["role"].(string))              // Récupération du rôle utilisateur
		}

		c.Next() // Continuer si le token est valide
	}
}

// Vérifie si un token est blacklisté
func isTokenBlacklisted(tokenString string) bool {
	var blacklistedToken models.BlacklistedToken
	err := config.DB.Where("token = ?", tokenString).First(&blacklistedToken).Error
	return err == nil // Retourne vrai si le token est trouvé dans la liste noire
}

// RoleMiddleware - Vérifie si l'utilisateur a le rôle requis pour accéder à la route
func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists || userRole != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Accès refusé : rôle insuffisant"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// RolesMiddleware - Vérifie si l'utilisateur a un des rôles autorisés
func RolesMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Rôle non trouvé"})
			c.Abort()
			return
		}

		// Vérifie si le rôle utilisateur est dans la liste des rôles autorisés
		for _, role := range allowedRoles {
			if userRole == role {
				c.Next() // Autorise si l'utilisateur a un rôle autorisé
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Rôle insuffisant"})
		c.Abort()
	}
}

// TokenBlacklistMiddleware - Vérifie si le token actuel est blacklisté
func TokenBlacklistMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Récupérer le token depuis l'en-tête Authorization
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header manquant"})
			c.Abort()
			return
		}

		// Vérifier que le token commence par "Bearer " puis le retirer
		if !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format de token invalide"})
			c.Abort()
			return
		}
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Vérifier si le token est blacklisté
		if isTokenBlacklisted(tokenString) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token invalide, veuillez vous reconnecter"})
			c.Abort()
			return
		}

		c.Next() // Si le token n'est pas blacklisté, continuer
	}
}

// controllers/admin_controller.go
package controllers

import (
	"net/http"

	"github.com/Kensei-s/back/config"
	"github.com/Kensei-s/back/models"

	"github.com/gin-gonic/gin"
)

// Liste des rôles valides qu'un administrateur peut attribuer
var adminValidRoles = []string{"organisateur", "teneur_de_stand", "parent", "eleve", "administrateur"}

// AdminChangeRole - Permet à un administrateur de changer le rôle d'un utilisateur
func AdminChangeRole(c *gin.Context) {
	var input struct {
		Role string `json:"role"`
	}

	// Lier les données reçues
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Valider que le rôle fourni est valide
	if !isValidAdminRole(input.Role) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Rôle invalide. Rôles acceptés : organisateur, teneur_de_stand, parent, eleve, administrateur"})
		return
	}

	// Récupérer l'utilisateur via son ID
	var user models.User
	if err := config.DB.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Utilisateur non trouvé"})
		return
	}

	// Mettre à jour le rôle de l'utilisateur
	user.Role = input.Role
	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour du rôle"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Rôle mis à jour avec succès", "new_role": user.Role})
}

// Fonction pour vérifier si le rôle est valide pour l'administrateur
func isValidAdminRole(role string) bool {
	for _, validRole := range adminValidRoles {
		if role == validRole {
			return true
		}
	}
	return false
}

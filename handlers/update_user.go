package handlers

import (
	"net/http"
	"osohub/db"
	"osohub/middleware"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UpdateUserRequest struct {
	Username          *string `json:"username,omitempty"`
	Bio               *string `json:"bio,omitempty"`
	ProfilePictureURL *string `json:"profile_picture_url,omitempty"`
	Password          *string `json:"password,omitempty"`
}

// UpdateOwnUser allows the authenticated user to update their profile
// @Summary Update the authenticated user's profile
// @Description Allows the user to update their username, bio, profile picture, and/or password
// @Tags Auth & Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body handlers.UpdateUserRequest true "Fields to update"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/me [patch]
func UpdateOwnUser(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	updateFields := make(map[string]interface{})
	if req.Username != nil {
		updateFields["username"] = *req.Username
	}
	if req.Bio != nil {
		updateFields["bio"] = *req.Bio
	}
	if req.ProfilePictureURL != nil {
		updateFields["profile_picture_url"] = *req.ProfilePictureURL
	}
	if req.Password != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
			return
		}
		updateFields["password_hash"] = string(hash)
	}

	if len(updateFields) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	session := db.GetSession()
	if session == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No DB session"})
		return
	}

	setParts := []string{}
	values := []interface{}{}
	for k, v := range updateFields {
		setParts = append(setParts, k+" = ?")
		values = append(values, v)
	}
	query := "UPDATE users_by_id SET " + joinComma(setParts) + " WHERE user_id = ?"
	values = append(values, userID)

	if err := session.Query(query, values...).Exec(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB update failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated"})
}

func joinComma(parts []string) string {
	result := ""
	for i, p := range parts {
		if i > 0 {
			result += ", "
		}
		result += p
	}
	return result
}

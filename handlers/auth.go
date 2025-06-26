package handlers

import (
	"net/http"
	"os"
	"osohub/db"
	"osohub/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// LoginRequest is the expected body for login
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login godoc
// @Summary Login by email and password
// @Accept json
// @Produce json
// @Param login body LoginRequest true "Login credentials"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Tags Auth & Users
// @Router /auth/login [post]
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":         "Invalid login data. Required fields: email, password.",
			"documentation": "https://docs.osohub.com/auth#login",
		})
		return
	}

	var user models.User
	query := `SELECT user_id, username, email, password_hash, profile_picture_url, bio, role, created_at FROM users_by_id WHERE email = ? LIMIT 1 ALLOW FILTERING`
	if err := db.GetSession().Query(query, req.Email).Consistency(gocql.One).Scan(
		&user.UserID, &user.Username, &user.Email, &user.PasswordHash,
		&user.ProfilePictureURL, &user.Bio, &user.Role, &user.CreatedAt,
	); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":         "Invalid credentials.",
			"documentation": "https://docs.osohub.com/auth#login",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":         "Invalid credentials.",
			"documentation": "https://docs.osohub.com/auth#login",
		})
		return
	}

	// Generar JWT
	claims := jwt.MapClaims{
		"user_id": user.UserID.String(),
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":         "Could not generate token",
			"documentation": "https://docs.osohub.com/errors#internal",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": tokenString, "user": user})
}

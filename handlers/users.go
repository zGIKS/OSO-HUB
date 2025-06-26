package handlers

import (
	"net/http"
	"osohub/db"
	"osohub/models"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"golang.org/x/crypto/bcrypt"
)

// CreateUserRequest is the expected body for creating a user
type CreateUserRequest struct {
	Username          string `json:"username" binding:"required"`
	Email             string `json:"email" binding:"required"`
	Password          string `json:"password" binding:"required"`
	ProfilePictureURL string `json:"profile_picture_url"`
	Bio               string `json:"bio"`
}

// CreateUser godoc
// @Summary Create a new user
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "User data"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Router /users [post]
// @Tags Auth & Users
func CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":         "Invalid user data. Required fields: username, email, password.",
			"documentation": "https://docs.osohub.com/users#create",
		})
		return
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":         "Error hashing password. Please try again later.",
			"documentation": "https://docs.osohub.com/errors#internal",
		})
		return
	}

	userID := gocql.TimeUUID()
	createdAt := userID.Time()

	// Check if email already exists
	var exists string
	err = db.GetSession().Query(`SELECT email FROM users_by_id WHERE email = ? ALLOW FILTERING`, req.Email).Scan(&exists)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	// Insert user
	if err := db.GetSession().Query(`INSERT INTO users_by_id (user_id, username, email, password_hash, profile_picture_url, bio, role, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		userID, req.Username, req.Email, string(hash), req.ProfilePictureURL, req.Bio, "user", createdAt).Exec(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}

	user := models.User{
		UserID:            userID,
		Username:          req.Username,
		Email:             req.Email,
		PasswordHash:      string(hash),
		ProfilePictureURL: req.ProfilePictureURL,
		Bio:               req.Bio,
		Role:              "user",
		CreatedAt:         createdAt,
	}
	c.JSON(http.StatusCreated, user)
}

// GetUserByID godoc
// @Summary Get user by ID
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /users/{user_id} [get]
// @Tags Auth & Users
func GetUserByID(c *gin.Context) {
	idStr := c.Param("user_id")
	userID, err := gocql.ParseUUID(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	var user models.User
	query := `SELECT user_id, username, email, password_hash, profile_picture_url, bio, role, created_at FROM users_by_id WHERE user_id = ? LIMIT 1`
	if err := db.GetSession().Query(query, userID).Consistency(gocql.One).Scan(
		&user.UserID, &user.Username, &user.Email, &user.PasswordHash,
		&user.ProfilePictureURL, &user.Bio, &user.Role, &user.CreatedAt,
	); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// BanUser godoc
// @Summary Ban or unban a user
// @Produce json
// @Param user_id path string true "User ID"
// @Param banned query bool true "Ban (true) or unban (false)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /users/{user_id}/ban [patch]
// @Tags Auth & Users
func BanUser(c *gin.Context) {
	idStr := c.Param("user_id")
	userID, err := gocql.ParseUUID(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}
	banned := c.DefaultQuery("banned", "true")
	newRole := "banned"
	if banned == "false" {
		newRole = "user"
	}
	if err := db.GetSession().Query(`UPDATE users_by_id SET role = ? WHERE user_id = ?`, newRole, userID).Exec(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User updated", "role": newRole})
}

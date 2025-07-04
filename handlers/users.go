package handlers

import (
	"fmt"
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

// GetCurrentUser godoc
// @Summary Get current authenticated user info
// @Description Get the profile information of the currently authenticated user
// @Tags Auth & Users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.User
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users/me [get]
func GetCurrentUser(c *gin.Context) {
	// Get user ID from JWT token
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":         "Unauthorized. Valid JWT token required.",
			"documentation": "https://docs.osohub.com/auth#jwt",
		})
		return
	}

	userID, err := gocql.ParseUUID(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id in token"})
		return
	}

	// Query user from database
	var user models.User
	query := `SELECT user_id, username, email, profile_picture_url, bio, role, created_at FROM users_by_id WHERE user_id = ?`
	if err := db.GetSession().Query(query, userID).Scan(
		&user.UserID, &user.Username, &user.Email, &user.ProfilePictureURL, &user.Bio, &user.Role, &user.CreatedAt,
	); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":         "User not found",
			"documentation": "https://docs.osohub.com/users#get",
		})
		return
	}

	// Don't return password hash
	user.PasswordHash = ""

	c.JSON(http.StatusOK, user)
}

// GetPublicProfile godoc
// @Summary Get public profile by username (no authentication required)
// @Produce json
// @Param username path string true "Username"
// @Success 200 {object} map[string]interface{} "Returns user profile and their images"
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /profile/{username} [get]
// @Tags Auth & Users
func GetPublicProfile(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":         "Username is required",
			"documentation": "https://docs.osohub.com/profile#public",
		})
		return
	}

	// Buscar usuario por username
	var user models.User
	query := `SELECT user_id, username, profile_picture_url, bio, created_at FROM users_by_id WHERE username = ? ALLOW FILTERING`
	if err := db.GetSession().Query(query, username).Scan(
		&user.UserID, &user.Username, &user.ProfilePictureURL, &user.Bio, &user.CreatedAt,
	); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":         "User not found",
			"documentation": "https://docs.osohub.com/profile#public",
		})
		return
	}
	// Obtener todas las imágenes del usuario
	var images []gin.H // Usamos gin.H para incluir likes_count
	imageQuery := `SELECT image_id, uploaded_at, user_profile_picture_url, image_url, title FROM images_by_user WHERE user_id = ?`
	iter := db.GetSession().Query(imageQuery, user.UserID).Iter()

	for {
		var image models.Image
		if !iter.Scan(&image.ImageID, &image.UploadedAt, &image.UserProfilePictureURL, &image.ImageURL, &image.Title) {
			break
		}
		// Agregar datos del usuario a cada imagen
		image.UserID = user.UserID
		image.Username = user.Username
		image.UserProfilePictureURL = user.ProfilePictureURL

		// Obtener contador de likes para cada imagen
		var likesCount int64
		countQuery := `SELECT likes FROM image_counters WHERE image_id = ?`
		if err := db.GetSession().Query(countQuery, image.ImageID).Scan(&likesCount); err != nil {
			likesCount = 0 // Si no existe contador, asumimos 0 likes
		}

		// Crear objeto con likes_count incluido
		imageWithLikes := gin.H{
			"image_id":                 image.ImageID,
			"uploaded_at":              image.UploadedAt,
			"user_id":                  image.UserID,
			"username":                 image.Username,
			"user_profile_picture_url": image.UserProfilePictureURL,
			"image_url":                image.ImageURL,
			"title":                    image.Title,
			"likes_count":              likesCount,
		}

		images = append(images, imageWithLikes)
	}

	if err := iter.Close(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":         "Error retrieving user images",
			"documentation": "https://docs.osohub.com/profile#public",
		})
		return
	}

	// Crear respuesta con perfil e imágenes
	response := gin.H{
		"user": gin.H{
			"user_id":             user.UserID,
			"username":            user.Username,
			"profile_picture_url": user.ProfilePictureURL,
			"bio":                 user.Bio,
			"created_at":          user.CreatedAt,
			"total_images":        len(images),
		},
		"images": images,
		"share_url": fmt.Sprintf("%s/profile/%s",
			c.Request.Header.Get("Origin"), username), // URL para compartir
	}

	c.JSON(http.StatusOK, response)
}

// GetMyShareLink godoc
// @Summary Get shareable link for current user's profile
// @Produce json
// @Success 200 {object} map[string]string "Returns share_url"
// @Failure 401 {object} map[string]interface{}
// @Security BearerAuth
// @Router /users/me/share-link [get]
// @Tags Auth & Users
func GetMyShareLink(c *gin.Context) {
	userIDStr, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":         "Unauthorized",
			"documentation": "https://docs.osohub.com/auth#jwt",
		})
		return
	}

	userID, err := gocql.ParseUUID(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":         "Invalid user_id",
			"documentation": "https://docs.osohub.com/auth#jwt",
		})
		return
	}

	// Obtener username del usuario actual
	var username string
	query := `SELECT username FROM users_by_id WHERE user_id = ?`
	if err := db.GetSession().Query(query, userID).Scan(&username); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":         "User not found",
			"documentation": "https://docs.osohub.com/users#share",
		})
		return
	}

	// Construir URL de compartir
	origin := c.Request.Header.Get("Origin")
	if origin == "" {
		origin = "http://localhost:5174" // Fallback para desarrollo
	}

	shareURL := fmt.Sprintf("%s/profile/%s", origin, username)

	c.JSON(http.StatusOK, gin.H{
		"share_url": shareURL,
		"username":  username,
		"message":   "Share this link to let others see your profile and photos!",
	})
}

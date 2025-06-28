package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"osohub/db"
	"osohub/middleware"
	"path/filepath"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
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
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param username formData string false "New username"
// @Param bio formData string false "New bio"
// @Param profile_picture formData file false "New profile picture (JPG, PNG, WebP, max 10MB)"
// @Param password formData string false "New password"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/me [patch]
func UpdateOwnUser(c *gin.Context) {
	// Get user ID from JWT token
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse form data
	username := c.PostForm("username")
	bio := c.PostForm("bio")
	password := c.PostForm("password")

	var profilePictureURL string

	// Handle profile picture upload to Cloudinary
	file, err := c.FormFile("profile_picture")
	if err == nil && file != nil {
		// Validate file type
		fileExt := strings.ToLower(filepath.Ext(file.Filename))
		allowedTypes := []string{".jpg", ".jpeg", ".png", ".webp"}
		isValidType := false
		for _, ext := range allowedTypes {
			if fileExt == ext {
				isValidType = true
				break
			}
		}

		if !isValidType {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Only JPEG, PNG, and WebP are allowed"})
			return
		}

		// Validate file size (max 10MB)
		if file.Size > 10<<20 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File too large. Maximum size is 10MB"})
			return
		}

		// Initialize Cloudinary
		cloudinaryURL := os.Getenv("CLOUDINARY_URL")
		if cloudinaryURL == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Cloudinary configuration missing"})
			return
		}

		cld, err := cloudinary.NewFromURL(cloudinaryURL)
		if err != nil {
			log.Printf("Error initializing Cloudinary: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize Cloudinary"})
			return
		}

		// Open uploaded file
		src, err := file.Open()
		if err != nil {
			log.Printf("Error opening uploaded file: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process uploaded file"})
			return
		}
		defer src.Close()

		// Upload to Cloudinary with profile-specific folder
		uploadResult, err := cld.Upload.Upload(
			context.Background(),
			src,
			uploader.UploadParams{
				Folder:         "osohub-profiles",
				Transformation: "c_fill,w_300,h_300,g_face", // Square crop focusing on face
				PublicID:       fmt.Sprintf("profile_%s_%d", userID, time.Now().Unix()),
			},
		)

		if err != nil {
			log.Printf("Error uploading to Cloudinary: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload profile picture"})
			return
		}

		profilePictureURL = uploadResult.SecureURL
		log.Printf("Profile picture uploaded successfully: %s", profilePictureURL)
	}

	// Build update fields dynamically
	updateFields := make(map[string]interface{})

	if username != "" {
		updateFields["username"] = username
	}

	if bio != "" {
		updateFields["bio"] = bio
	}

	if profilePictureURL != "" {
		updateFields["profile_picture_url"] = profilePictureURL
	}

	if password != "" {
		// Hash the new password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		updateFields["password_hash"] = string(hashedPassword)
	}

	if len(updateFields) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	// Build and execute query
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
	query := "UPDATE users_by_id SET " + strings.Join(setParts, ", ") + " WHERE user_id = ?"
	values = append(values, userID)

	if err := session.Query(query, values...).Exec(); err != nil {
		log.Printf("Error updating user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB update failed"})
		return
	}
	// Si se actualizó la foto de perfil o el username, actualizar todas las imágenes del usuario
	if profilePictureURL != "" || username != "" {
		log.Printf("Updating user info in all user images...")

		// Obtener el username y profile_picture_url actuales del usuario después de la actualización
		var currentUsername, currentProfilePictureURL string
		userQuery := `SELECT username, profile_picture_url FROM users_by_id WHERE user_id = ?`
		if err := session.Query(userQuery, userID).Scan(&currentUsername, &currentProfilePictureURL); err != nil {
			log.Printf("Error getting updated user info: %v", err)
		} else {
			// Primero obtener todas las imágenes del usuario desde images_by_user
			// porque es la única tabla que permite filtrar eficientemente por user_id
			userImagesQuery := `SELECT uploaded_at, image_id FROM images_by_user WHERE user_id = ?`
			iter := session.Query(userImagesQuery, userID).Iter()

			var uploadedAt time.Time
			var imageID gocql.UUID
			var imagesToUpdate []struct {
				ImageID    gocql.UUID
				UploadedAt time.Time
				DayBucket  string
			}

			// Recopilar todas las imágenes del usuario
			for iter.Scan(&uploadedAt, &imageID) {
				dayBucket := uploadedAt.Format("2006-01-02")
				imagesToUpdate = append(imagesToUpdate, struct {
					ImageID    gocql.UUID
					UploadedAt time.Time
					DayBucket  string
				}{imageID, uploadedAt, dayBucket})
			}
			iter.Close()

			if len(imagesToUpdate) > 0 {
				// Actualizar images_by_user
				log.Printf("Updating %d images in images_by_user...", len(imagesToUpdate))
				for _, img := range imagesToUpdate {
					updateUserImageQuery := `UPDATE images_by_user SET user_profile_picture_url = ? WHERE user_id = ? AND uploaded_at = ? AND image_id = ?`
					if err := session.Query(updateUserImageQuery, currentProfilePictureURL, userID, img.UploadedAt, img.ImageID).Exec(); err != nil {
						log.Printf("Error updating images_by_user for image %v: %v", img.ImageID, err)
					}
				}

				// Actualizar images_by_id (una por una usando image_id)
				log.Printf("Updating %d images in images_by_id...", len(imagesToUpdate))
				for _, img := range imagesToUpdate {
					updateImageQuery := `UPDATE images_by_id SET username = ?, user_profile_picture_url = ? WHERE image_id = ?`
					if err := session.Query(updateImageQuery, currentUsername, currentProfilePictureURL, img.ImageID).Exec(); err != nil {
						log.Printf("Error updating images_by_id for image %v: %v", img.ImageID, err)
					}
				}

				// Actualizar images_by_date
				log.Printf("Updating %d images in images_by_date...", len(imagesToUpdate))
				for _, img := range imagesToUpdate {
					updateDateImageQuery := `UPDATE images_by_date SET username = ?, user_profile_picture_url = ? WHERE day_bucket = ? AND uploaded_at = ? AND image_id = ?`
					if err := session.Query(updateDateImageQuery, currentUsername, currentProfilePictureURL, img.DayBucket, img.UploadedAt, img.ImageID).Exec(); err != nil {
						log.Printf("Error updating images_by_date for image %v: %v", img.ImageID, err)
					}
				}

				log.Printf("Finished updating user info in all images")
			} else {
				log.Printf("No images found for user %s", userID)
			}
		}
	}

	response := gin.H{"message": "Profile updated successfully"}
	if profilePictureURL != "" {
		response["profile_picture_url"] = profilePictureURL
	}

	c.JSON(http.StatusOK, response)
}

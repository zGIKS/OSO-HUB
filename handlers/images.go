package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"osohub/db"
	"osohub/models"
	"path/filepath"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
)

// LikeImage godoc
// @Summary Like an image (one like per user)
// @Param image_id path string true "Image ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Security BearerAuth
// @Router /images/{image_id}/like [post]
// @Tags Images
func LikeImage(c *gin.Context) {
	imageIDStr := c.Param("image_id")
	imageID, err := gocql.ParseUUID(imageIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":         "Invalid image_id. Must be a valid UUID.",
			"documentation": "https://docs.osohub.com/images#like",
		})
		return
	}
	userIDStr, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":         "Unauthorized.",
			"documentation": "https://docs.osohub.com/auth#jwt",
		})
		return
	}
	userID, err := gocql.ParseUUID(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":         "Invalid user_id. Must be a valid UUID.",
			"documentation": "https://docs.osohub.com/auth#jwt",
		})
		return
	}
	var exists gocql.UUID
	err = db.GetSession().Query(`SELECT user_id FROM likes_by_image WHERE image_id = ? AND user_id = ?`, imageID, userID).Scan(&exists)
	if err == nil {
		c.Status(http.StatusNoContent)
		return
	}
	now := time.Now()
	if err := db.GetSession().Query(`INSERT INTO likes_by_image (image_id, user_id, liked_at) VALUES (?, ?, ?)`, imageID, userID, now).Exec(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":         "Error liking image. Please try again later.",
			"documentation": "https://docs.osohub.com/errors#internal",
		})
		return
	}
	if err := db.GetSession().Query(`UPDATE image_counters SET likes = likes + 1 WHERE image_id = ?`, imageID).Exec(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":         "Error updating like counter. Please try again later.",
			"documentation": "https://docs.osohub.com/errors#internal",
		})
		return
	}
	c.Status(http.StatusNoContent)
}

// UnlikeImage godoc
// @Summary Remove like from an image
// @Param image_id path string true "Image ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Security BearerAuth
// @Router /images/{image_id}/like [delete]
// @Tags Images
func UnlikeImage(c *gin.Context) {
	imageIDStr := c.Param("image_id")
	imageID, err := gocql.ParseUUID(imageIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":         "Invalid image_id. Must be a valid UUID.",
			"documentation": "https://docs.osohub.com/images#like",
		})
		return
	}
	userIDStr, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":         "Unauthorized.",
			"documentation": "https://docs.osohub.com/auth#jwt",
		})
		return
	}
	userID, err := gocql.ParseUUID(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":         "Invalid user_id. Must be a valid UUID.",
			"documentation": "https://docs.osohub.com/auth#jwt",
		})
		return
	}
	if err := db.GetSession().Query(`DELETE FROM likes_by_image WHERE image_id = ? AND user_id = ?`, imageID, userID).Exec(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":         "Error removing like. Please try again later.",
			"documentation": "https://docs.osohub.com/errors#internal",
		})
		return
	}
	if err := db.GetSession().Query(`UPDATE image_counters SET likes = likes - 1 WHERE image_id = ?`, imageID).Exec(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":         "Error updating like counter. Please try again later.",
			"documentation": "https://docs.osohub.com/errors#internal",
		})
		return
	}
	c.Status(http.StatusNoContent)
}

// GetImageLikesCount godoc
// @Summary Get like count for an image
// @Param image_id path string true "Image ID"
// @Success 200 {object} map[string]int64
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /images/{image_id}/likes/count [get]
// @Tags Images
func GetImageLikesCount(c *gin.Context) {
	imageIDStr := c.Param("image_id")
	imageID, err := gocql.ParseUUID(imageIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":         "Invalid image_id. Must be a valid UUID.",
			"documentation": "https://docs.osohub.com/images#like",
		})
		return
	}
	var likes int64
	err = db.GetSession().Query(`SELECT likes FROM image_counters WHERE image_id = ?`, imageID).Scan(&likes)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":         "Image not found or no likes.",
			"documentation": "https://docs.osohub.com/images#like",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"likes": likes})
}

// DeleteImage godoc
// @Summary Delete an image by ID (solo el dueño puede borrar)
// @Param image_id path string true "Image ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Security BearerAuth
// @Router /images/{image_id} [delete]
// @Tags Images
func DeleteImage(c *gin.Context) {
	imageID := c.Param("image_id")
	if _, err := uuid.Parse(imageID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":         "Invalid image_id. Must be a valid UUID.",
			"documentation": "https://docs.osohub.com/images#delete",
		})
		return
	}

	// Validar que el usuario autenticado es el dueño de la imagen
	userIDStr, ok := c.Get("user_id")
	if !ok {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	var ownerID gocql.UUID
	var dayBucket string
	var uploadedAt time.Time
	err := db.GetSession().Query(`SELECT user_id, day_bucket, uploaded_at FROM images_by_id WHERE image_id = ?`, imageID).Scan(&ownerID, &dayBucket, &uploadedAt)
	if err != nil {
		c.JSON(404, gin.H{"error": "Image not found"})
		return
	}
	if ownerID.String() != userIDStr {
		c.JSON(403, gin.H{"error": "You are not the owner of this image"})
		return
	}

	// Borra de images_by_id
	if err := db.GetSession().Query(`DELETE FROM images_by_id WHERE image_id = ?`, imageID).Exec(); err != nil {
		c.JSON(500, gin.H{"error": "Error deleting image (by_id)"})
		return
	}
	// Borra de images_by_date
	if err := db.GetSession().Query(`DELETE FROM images_by_date WHERE day_bucket = ? AND uploaded_at = ? AND image_id = ?`, dayBucket, uploadedAt, imageID).Exec(); err != nil {
		c.JSON(500, gin.H{
			"error":         "Error deleting image (by_date). Please try again later.",
			"documentation": "https://docs.osohub.com/errors#internal",
		})
		return
	}
	// Borra de images_by_user
	if err := db.GetSession().Query(`DELETE FROM images_by_user WHERE user_id = ? AND uploaded_at = ? AND image_id = ?`, ownerID, uploadedAt, imageID).Exec(); err != nil {
		c.JSON(500, gin.H{"error": "Error deleting image (by_user)"})
		return
	}
	// Borra de image_counters
	if err := db.GetSession().Query(`DELETE FROM image_counters WHERE image_id = ?`, imageID).Exec(); err != nil {
		c.JSON(500, gin.H{"error": "Error deleting image (counters)"})
		return
	}
	// Borra de likes_by_image
	if err := db.GetSession().Query(`DELETE FROM likes_by_image WHERE image_id = ?`, imageID).Exec(); err != nil {
		c.JSON(500, gin.H{"error": "Error deleting image (likes)"})
		return
	}
	// Borra de reports_by_image
	if err := db.GetSession().Query(`DELETE FROM reports_by_image WHERE image_id = ?`, imageID).Exec(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":         "Error deleting image (reports). Please try again later.",
			"documentation": "https://docs.osohub.com/errors#internal",
		})
		return
	}
	c.Status(204)
}

// UploadImageRequest is the expected body for uploading an image (deprecated, now uses form-data)
type UploadImageRequest struct {
	ImageURL string `json:"image_url" binding:"required"`
	Title    string `json:"title" binding:"required"`
}

// UploadImage godoc
// @Summary Upload a new image file to Cloudinary
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param image formData file true "Image file (JPG, PNG, GIF, WebP, max 10MB)"
// @Param title formData string true "Image title (max 100 characters)"
// @Success 201 {object} models.Image
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /images [post]
// @Tags Images
func UploadImage(c *gin.Context) {
	// Obtener user_id del token JWT
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

	// Obtener archivo de imagen
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No image file uploaded. Use 'image' field in form-data.",
		})
		return
	}

	// Validar tipo de archivo
	allowedTypes := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	fileExt := strings.ToLower(filepath.Ext(file.Filename))
	isValidType := false
	for _, ext := range allowedTypes {
		if fileExt == ext {
			isValidType = true
			break
		}
	}
	if !isValidType {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid file type. Allowed: JPG, PNG, GIF, WebP",
		})
		return
	}

	// Validar tamaño (máximo 10MB)
	if file.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "File too large. Maximum size: 10MB",
		})
		return
	}

	// Obtener título
	title := c.PostForm("title")
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Title is required",
		})
		return
	}

	// Validar longitud del título
	if len(title) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Title too long. Maximum 100 characters.",
		})
		return
	}
	// Configurar Cloudinary usando CLOUDINARY_URL (método recomendado)
	cloudinaryURL := os.Getenv("CLOUDINARY_URL")
	if cloudinaryURL == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "CLOUDINARY_URL not configured"})
		return
	}

	fmt.Printf("DEBUG Cloudinary URL: %s\n", cloudinaryURL)

	cld, err := cloudinary.NewFromURL(cloudinaryURL)
	if err != nil {
		fmt.Printf("DEBUG Cloudinary config error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cloudinary configuration error"})
		return
	}

	// Abrir archivo
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not open uploaded file"})
		return
	}
	defer src.Close()

	// Subir a Cloudinary
	uploadParams := uploader.UploadParams{
		PublicID: fmt.Sprintf("osohub/%d_%s", time.Now().Unix(), strings.TrimSuffix(file.Filename, fileExt)),
		Folder:   "osohub-images",
	}

	fmt.Printf("DEBUG Upload params: %+v\n", uploadParams)

	result, err := cld.Upload.Upload(context.Background(), src, uploadParams)
	if err != nil {
		fmt.Printf("DEBUG Cloudinary upload error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image to Cloudinary", "details": err.Error()})
		return
	}

	fmt.Printf("DEBUG Cloudinary result: %+v\n", result)

	// URL de Cloudinary (ya optimizada)
	imageURL := result.SecureURL
	fmt.Printf("DEBUG imageURL: %s\n", imageURL)

	// Obtener información del usuario para el username
	var username string
	if err := db.GetSession().Query(`SELECT username FROM users_by_id WHERE user_id = ?`, userID).Scan(&username); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	imageID := gocql.TimeUUID()
	uploadedAt := imageID.Time()
	dayBucket := uploadedAt.Format("2006-01-02")

	// Insert into images_by_id
	if err := db.GetSession().Query(`INSERT INTO images_by_id (image_id, day_bucket, uploaded_at, user_id, username, image_url, title) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		imageID, dayBucket, uploadedAt, userID, username, imageURL, title).Exec(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving image (by_id)"})
		return
	}

	// Insert into images_by_date
	if err := db.GetSession().Query(`INSERT INTO images_by_date (day_bucket, uploaded_at, image_id, user_id, username, image_url, title) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		dayBucket, uploadedAt, imageID, userID, username, imageURL, title).Exec(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":         "Could not save image (by_date). Please try again later.",
			"documentation": "https://docs.osohub.com/errors#internal",
		})
		return
	}

	// Insert into images_by_user
	if err := db.GetSession().Query(`INSERT INTO images_by_user (user_id, uploaded_at, image_id, image_url, title) VALUES (?, ?, ?, ?, ?)`,
		userID, uploadedAt, imageID, imageURL, title).Exec(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving image (by_user)"})
		return
	}

	image := models.Image{
		ImageID:    imageID,
		DayBucket:  dayBucket,
		UploadedAt: uploadedAt,
		UserID:     userID,
		Username:   username,
		ImageURL:   imageURL,
		Title:      title,
	}
	c.JSON(http.StatusCreated, image)
}

// GetImagesByUser godoc
// @Summary Get all images for a user (profile)
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {array} models.Image
// @Failure 400 {object} map[string]interface{}
// @Router /users/{user_id}/images [get]
// @Tags Images
func GetImagesByUser(c *gin.Context) {
	idStr := c.Param("user_id")
	userID, err := gocql.ParseUUID(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":         "Invalid user_id. Must be a valid UUID.",
			"documentation": "https://docs.osohub.com/users#images",
		})
		return
	}
	query := `SELECT uploaded_at, image_id, image_url, title FROM images_by_user WHERE user_id = ?`
	iter := db.GetSession().Query(query, userID).Iter()
	var images []models.Image
	var img models.Image
	for iter.Scan(&img.UploadedAt, &img.ImageID, &img.ImageURL, &img.Title) {
		img.UserID = userID
		images = append(images, img)
	}
	if err := iter.Close(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":         "Could not fetch images. Please try again later.",
			"documentation": "https://docs.osohub.com/errors#internal",
		})
		return
	}
	c.JSON(http.StatusOK, images)
}

// GetImageByIDByOnlyID godoc
// @Summary Get image by ID (direct, for Swagger compatibility)
// @Produce json
// @Param image_id path string true "Image ID"
// @Success 200 {object} models.Image
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /images/byid/{image_id} [get]
// @Tags Images
func GetImageByIDByOnlyID(c *gin.Context) {
	idStr := c.Param("image_id")
	imageID, err := gocql.ParseUUID(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":         "Invalid image_id. Must be a valid UUID.",
			"documentation": "https://docs.osohub.com/images#byid",
		})
		return
	}
	var image models.Image
	query := `SELECT image_id, day_bucket, uploaded_at, user_id, username, image_url, title FROM images_by_id WHERE image_id = ? LIMIT 1`
	if err := db.GetSession().Query(query, imageID).Consistency(gocql.One).Scan(
		&image.ImageID, &image.DayBucket, &image.UploadedAt, &image.UserID, &image.Username, &image.ImageURL, &image.Title,
	); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":         "Image not found.",
			"documentation": "https://docs.osohub.com/images#byid",
		})
		return
	}
	c.JSON(http.StatusOK, image)
}

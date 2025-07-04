package handlers

import (
	"net/http"
	"osohub/db"
	"osohub/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetFeed godoc
// @Summary Get global image feed (latest images)
// @Produce json
// @Tags Images
// @Param day_bucket query string false "Day bucket (YYYY-MM-DD)"
// @Param limit query int false "Limit"
// @Success 200 {array} models.Image
// @Router /feed [get]
func GetFeed(c *gin.Context) {
	dayBucket := c.DefaultQuery("day_bucket", "")
	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 20
	}
	if dayBucket == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":         "day_bucket is required (YYYY-MM-DD)",
			"documentation": "https://docs.osohub.com/images#feed",
		})
		return
	}

	// Get images from images_by_date
	query := `SELECT image_id, user_id, username, image_url, title, uploaded_at FROM images_by_date WHERE day_bucket = ? LIMIT ?`
	iter := db.GetSession().Query(query, dayBucket, limit).Iter()
	var images []models.Image
	var img models.Image

	// Store unique user IDs to fetch their current profile pictures
	userProfilePictures := make(map[string]string)

	// First, collect all images and unique user IDs
	var tempImages []models.Image
	for iter.Scan(&img.ImageID, &img.UserID, &img.Username, &img.ImageURL, &img.Title, &img.UploadedAt) {
		tempImages = append(tempImages, img)
		userProfilePictures[img.UserID.String()] = "" // Initialize with empty string
	}
	if err := iter.Close(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":         "Could not fetch feed. Please try again later.",
			"documentation": "https://docs.osohub.com/errors#internal",
		})
		return
	}

	// Fetch current profile pictures for all unique users
	for userIDStr := range userProfilePictures {
		var profilePictureURL string
		userQuery := `SELECT profile_picture_url FROM users_by_id WHERE user_id = ?`
		if err := db.GetSession().Query(userQuery, userIDStr).Scan(&profilePictureURL); err == nil {
			userProfilePictures[userIDStr] = profilePictureURL
		}
	}

	// Now build the final images array with current profile pictures
	for _, img := range tempImages {
		img.UserProfilePictureURL = userProfilePictures[img.UserID.String()]
		images = append(images, img)
	}

	c.JSON(http.StatusOK, images)
}

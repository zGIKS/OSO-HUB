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

	query := `SELECT image_id, user_id, username, image_url, title, uploaded_at FROM images_by_date WHERE day_bucket = ? LIMIT ?`
	iter := db.Session.Query(query, dayBucket, limit).Iter()
	var images []models.Image
	var img models.Image
	for iter.Scan(&img.ImageID, &img.UserID, &img.Username, &img.ImageURL, &img.Title, &img.UploadedAt) {
		images = append(images, img)
	}
	if err := iter.Close(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":         "Could not fetch feed. Please try again later.",
			"documentation": "https://docs.osohub.com/errors#internal",
		})
		return
	}
	c.JSON(http.StatusOK, images)
}

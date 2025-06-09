package handlers

import (
	"net/http"
	"osohub/db"
	"osohub/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
)

// ReportImage godoc
// @Summary Report an image
// @Description Report an image for inappropriate content or other reasons
// @Tags Images
// @Security BearerAuth
// @Param image_id path string true "Image ID"
// @Accept json
// @Produce json
// @Param report body models.ReportRequest true "Report reason"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /images/{image_id}/report [post]
func ReportImage(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	imageIDStr := c.Param("image_id")
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":         "Unauthorized. Missing user_id in token.",
			"documentation": "https://docs.osohub.com/auth#jwt",
		})
		return
	}
	imageUUID, err := uuid.Parse(imageIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":         "Invalid image_id. Must be a valid UUID.",
			"documentation": "https://docs.osohub.com/images#report",
		})
		return
	}
	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":         "Invalid user_id. Must be a valid UUID.",
			"documentation": "https://docs.osohub.com/auth#jwt",
		})
		return
	}
	var req models.ReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":         "Reason is required and must be a non-empty string.",
			"documentation": "https://docs.osohub.com/images#report",
		})
		return
	}

	timeUUID := gocql.TimeUUID().String()
	now := time.Now().UTC()

	if err := db.InsertReport(timeUUID, imageUUID.String(), userUUID.String(), req.Reason, now); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":         "Could not report image. Please try again later.",
			"documentation": "https://docs.osohub.com/errors#internal",
			"debug":         err.Error(),
		})
		return
	}
	if err := db.IncrementImageReportCounter(imageIDStr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":         "Could not update report counter. Please try again later.",
			"documentation": "https://docs.osohub.com/errors#internal",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Image reported"})
}

// GetImageReportsCount godoc
// @Summary Get report count for an image
// @Description Get the number of reports for a specific image
// @Tags Images
// @Param image_id path string true "Image ID"
// @Produce json
// @Success 200 {object} map[string]int
// @Failure 500 {object} map[string]string
// @Router /images/{image_id}/reports/count [get]
func GetImageReportsCount(c *gin.Context) {
	imageID := c.Param("image_id")
	if _, err := uuid.Parse(imageID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":         "Invalid image_id. Must be a valid UUID.",
			"documentation": "https://docs.osohub.com/images#report",
		})
		return
	}
	count, err := db.GetImageReportCount(imageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":         "Could not fetch report count. Please try again later.",
			"documentation": "https://docs.osohub.com/errors#internal",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": count})
}

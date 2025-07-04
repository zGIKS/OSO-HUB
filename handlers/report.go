package handlers

import (
	"net/http"
	"osohub/db"
	"osohub/models"
	"strconv"
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
			"error":         "Category is required. Reason is optional.",
			"documentation": "https://docs.osohub.com/images#report",
		})
		return
	}

	// Validar que la categoría sea válida
	validCategories := map[string]bool{
		"harassment":     true,
		"hate":           true,
		"spam":           true,
		"inappropriate":  true,
		"violence":       true,
		"misinformation": true,
		"copyright":      true,
		"other":          true,
	}

	if !validCategories[req.Category] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":         "Invalid category. Use /reports/categories to get valid categories.",
			"documentation": "https://docs.osohub.com/images#report",
		})
		return
	}

	timeUUID := gocql.TimeUUID()
	now := time.Now().UTC()

	// Insertar en reports_by_image (tabla principal)
	query1 := `INSERT INTO reports_by_image (image_id, report_id, reporter_id, category, reason, reported_at) VALUES (?, ?, ?, ?, ?, ?)`
	if err := db.GetSession().Query(query1, gocql.UUID(imageUUID), timeUUID, gocql.UUID(userUUID), req.Category, req.Reason, now).Exec(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":         "Could not report image. Please try again later.",
			"documentation": "https://docs.osohub.com/errors#internal",
		})
		return
	}

	// Insertar en reports_by_category (para análisis)
	query2 := `INSERT INTO reports_by_category (category, reported_at, report_id, image_id, reporter_id, reason) VALUES (?, ?, ?, ?, ?, ?)`
	if err := db.GetSession().Query(query2, req.Category, now, timeUUID, gocql.UUID(imageUUID), gocql.UUID(userUUID), req.Reason).Exec(); err != nil {
		// No fallar si esto falla, pero log el error
		// Este insert es para análisis, no crítico
	}

	// Incrementar contador
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

// ReportCategory representa una categoría de reporte
type ReportCategory struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// GetReportCategories godoc
// @Summary Get available report categories
// @Description Get list of available categories for reporting content
// @Tags Images
// @Produce json
// @Success 200 {object} map[string][]ReportCategory
// @Router /reports/categories [get]
func GetReportCategories(c *gin.Context) {
	categories := []ReportCategory{
		{
			ID:          "harassment",
			Name:        "Harassment",
			Description: "Content that harasses, intimidates or bothers other users",
		},
		{
			ID:          "hate",
			Name:        "Hate Speech",
			Description: "Content that promotes hatred against groups or individuals",
		},
		{
			ID:          "spam",
			Name:        "Spam",
			Description: "Repetitive, unwanted or promotional content",
		},
		{
			ID:          "inappropriate",
			Name:        "Inappropriate Content",
			Description: "Explicit sexual content or material not suitable for all audiences",
		},
		{
			ID:          "violence",
			Name:        "Violence",
			Description: "Content that shows or promotes violence",
		},
		{
			ID:          "misinformation",
			Name:        "Misinformation",
			Description: "False or misleading information",
		},
		{
			ID:          "copyright",
			Name:        "Copyright",
			Description: "Unauthorized use of copyrighted content",
		},
		{
			ID:          "other",
			Name:        "Other",
			Description: "Other reasons not listed above",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
	})
}

// GetReportsByCategory godoc
// @Summary Get reports grouped by category (Admin only)
// @Description Get reports grouped by category for moderation purposes
// @Tags Images
// @Security BearerAuth
// @Param category query string false "Filter by specific category"
// @Param limit query int false "Limit number of results (default: 50)"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /reports/by-category [get]
func GetReportsByCategory(c *gin.Context) {
	// TODO: Agregar verificación de rol de admin
	// userRole := c.GetString("user_role")
	// if userRole != "admin" {
	//     c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
	//     return
	// }

	category := c.Query("category")
	limit := 50 // default
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 100 {
			limit = parsedLimit
		}
	}

	var query string
	var args []interface{}

	if category != "" {
		query = `SELECT category, reported_at, report_id, image_id, reporter_id, reason FROM reports_by_category WHERE category = ? LIMIT ?`
		args = []interface{}{category, limit}
	} else {
		query = `SELECT category, reported_at, report_id, image_id, reporter_id, reason FROM reports_by_category LIMIT ?`
		args = []interface{}{limit}
	}

	iter := db.GetSession().Query(query, args...).Iter()

	reports := make(map[string][]gin.H)

	for {
		var reportCategory, reason string
		var reportedAt time.Time
		var reportID, imageID, reporterID gocql.UUID

		if !iter.Scan(&reportCategory, &reportedAt, &reportID, &imageID, &reporterID, &reason) {
			break
		}

		report := gin.H{
			"report_id":   reportID,
			"image_id":    imageID,
			"reporter_id": reporterID,
			"reason":      reason,
			"reported_at": reportedAt,
		}

		reports[reportCategory] = append(reports[reportCategory], report)
	}

	if err := iter.Close(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error retrieving reports",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"reports_by_category": reports,
		"total_categories":    len(reports),
	})
}

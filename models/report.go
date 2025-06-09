package models

// ReportRequest represents the request body for reporting an image
// swagger:model
type ReportRequest struct {
	Reason string `json:"reason" binding:"required"`
}

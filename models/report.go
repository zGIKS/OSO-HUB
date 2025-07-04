package models

// ReportRequest represents the request body for reporting an image
// swagger:model
type ReportRequest struct {
	Category string `json:"category" binding:"required"` // ID de la categoría (ej: "harassment", "hate", etc.)
	Reason   string `json:"reason"`                      // Descripción adicional opcional
}

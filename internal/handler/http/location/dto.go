package location

import (
	"time"

	"github.com/google/uuid"
)

type CheckLocationRequest struct {
	UserID uuid.UUID `json:"user_id"`
	Lat    float64   `json:"lat"`
	Lng    float64   `json:"lng"`
}

type LocationResponse struct {
	ID          uuid.UUID   `json:"id"`
	UserID      uuid.UUID   `json:"user_id"`
	Lat         float64     `json:"lat"`
	Lng         float64     `json:"lng"`
	Timestamp   time.Time   `json:"timestamp"`
	IsCheck     bool        `json:"is_check"`
	IncidentIDs []uuid.UUID `json:"incident_ids"`
}

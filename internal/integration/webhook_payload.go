package integration

import "github.com/google/uuid"

type WebhookPayload struct {
	UserID      uuid.UUID   `json:"user_id"`
	Lat         float64     `json:"lat"`
	Lng         float64     `json:"lng"`
	IncidentIDs []uuid.UUID `json:"incident_ids"`
	Timestamp   int64       `json:"timestamp"`
	Retry       int         `json:"retry"`
}

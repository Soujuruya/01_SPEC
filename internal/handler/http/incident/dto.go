package incident

type CreateIncidentRequest struct {
	Title    string  `json:"title"`
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
	Radius   float64 `json:"radius"`
	IsActive bool    `json:"is_active"`
}

type UpdateIncidentRequest struct {
	Title    *string  `json:"title"`
	Lat      *float64 `json:"lat"`
	Lng      *float64 `json:"lng"`
	Radius   *float64 `json:"radius"`
	IsActive *bool    `json:"is_active"`
}

type IncidentResponse struct {
	ID        string  `json:"id"`
	Title     string  `json:"title"`
	Lat       float64 `json:"lat"`
	Lng       float64 `json:"lng"`
	Radius    float64 `json:"radius"`
	IsActive  bool    `json:"is_active"`
	CreatedAt string  `json:"created_at"`
}

type IncidentListResponse struct {
	Incidents []IncidentResponse `json:"incidents"`
	Limit     int                `json:"limit"`
	Offset    int                `json:"offset"`
	Total     int                `json:"total"`
}
